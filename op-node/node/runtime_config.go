package node

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum-optimism/optimism/op-node/p2p"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

var (
	// UnsafeBlockSignerAddressSystemConfigStorageSlot is the storage slot identifier of the unsafeBlockSigner
	// `address` storage value in the SystemConfig L1 contract. Computed as `keccak256("systemconfig.unsafeblocksigner")`
	UnsafeBlockSignerAddressSystemConfigStorageSlot = common.HexToHash("0x65a7ed542fb37fe237fdfbdd70b31598523fe5b32879e307bae27a0bd9581c08")

	// RequiredProtocolVersionStorageSlot is the storage slot that the required protocol version is stored at.
	// Computed as: `bytes32(uint256(keccak256("protocolversion.required")) - 1)`
	RequiredProtocolVersionStorageSlot = common.HexToHash("0x4aaefe95bd84fd3f32700cf3b7566bc944b73138e41958b5785826df2aecace0")

	// RecommendedProtocolVersionStorageSlot is the storage slot that the recommended protocol version is stored at.
	// Computed as: `bytes32(uint256(keccak256("protocolversion.recommended")) - 1)`
	RecommendedProtocolVersionStorageSlot = common.HexToHash("0xe314dfc40f0025322aacc0ba8ef420b62fb3b702cf01e0cdf3d829117ac2ff1a")
)

type RuntimeCfgL1Source interface {
	ReadStorageAt(ctx context.Context, address common.Address, storageSlot common.Hash, blockHash common.Hash) (common.Hash, error)
}

type ReadonlyRuntimeConfig interface {
	P2PSequencerAddress(eth.L2BlockRef) common.Address
	RequiredProtocolVersion() params.ProtocolVersion
	RecommendedProtocolVersion() params.ProtocolVersion
}

// RuntimeConfig maintains runtime-configurable options.
// These options are loaded based on initial loading + updates for every subsequent L1 block.
// Only the *latest* values are maintained however, the runtime config has no concept of chain history,
// does not require any archive data, and may be out of sync with the rollup derivation process.
type RuntimeConfig struct {
	log log.Logger

	l1Client  RuntimeCfgL1Source
	rollupCfg *rollup.Config

	sc systemConfigData
	pv protocolVersionData
}

type systemConfigData struct {
	mu                    sync.RWMutex
	p2pBlockSignerSafeLag uint64

	// preP2PBlockSignerAddr and L1Ref records the previous P2P block signer address and when it was set in effect, it could be empty during the first load.
	preP2PBlockSignerAddr  common.Address
	preP2PBlockSignerL1Ref eth.L1BlockRef
	// curP2PBlockSignerAddr and L1Ref records the current P2P block signer address and when it was set in effect.
	curP2PBlockSignerAddr  common.Address
	curP2PBlockSignerL1Ref eth.L1BlockRef
}

func (sc *systemConfigData) P2PBlockSignerAddr(l2Ref eth.L2BlockRef) common.Address {
	diff := l2Ref.L1Origin.Number - sc.curP2PBlockSignerL1Ref.Number
	if diff >= sc.p2pBlockSignerSafeLag {
		return sc.curP2PBlockSignerAddr
	}

	return sc.preP2PBlockSignerAddr
}

type protocolVersionData struct {
	mu    sync.RWMutex
	l1Ref eth.L1BlockRef

	// superchain protocol version signals
	recommended params.ProtocolVersion
	required    params.ProtocolVersion
}

var _ p2p.GossipRuntimeConfig = (*RuntimeConfig)(nil)

func NewRuntimeConfig(log log.Logger, cfg *Config, l1Client RuntimeCfgL1Source, rollupCfg *rollup.Config) *RuntimeConfig {
	return &RuntimeConfig{
		log:       log,
		l1Client:  l1Client,
		rollupCfg: rollupCfg,
		sc: systemConfigData{
			p2pBlockSignerSafeLag: cfg.P2PBlockSignerAddrSafeLag,
		},
		pv: protocolVersionData{},
	}
}

func (r *RuntimeConfig) P2PSequencerAddress(l2Ref eth.L2BlockRef) common.Address {
	r.sc.mu.RLock()
	defer r.sc.mu.RUnlock()
	return r.sc.P2PBlockSignerAddr(l2Ref)
}

