package safedb

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/cockroachdb/pebble"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrInvalidEntry = errors.New("invalid db entry")
)

const (
	// Keys are prefixed with a constant byte to allow us to differentiate different "columns" within the data
	keyPrefixSafeByL1BlockNum byte = 0
)

var (
	SafeByL1BlockNumKey = uint64Key{prefix: keyPrefixSafeByL1BlockNum}
)

type uint64Key struct {
	prefix byte
}

func (c uint64Key) Of(num uint64) []byte {
	key := make([]byte, 0, 9)
	key = append(key, c.prefix)
	key = binary.BigEndian.AppendUint64(key, num)
	return key
}
func (c uint64Key) Max() []byte {
	return c.Of(math.MaxUint64)
}

func (c uint64Key) IterRange() *pebble.IterOptions {
	return &pebble.IterOptions{
		LowerBound: c.Of(0),
		UpperBound: c.Max(),
	}
}

type SafeDB struct {
	// m ensures all read iterators are closed before closing the database by preventing concurrent read and write
	// operations (with close considered a write operation).
	m   sync.RWMutex
	log log.Logger
	db  *pebble.DB

	writeOpts *pebble.WriteOptions

	closed bool
}

func ValueL1BlockNum(l1Hash common.Hash, l2Hash common.Hash, l2BlockNum uint64) []byte {
	val := make([]byte, 0, 72)
	val = append(val, l1Hash.Bytes()...)
	val = append(val, l2Hash.Bytes()...)
	val = binary.BigEndian.AppendUint64(val, l2BlockNum)
	return val
}

func DecodeValueL1BlockNum(val []byte) (l1Hash common.Hash, l2 eth.BlockID, err error) {
	if len(val) != 72 {
		err = ErrInvalidEntry
		return
	}
	copy(l1Hash[:], val[:32])
	copy(l2.Hash[:], val[32:64])
	l2.Number = binary.BigEndian.Uint64(val[64:])
	return
}

func NewSafeDB(logger log.Logger, path string) (*SafeDB, error) {
	db, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		return nil, err
	}
	return &SafeDB{
		log:       logger,
		db:        db,
		writeOpts: &pebble.WriteOptions{Sync: true},
	}, nil
}

func (d *SafeDB) SafeHeadUpdated(safeHead eth.L2BlockRef, l1Head eth.BlockID) error {
	d.m.Lock()
	defer d.m.Unlock()
	d.log.Info("Update safe head", "l2", safeHead.ID(), "l1", l1Head)
	batch := d.db.NewBatch()
	defer batch.Close()
	// Delete any entries after this L1 block. Normally the l1Head continuously increases and this does nothing
	// However when the pipeline resets the L1 head may drop back and we need to remove later entries and allow them
	// to be repopulated as derivation progresses again. The resulting data may be different if L1 reorged.
	if err := batch.DeleteRange(SafeByL1BlockNumKey.Of(l1Head.Number+1), SafeByL1BlockNumKey.Max(), d.writeOpts); err != nil {
		return fmt.Errorf("failed to truncate safe head entries: %w", err)
	}
	if err := batch.Set(SafeByL1BlockNumKey.Of(l1Head.Number), ValueL1BlockNum(l1Head.Hash, safeHead.Hash, safeHead.Number), d.writeOpts); err != nil {
		// TODO(client-pod#593): Add tests to ensure we don't lose data here
		// We do in fact lose this update here. Even if we didn't the correct behaviour is to retry the exact same write
		// so maybe we should just keep retrying here instead of returning an error?
		return fmt.Errorf("failed to record safe head update: %w", err)
	}
	if err := batch.Commit(d.writeOpts); err != nil {
		return fmt.Errorf("failed to commit safe head update: %w", err)
	}
	return nil
}

func (d *SafeDB) SafeHeadAtL1(ctx context.Context, l1BlockNum uint64) (l1Hash common.Hash, safeHead eth.BlockID, err error) {
	d.m.RLock()
	defer d.m.RUnlock()
	iter, err := d.db.NewIterWithContext(ctx, SafeByL1BlockNumKey.IterRange())
	if err != nil {
		return
	}
	defer iter.Close()
	if valid := iter.SeekLT(SafeByL1BlockNumKey.Of(l1BlockNum + 1)); !valid {
		err = ErrNotFound
		return
	}
	// Found an entry at or before the requested L1 block
	val, err := iter.ValueAndErr()
	if err != nil {
		return
	}
	l1Hash, safeHead, err = DecodeValueL1BlockNum(val)
	return
}

func (d *SafeDB) Close() error {
	d.m.Lock()
	defer d.m.Unlock()
	if d.closed {
		// Already closed
		return nil
	}
	d.closed = true
	return d.db.Close()
}
