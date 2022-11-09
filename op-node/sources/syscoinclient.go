package sources

import (
	"context"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"runtime/debug"
	"time"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)
// JSONMarshalerV2 is used for marshalling requests to newer Syscoin Type RPC interfaces
type JSONMarshalerV2 struct{}

// Marshal converts struct passed by parameter to JSON
func (JSONMarshalerV2) Marshal(v interface{}) ([]byte, error) {
	d, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return d, nil
}
// SyscoinRPC is an interface to JSON-RPC syscoind service.
type SyscoinRPC struct {
	client       http.Client
	rpcURL       string
	user         string
	password     string
	RPCMarshaler JSONMarshalerV2
}
type SyscoinClient struct {
	client *SyscoinRPC
}
func NewSyscoinClient() SyscoinClient {
	transport := &http.Transport{
		Dial:                (&net.Dialer{KeepAlive: 600 * time.Second}).Dial,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100, // necessary to not to deplete ports
	}

	s := &SyscoinRPC{
		client:       http.Client{Timeout: time.Duration(25) * time.Second, Transport: transport},
		rpcURL:       "http://l1:18370/wallet/wallet",
		user:         "u",
		password:     "p",
		RPCMarshaler: JSONMarshalerV2{},
	}

	return SyscoinClient{s}
}
// RPCError defines rpc error returned by backend
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
func (e *RPCError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
func safeDecodeResponse(body io.ReadCloser, res interface{}) (err error) {
	var data []byte
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			if len(data) > 0 && len(data) < 2048 {
				err = errors.New(fmt.Sprintf("Error %v", string(data)))
			} else {
				err = errors.New("Internal error")
			}
		}
	}()
	data, err = ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &res)
}

// Call calls Backend RPC interface, using RPCMarshaler interface to marshall the request
func (s *SyscoinClient) Call(req interface{}, res interface{}) error {
	httpData, err := s.client.RPCMarshaler.Marshal(req)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequest("POST", s.client.rpcURL, bytes.NewBuffer(httpData))
	if err != nil {
		return err
	}
	httpReq.SetBasicAuth(s.client.user, s.client.password)
	httpRes, err := s.client.client.Do(httpReq)
	// in some cases the httpRes can contain data even if it returns error
	// see http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/
	if httpRes != nil {
		defer httpRes.Body.Close()
	}
	if err != nil {
		return err
	}
	// if server returns HTTP error code it might not return json with response
	// handle both cases
	if httpRes.StatusCode != 200 {
		err = safeDecodeResponse(httpRes.Body, &res)
		if err != nil {
			return errors.New(fmt.Sprintf("Error %v %v", httpRes.Status, err))
		}
		return nil
	}
	return safeDecodeResponse(httpRes.Body, &res)
}

func (s *SyscoinClient) CreateBlob(data []byte) (common.Hash, error) {
	type ResCreateBlob struct {
		Error  *RPCError `json:"error"`
		Result struct {
			TXID string `json:"txid"`
		} `json:"result"`
	}

	res := ResCreateBlob{}
	type CmdCreateBlob struct {
		Method string `json:"method"`
		Params struct {
			Data string `json:"data"`
		} `json:"params"`
	}
	req := CmdCreateBlob{Method: "syscoincreatenevmblob"}
	req.Params.Data = string(data)
	err := s.Call(&req, &res)
	if err != nil {
		return common.Hash{}, err
	}
	if res.Error != nil {
		return common.Hash{}, res.Error
	}
	return common.HexToHash(res.Result.TXID), err
}
// SYSCOIN used to get blob confirmation by checking block number then tx receipt below to get block height of blob confirmation
func (s *SyscoinClient) BlockNumber(ctx context.Context) (uint64, error) {
	type ResGetBlockNumber struct {
		Error  *RPCError `json:"error"`
		BlockNumber uint64 `json:"blocknumber"`
	}
	res := ResGetBlockNumber{}
	type CmdGetBlockNumber struct {
		Method string `json:"method"`
		Params struct {
		} `json:"params"`
	}
	req := CmdGetBlockNumber{Method: "getblockcount"}
	err := s.Call(&req, &res)
	if err != nil {
		return 0, err
	}
	if res.Error != nil {
		return 0, res.Error
	}
	return res.BlockNumber, err
}
// SYSCOIN used to get blob receipt
func (s *SyscoinClient) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	type ResGetBlobReceipt struct {
		Error  *RPCError `json:"error"`
		Result struct {
			VH string `json:"versionhash"`
			BlockHash string `json:"blockhash"`
			BlockHeight int64 `json:"height"`
			MPT int64 `json:"mpt"`
		} `json:"result"`
	}
	res := ResGetBlobReceipt{}
	type CmdGetBlobReceipt struct {
		Method string `json:"method"`
		Params struct {
			TXID string `json:"versionhash_or_txid"`
		} `json:"params"`
	}
	req := CmdGetBlobReceipt{Method: "getnevmblobdata"}
	req.Params.TXID = txHash.String()[2:]
	err := s.Call(&req, &res)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, res.Error
	}
	receipt := types.Receipt{}
	if res.Result.MPT > 0 && len(res.Result.BlockHash) > 0 {
		// store VH in txhash used by driver to put into BatchInbox
		receipt = types.Receipt{
			TxHash:      common.HexToHash(res.Result.VH),
			BlockNumber: big.NewInt(res.Result.BlockHeight),
			BlockHash:   common.HexToHash(res.Result.BlockHash),
			Status:      types.ReceiptStatusSuccessful,
		}
	}
	return &receipt, err
}

func (s *SyscoinClient) GetBlobFromRPC(vh common.Hash) (string, error) {
	type ResGetBlobData struct {
		Error  *RPCError `json:"error"`
		Result struct {
			Data string `json:"data"`
		} `json:"result"`
	}
	res := ResGetBlobData{}
	type CmdGetBlobData struct {
		Method string `json:"method"`
		Params struct {
			VersionHash string `json:"versionhash_or_txid"`
			Verbose   bool   `json:"getdata"`
		} `json:"params"`
	}
	req := CmdGetBlobData{Method: "getnevmblobdata"}
	req.Params.VersionHash = vh.String()[2:]
	req.Params.Verbose = true
	err := s.Call(&req, &res)
	if err != nil {
		return "", err
	}
	if res.Error != nil {
		return "", res.Error
	}
	return res.Result.Data, err
}

func (s *SyscoinClient) GetBlobFromCloud(vh common.Hash) (string, error) {
	return "", nil
}