func (r *RuntimeConfig) RequiredProtocolVersion() params.ProtocolVersion {
	r.pv.mu.RLock()
	defer r.pv.mu.RUnlock()
	return r.pv.required
}

func (r *RuntimeConfig) RecommendedProtocolVersion() params.ProtocolVersion {
	r.pv.mu.RLock()
	defer r.pv.mu.RUnlock()
	return r.pv.recommended
}

// Load resets the runtime configuration by fetching the latest config data from L1 at the given L1 block.
// Load is safe to call concurrently, but will lock the runtime configuration modifications only,
// and will thus not block other Load calls with possibly alternative L1 block views.
func (r *RuntimeConfig) Load(ctx context.Context, l1Ref eth.L1BlockRef) error {
	if err := r.loadSystemConfig(ctx, l1Ref); err != nil {
		return err
	}
	if err := r.loadProtocolVersions(ctx, l1Ref); err != nil {
		return err
	}
	return nil
}

func (r *RuntimeConfig) OnP2PBlockSignerAddressUpdated(addr common.Address, l1Ref eth.L1BlockRef) {
	r.sc.mu.Lock()
	defer r.sc.mu.Unlock()
	if l1Ref.Time <= r.sc.curP2PBlockSignerL1Ref.Time {
		r.log.Warn("ignoring outdated P2P signer address update", "current", r.sc.curP2PBlockSignerAddr, "new", addr, "current_l1_ref", r.sc.curP2PBlockSignerL1Ref, "new_l1_ref", l1Ref)
		return
	}

	r.sc.preP2PBlockSignerAddr = r.sc.curP2PBlockSignerAddr
	r.sc.preP2PBlockSignerL1Ref = r.sc.curP2PBlockSignerL1Ref
	r.sc.curP2PBlockSignerAddr = addr
	r.sc.curP2PBlockSignerL1Ref = l1Ref
}

func (r *RuntimeConfig) loadProtocolVersions(ctx context.Context, l1Ref eth.L1BlockRef) error {
	// The superchain protocol version data is optional; only applicable to rollup configs that specify a ProtocolVersions address.
	var requiredProtVersion, recommendedProtoVersion params.ProtocolVersion
	if r.rollupCfg.ProtocolVersionsAddress != (common.Address{}) {
		requiredVal, err := r.l1Client.ReadStorageAt(ctx, r.rollupCfg.ProtocolVersionsAddress, RequiredProtocolVersionStorageSlot, l1Ref.Hash)
		if err != nil {
			return fmt.Errorf("required-protocol-version value failed to load from L1 contract: %w", err)
		}
		requiredProtVersion = params.ProtocolVersion(requiredVal)
		recommendedVal, err := r.l1Client.ReadStorageAt(ctx, r.rollupCfg.ProtocolVersionsAddress, RecommendedProtocolVersionStorageSlot, l1Ref.Hash)
		if err != nil {
			return fmt.Errorf("recommended-protocol-version value failed to load from L1 contract: %w", err)
		}
		recommendedProtoVersion = params.ProtocolVersion(recommendedVal)
	}

	r.pv.mu.Lock()
	defer r.pv.mu.Unlock()
	r.pv.l1Ref = l1Ref
	r.pv.required = requiredProtVersion
	r.pv.recommended = recommendedProtoVersion
	return nil
}

func (r *RuntimeConfig) loadSystemConfig(ctx context.Context, l1Ref eth.L1BlockRef) error {
	p2pSignerVal, err := r.l1Client.ReadStorageAt(ctx, r.rollupCfg.L1SystemConfigAddress, UnsafeBlockSignerAddressSystemConfigStorageSlot, l1Ref.Hash)
	if err != nil {
		return fmt.Errorf("failed to fetch unsafe block signing address from system config: %w", err)
	}

	r.sc.mu.Lock()
	defer r.sc.mu.Unlock()
	r.sc.curP2PBlockSignerAddr = common.BytesToAddress(p2pSignerVal[:])
	r.sc.curP2PBlockSignerL1Ref = l1Ref
	r.log.Info("loaded new runtime system config values!", "p2p_seq_address", r.sc.curP2PBlockSignerAddr)
	return nil
}
