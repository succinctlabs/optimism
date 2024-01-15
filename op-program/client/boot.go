package client

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"math"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	preimage "github.com/ethereum-optimism/optimism/op-preimage"
	"github.com/ethereum-optimism/optimism/op-program/chainconfig"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

const (
	L1HeadLocalIndex preimage.LocalIndexKey = iota + 1
	L2OutputRootLocalIndex
	L2ClaimLocalIndex
	L2ClaimBlockNumberLocalIndex
	L2ChainIDLocalIndex
	L1ChainIDLocalIndex

	// These local keys are only used for custom chains
	L2ChainConfigLocalIndex
	RollupConfigLocalIndex
	L1ChainConfigLocalIndex
)

// CustomChainIDIndicator is used to detect when the program should load custom chain configuration
const CustomChainIDIndicator = uint64(math.MaxUint64)

type Game interface {
	Run(logger log.Logger, pClient *preimage.OracleClient, hClient *preimage.HintWriter) error
}

type BootInfo struct {
	L1Head common.Hash

	l1ChainID uint64

	// All agreed-upon nodes between root and SPLIT_DEPTH. Zeroed hash if no agreement.
	// The first entry is the root claim itself.
	AncestorClaims []common.Hash

	// gindex at the SPLIT_DEPTH (separation between context and VM)
	SplitGindex *big.Int
}

func (bootInfo *BootInfo) Boot(r *preimage.OracleClient) Game {
	// infer SPLIT_DEPTH from the provided gindex, so we do not have to hardcode it into the program
	splitDepth := bootInfo.SplitGindex.BitLen() - 1
	// Choice between a L1 game or a L2 kind of game.
	gameChoice := bootInfo.SplitGindex.Bit(splitDepth)

	if gameChoice == 0 {
		var post common.Hash
		// TODO traverse up to find the post-state

		// TODO if the split-depth was really deep (it might get deeper with interop), we may have some padding,
		// and need to just disable part of the sub-tree, by forcing post == 0.

		// Mask out the two bits (gindex root bit, and game-choice bit), to get the index of the disputed thing
		outputIndex := new(big.Int).AndNot(bootInfo.SplitGindex, new(big.Int).Lsh(big.NewInt(0b11), uint(splitDepth)))

		l1ChainConfig := l1Config(bootInfo.l1ChainID, r)
		return &L1Game{
			L1Head:        bootInfo.L1Head,
			L1Claim:       post,
			L1ClaimNumber: outputIndex.Uint64(),
			L1ChainConfig: l1ChainConfig,
		}
	} else {
		l1SuperRoot := bootInfo.AncestorClaims[1]

		var pre, post common.Hash
		// TODO traverse up to find pre/post

		// Mask out the two bits (gindex root bit, and game-choice bit), to get the index of the disputed thing
		outputIndex := new(big.Int).AndNot(bootInfo.SplitGindex, new(big.Int).Lsh(big.NewInt(0b11), uint(splitDepth)))

		// TODO this should just be part of the disputed path (in gindex), once we have interop covering multiple L2s.
		l2ChainID := binary.BigEndian.Uint64(r.Get(L2ChainIDLocalIndex))
		l2ChainConfig, rollupConfig := l2Configs(l2ChainID, r)
		return &L2Game{
			L1Head:             bootInfo.L1Head,
			L1SuperRoot:        l1SuperRoot,
			L2Claim:            post,
			L2ClaimBlockNumber: outputIndex.Uint64(),
			L2Prestate:         pre,
			L2ChainConfig:      l2ChainConfig,
			RollupConfig:       rollupConfig,
		}
	}
}

func l1Config(l1ChainID uint64, r *preimage.OracleClient) *params.ChainConfig {
	if l1ChainID == CustomChainIDIndicator {
		var conf params.ChainConfig
		err := json.Unmarshal(r.Get(L1ChainConfigLocalIndex), &conf)
		if err != nil {
			panic("failed to bootstrap l2ChainConfig")
		}
		return &conf
	} else {
		switch l1ChainID {
		case params.GoerliChainConfig.ChainID.Uint64():
			return params.GoerliChainConfig
		case params.MainnetChainConfig.ChainID.Uint64():
			return params.MainnetChainConfig
		case params.SepoliaChainConfig.ChainID.Uint64():
			return params.SepoliaChainConfig
		case params.HoleskyChainConfig.ChainID.Uint64():
			return params.HoleskyChainConfig
		default:
			panic(fmt.Errorf("unrecognized L1 chain: %d"))
		}
	}
}

func l2Configs(l2ChainID uint64, r *preimage.OracleClient) (*params.ChainConfig, *rollup.Config) {
	var l2ChainConfig *params.ChainConfig
	var rollupConfig *rollup.Config
	if l2ChainID == CustomChainIDIndicator {
		l2ChainConfig = new(params.ChainConfig)
		err := json.Unmarshal(r.Get(L2ChainConfigLocalIndex), l2ChainConfig)
		if err != nil {
			panic("failed to bootstrap l2ChainConfig")
		}
		rollupConfig = new(rollup.Config)
		err = json.Unmarshal(r.Get(RollupConfigLocalIndex), rollupConfig)
		if err != nil {
			panic("failed to bootstrap rollup config")
		}
	} else {
		var err error
		rollupConfig, err = chainconfig.RollupConfigByChainID(l2ChainID)
		if err != nil {
			panic(err)
		}
		l2ChainConfig, err = chainconfig.ChainConfigByChainID(l2ChainID)
		if err != nil {
			panic(err)
		}
	}
	return l2ChainConfig, rollupConfig
}

type oracleClient interface {
	Get(key preimage.Key) []byte
}

type BootstrapClient struct {
	r oracleClient
}

func NewBootstrapClient(r oracleClient) *BootstrapClient {
	return &BootstrapClient{r: r}
}

func (br *BootstrapClient) BootInfo() *BootInfo {
	l1Head := common.BytesToHash(br.r.Get(L1HeadLocalIndex))
	l1ChainID := binary.BigEndian.Uint64(br.r.Get(L1ChainIDLocalIndex))
	// TODO read branch nodes between root-claim and SPLIT_DEPTH
	var ancestors []common.Hash

	// TODO read split gindex
	var splitGindex common.Hash

	return &BootInfo{
		L1Head:         l1Head,
		l1ChainID:      l1ChainID,
		AncestorClaims: ancestors,
		SplitGindex:    splitGindex.Big(),
	}
}
