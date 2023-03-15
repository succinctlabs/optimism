package ether

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"math/big"
	"sync"

	"github.com/ethereum-optimism/optimism/op-chain-ops/crossdomain"
	"github.com/ethereum-optimism/optimism/op-chain-ops/util"

	"github.com/ethereum-optimism/optimism/op-bindings/predeploys"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
)

const (
	checkJobs = 16
)

var maxSlot = common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

// PreCheckBalances checks that the given list of addresses and allowances represents all storage
// slots in the LegacyERC20ETH contract. We don't have to filter out extra addresses like we do for
// withdrawals because we'll simply carry the balance of a given address to the new system, if the
// account is extra then it won't have any balance and nothing will happen.
func PreCheckBalances(ldb ethdb.Database, db *state.StateDB, addresses []common.Address, allowances []*crossdomain.Allowance, chainID int, noCheck bool) (FilteredOVMETHAddresses, error) {
	// Chain params to use for integrity checking.
	params := crossdomain.ParamsByChainID[chainID]
	if params == nil {
		return nil, fmt.Errorf("no chain params for %d", chainID)
	}

	return doMigration(db, addresses, allowances, params.ExpectedSupplyDelta, noCheck)
}

func doMigration(db *state.StateDB, addresses []common.Address, allowances []*crossdomain.Allowance, expDiff *big.Int, noCheck bool) (FilteredOVMETHAddresses, error) {
	// We'll need to maintain a list of all addresses that we've seen along with all of the storage
	// slots based on the witness data.
	addrs := make([]common.Address, 0)
	slotsAddrs := make(map[common.Hash]common.Address)
	slotsInp := make(map[common.Hash]int)

	// For each known address, compute its balance key and add it to the list of addresses.
	// Mint events are instrumented as regular ETH events in the witness data, so we no longer
	// need to iterate over mint events during the migration.
	for _, addr := range addresses {
		sk := CalcOVMETHStorageKey(addr)
		slotsAddrs[sk] = addr
		slotsInp[sk] = 1
	}

	// For each known allowance, compute its storage key and add it to the list of addresses.
	for _, allowance := range allowances {
		sk := CalcAllowanceStorageKey(allowance.From, allowance.To)
		slotsAddrs[sk] = allowance.From
		slotsInp[sk] = 2
	}

	// Add the old SequencerEntrypoint because someone sent it ETH a long time ago and it has a
	// balance but none of our instrumentation could easily find it. Special case.
	sequencerEntrypointAddr := common.HexToAddress("0x4200000000000000000000000000000000000005")
	entrySK := CalcOVMETHStorageKey(sequencerEntrypointAddr)
	slotsAddrs[entrySK] = sequencerEntrypointAddr
	slotsInp[entrySK] = 1

	// Build a mapping of every storage slot in the LegacyERC20ETH contract, except the list of
	// slots that we know we can ignore (totalSupply, name, symbol).
	slotsAct := make(map[common.Hash]common.Hash)

	// WaitGroup to wait on each iteration job to finish.
	var wg sync.WaitGroup
	// Channel to receive storage slot keys and values from each iteration job.
	kvCh := make(chan [2]common.Hash)

	// Divide the key space into partitions by dividing the key space by the number
	// of jobs. This will leave some slots left over, which we handle below.
	partSize := new(big.Int).Div(maxSlot.Big(), big.NewInt(checkJobs))

	// Define a worker function to iterate over each partition.
	worker := func(start, end common.Hash) {
		// Decrement the WaitGroup when the function returns.
		defer wg.Done()

		// Create a new storage trie. Each trie returned by db.StorageTrie
		// is a copy, so this is safe for concurrent use.
		st, err := db.StorageTrie(predeploys.LegacyERC20ETHAddr)
		if err != nil {
			// Should never happen, so explode if it does.
			log.Crit("cannot get storage trie for LegacyERC20ETHAddr", "err", err)
		}

		it := trie.NewIterator(st.NodeIterator(start.Bytes()))

		// Below code is largely based on db.ForEachStorage. We can't use that
		// because it doesn't allow us to specify a start and end key.
		for it.Next() {
			// Use the raw (i.e., secure hashed) key to check if we've reached
			// the end of the partition.
			if new(big.Int).SetBytes(it.Key).Cmp(end.Big()) >= 0 {
				break
			}

			// Skip if the value is empty.
			rawValue := it.Value
			if len(rawValue) == 0 {
				continue
			}

			// Get the preimage.
			key := common.BytesToHash(st.GetKey(it.Key))

			// Parse the raw value.
			_, content, _, err := rlp.Split(rawValue)
			if err != nil {
				// Should never happen, so explode if it does.
				log.Crit("mal-formed data in state: %v", err)
			}

			// Convert the value to a common.Hash, then send to the channel.
			value := common.BytesToHash(content)
			kvCh <- [2]common.Hash{key, value}
		}
	}

	for i := 0; i < checkJobs; i++ {
		wg.Add(1)

		// Compute the start and end keys for this partition.
		start := common.BigToHash(new(big.Int).Mul(big.NewInt(int64(i)), partSize))
		var end common.Hash
		if i < checkJobs-1 {
			// If this is not the last partition, use the next partition's start key as the end.
			end = common.BigToHash(new(big.Int).Mul(big.NewInt(int64(i+1)), partSize))
		} else {
			// If this is the last partition, use the max slot as the end.
			end = maxSlot
		}

		// Kick off our worker.
		go worker(start, end)
	}

	// Make a channel to make sure that the collector process completes.
	finCh := make(chan struct{})

	// Kick off another background process to collect
	// values from the channel and add them to the map.
	var count int
	progress := util.ProgressLogger(1000, "Collected OVM_ETH storage slot")
	go func() {
		for kv := range kvCh {
			key := kv[0]
			value := kv[1]

			progress()

			// We can safely ignore specific slots (totalSupply, name, symbol).
			if ignoredSlots[key] {
				continue
			}

			// Slot exists, so add it to the map.
			slotsAct[key] = value
			count++
		}
		finCh <- struct{}{}
	}()

	wg.Wait()
	close(kvCh)
	<-finCh

	// Log how many slots were iterated over.
	log.Info("Iterated legacy balances", "count", count)

	// Iterate over the list of known slots and check that we have a slot for each one. We'll also
	// keep track of the total balance to be migrated and throw if the total supply exceeds the
	// expected supply delta.
	totalFound := new(big.Int)
	var unknown bool
	progress = util.ProgressLogger(1000, "Checked OVM_ETH storage slot")
	for slot := range slotsAct {
		slotType, ok := slotsInp[slot]
		if !ok {
			if noCheck {
				log.Error("ignoring unknown storage slot in state", "slot", slot.String())
			} else {
				unknown = true
				log.Error("unknown storage slot in state", "slot", slot.String())
				continue
			}
		}

		// No accounts should have a balance in state. If they do, bail.
		addr := slotsAddrs[slot]
		bal := db.GetBalance(addr)
		if bal.Sign() != 0 {
			log.Error(
				"account has non-zero balance in state - should never happen",
				"addr", addr,
				"balance", bal.String(),
			)
			if !noCheck {
				return nil, fmt.Errorf("account has non-zero balance in state - should never happen: %s", addr.String())
			}
		}

		// Add balances to the total found.
		switch slotType {
		case 1:
			// Balance slot.
			totalFound.Add(totalFound, slotsAct[slot].Big())
			addrs = append(addrs, addr)
		case 2:
			// Allowance slot.
			continue
		default:
			// Should never happen.
			if noCheck {
				log.Error("unknown slot type", "slot", slot, "type", slotType)
			} else {
				log.Crit("unknown slot type: %d", slotType)
			}
		}

		progress()
	}
	if unknown {
		return nil, errors.New("unknown storage slots in state (see logs for details)")
	}

	// Verify the supply delta. Recorded total supply in the LegacyERC20ETH contract may be higher
	// than the actual migrated amount because self-destructs will remove ETH supply in a way that
	// cannot be reflected in the contract. This is fine because self-destructs just mean the L2 is
	// actually *overcollateralized* by some tiny amount.
	totalSupply := getOVMETHTotalSupply(db)
	delta := new(big.Int).Sub(totalSupply, totalFound)
	if delta.Cmp(expDiff) != 0 {
		log.Error(
			"supply mismatch",
			"migrated", totalFound.String(),
			"supply", totalSupply.String(),
			"delta", delta.String(),
			"exp_delta", expDiff.String(),
		)
		if !noCheck {
			return nil, fmt.Errorf("supply mismatch: %s", delta.String())
		}
	}

	// Supply is verified.
	log.Info(
		"supply verified OK",
		"migrated", totalFound.String(),
		"supply", totalSupply.String(),
		"delta", delta.String(),
		"exp_delta", expDiff.String(),
	)

	// We know we have at least a superset of all addresses here since we know that we have every
	// storage slot. It's fine to have extras because they won't have any balance.
	return addrs, nil
}
