package derive

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-service/eth"
)

type dataItem interface {
	Data() (eth.Data, error)
}

type calldataItem struct {
	calldata eth.Data
}

func (i calldataItem) Data() (eth.Data, error) {
	return i.calldata, nil
}

type blobItem struct {
	blobRef eth.IndexedBlobHash
	blob    *eth.Blob
}

func (i blobItem) Data() (eth.Data, error) {
	return i.blob.ToData()
}

func (i *blobItem) SetBlob(blobs []*eth.Blob) error {
	idx := int(i.blobRef.Index)
	if len(blobs) <= idx {
		return fmt.Errorf("blob out of range, len(blobs): %d, idx: %d", len(blobs), idx)
	}
	i.blob = blobs[idx]
	if i.blob == nil {
		return fmt.Errorf("nil blob for idx: %d", idx)
	}
	return nil
}

// BlobDataSource fetches blobs or calldata as appropriate and transforms them into usable rollup
// data.
type BlobDataSource struct {
	data         []dataItem
	ref          eth.L1BlockRef
	batcherAddr  common.Address
	dsCfg        DataSourceConfig
	fetcher      L1TransactionFetcher
	blobsFetcher L1BlobsFetcher
	log          log.Logger
}

// NewBlobDataSource creates a new blob data source.
func NewBlobDataSource(ctx context.Context, log log.Logger, dsCfg DataSourceConfig, fetcher L1TransactionFetcher, blobsFetcher L1BlobsFetcher, ref eth.L1BlockRef, batcherAddr common.Address) DataIter {
	return &BlobDataSource{
		ref:          ref,
		dsCfg:        dsCfg,
		fetcher:      fetcher,
		log:          log.New("origin", ref),
		batcherAddr:  batcherAddr,
		blobsFetcher: blobsFetcher,
	}
}

// Next returns the next piece of batcher data if any remains. It returns ResetError if it cannot
// find the referenced block or a referenced blob, or TemporaryError for any other failure to fetch
// a block or blob.
func (ds *BlobDataSource) Next(ctx context.Context) (eth.Data, error) {
	if ds.data == nil {
		var err error
		if ds.data, err = ds.open(ctx); err != nil {
			return nil, err
		}
	}

	if len(ds.data) == 0 {
		return nil, io.EOF
	}

	next := ds.data[0]
	ds.data = ds.data[1:]

	data, err := next.Data()
	if err != nil {
		ds.log.Error("ignoring data item due to error", "err", err)
		return ds.Next(ctx)
	}
	return data, nil
}

// open fetches and returns the blob or calldata (as appropriate) from all valid batcher
// transactions in the referenced block. Returns an empty (non-nil) array if no batcher
// transactions are found. Returns nil array whenever an error is returned.
func (ds *BlobDataSource) open(ctx context.Context) ([]dataItem, error) {
	_, txs, err := ds.fetcher.InfoAndTxsByHash(ctx, ds.ref.Hash)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return nil, NewResetError(fmt.Errorf("%w: failed to open blob data source", err))
		}
		return nil, NewTemporaryError(fmt.Errorf("%w: failed to open blob data source", err))
	}

	data, hashes := dataAndHashesFromTxs(txs, &ds.dsCfg, ds.batcherAddr)

	if len(hashes) == 0 {
		// there are no blobs to fetch so we can return immediately
		return data, nil
	}

	// download the actual blob bodies corresponding to the indexed blob hashes
	blobs, err := ds.blobsFetcher.GetBlobs(ctx, ds.ref, hashes)
	if errors.Is(err, ethereum.NotFound) {
		// If the L1 block was available, then the blobs should be available too. The only
		// exception is if the blob retention window has expired, which we will ultimately handle
		// by failing over to a blob archival service.
		return nil, NewResetError(fmt.Errorf("%w: failed to fetch blobs", err))
	} else if err != nil {
		return nil, NewTemporaryError(fmt.Errorf("%w: failed to fetch blobs", err))
	}

	// go back over the data array and populate the blob pointers
	if err := fillBlobPointers(data, blobs); err != nil {
		// this shouldn't happen unless there is a bug in the blobs fetcher
		return nil, NewResetError(fmt.Errorf("%w: failed to fill blob pointers", err))
	}
	return data, nil
}

// dataAndHashesFromTxs extracts calldata and datahashes from the input transactions and returns them. It
// creates a placeholder blobOrCalldata element for each returned blob hash that must be populated
// by fillBlobPointers after blob bodies are retrieved.
func dataAndHashesFromTxs(txs types.Transactions, config *DataSourceConfig, batcherAddr common.Address) ([]dataItem, []eth.IndexedBlobHash) {
	data := make([]dataItem, 0, len(txs))
	var hashes []eth.IndexedBlobHash
	blobIndex := 0 // index of each blob in the block's blob sidecar
	for _, tx := range txs {
		// skip any non-batcher transactions
		if !isValidBatchTx(tx, config.l1Signer, config.batchInboxAddress, batcherAddr) {
			blobIndex += len(tx.BlobHashes())
			continue
		}
		// handle non-blob batcher transactions by extracting their calldata
		if tx.Type() != types.BlobTxType {
			calldata := eth.Data(tx.Data())
			data = append(data, &calldataItem{calldata: calldata})
			continue
		}
		// handle blob batcher transactions by extracting their blob hashes, ignoring any calldata.
		if len(tx.Data()) > 0 {
			log.Warn("blob tx has calldata, which will be ignored", "txhash", tx.Hash())
		}
		for _, h := range tx.BlobHashes() {
			idh := eth.IndexedBlobHash{
				Index: uint64(blobIndex),
				Hash:  h,
			}
			hashes = append(hashes, idh)
			data = append(data, &blobItem{blobRef: idh}) // will fill in blob pointers after we download them below
			blobIndex++
		}
	}
	return data, hashes
}

// fillBlobPointers goes back through the data array and fills in the pointers to the fetched blob
// bodies. There should be exactly one placeholder blobOrCalldata element for each blob, otherwise
// error is returned.
func fillBlobPointers(data []dataItem, blobs []*eth.Blob) error {
	var blobCount int
	for _, di := range data {
		if bi, ok := di.(*blobItem); ok {
			blobCount++
			if err := bi.SetBlob(blobs); err != nil {
				return err
			}
		}
	}
	if blobCount != len(blobs) {
		return fmt.Errorf("too many blobs, %d != %d", len(blobs), blobCount)
	}
	return nil
}
