package ether

import (
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
	checkJobs = 64
)

var maxSlot = common.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

type accountData struct {
	balance *big.Int
	address common.Address
}

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

	// WaitGroup to wait on each iteration job to finish.
	var wg sync.WaitGroup
	// Channel to receive storage slot keys and values from each iteration job.
	outCh := make(chan accountData)
	// Channel to receive errors from each iteration job.
	errCh := make(chan error, checkJobs)
	// Channel to cancel all iteration jobs.
	cancelCh := make(chan struct{})

	// Keep track of the total migrated supply.
	totalFound := new(big.Int)

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
			select {
			case <-cancelCh:
				// If we've already received an error, stop processing.
				return
			default:
				break
			}

			// Use the raw (i.e., secure hashed) key to check if we've reached
			// the end of the partition.
			if new(big.Int).SetBytes(it.Key).Cmp(end.Big()) >= 0 {
				return
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

			// We can safely ignore specific slots (totalSupply, name, symbol).
			if ignoredSlots[key] {
				continue
			}

			// No accounts should have a balance in state. If they do, bail.
			addr := slotsAddrs[key]
			bal := db.GetBalance(addr)
			if bal.Sign() != 0 {
				log.Error(
					"account has non-zero balance in state - should never happen",
					"addr", addr,
					"balance", bal.String(),
				)
				if !noCheck {
					errCh <- fmt.Errorf("account has non-zero balance in state - should never happen: %s", addr.String())
					return
				}
			}

			slotType, ok := slotsInp[key]
			if !ok {
				if noCheck {
					log.Error("ignoring unknown storage slot in state", "slot", key.String())
				} else {
					errCh <- fmt.Errorf("unknown storage slot in state: %s", key.String())
					return
				}
			}

			// Add balances to the total found.
			switch slotType {
			case 1:
				// Convert the value to a common.Hash, then send to the channel.
				value := common.BytesToHash(content)
				outCh <- accountData{
					balance: value.Big(),
					address: addr,
				}
			case 2:
				// Allowance slot.
				continue
			default:
				// Should never happen.
				if noCheck {
					log.Error("unknown slot type", "slot", key, "type", slotType)
				} else {
					log.Crit("unknown slot type %d, should never happen", slotType)
				}
			}
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

	// Keep track of the last error seen.
	var lastErr error

	// Only cancel things once.
	var cancelOnce sync.Once

	// Kick off another background process to collect
	// values from the channel and add them to the map.
	var count int
	progress := util.ProgressLogger(1000, "Collected OVM_ETH storage slot")
	go func() {
		defer func() {
			finCh <- struct{}{}
		}()
		for {
			select {
			case account := <-outCh:
				progress()
				// Accumulate addresses and total supply.
				addrs = append(addrs, account.address)
				totalFound = new(big.Int).Add(totalFound, account.balance)
			case err := <-errCh:
				lastErr = err
				cancelOnce.Do(func() {
					close(cancelCh)
				})
			case <-cancelCh:
				return
			}
		}
	}()

	wg.Wait()
	cancelOnce.Do(func() {
		close(cancelCh)
	})
	<-finCh

	if lastErr != nil {
		return nil, lastErr
	}

	// Log how many slots were iterated over.
	log.Info("Iterated legacy balances", "count", count)

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
