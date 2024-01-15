package ssz

import (
	"encoding/binary"
	"fmt"

	preimage "github.com/ethereum-optimism/optimism/op-preimage"
	"github.com/ethereum-optimism/optimism/op-service/eth"
)

type PairPreimageFn func(root eth.Bytes32) (left, right eth.Bytes32)

func (fn PairPreimageFn) TraverseBranch(root eth.Bytes32, depth uint8, index uint8) eth.Bytes32 {
	for depth > 0 {
		left, right := fn(root)
		if index&(1<<(depth-1)) == 0 {
			root = left
		} else {
			root = right
		}
		depth -= 1
	}
	return root
}

func ListLeaves[L ~[32]byte](fn PairPreimageFn, root eth.Bytes32, depth uint8) []L {
	subRoot, lengthMixin := fn(root)
	if ([24]byte)(lengthMixin[:24]) != ([24]byte{}) {
		panic(fmt.Errorf("expected uint64 length, but got %x", lengthMixin))
	}
	length := binary.LittleEndian.Uint64(lengthMixin[24:])
	if capacity := uint64(1) << depth; length > capacity {
		panic(fmt.Errorf("length %d larger than capacity %d", length, capacity))
	}
	result := make([]L, length)
	collectSubtree[L](fn, subRoot, depth, result)
	return result
}

// collectSubtree recursively traverses down the binary merkle tree, and puts the leaf values into the dest slice.
func collectSubtree[L ~[32]byte](fn PairPreimageFn, root eth.Bytes32, depth uint8, dest []L) {
	if len(dest) == 0 { // no need to expand unused sub-trees
		return
	}
	if depth == 0 {
		if len(dest) != 1 {
			panic("expected single node")
		}
		dest[0] = L(root)
		return
	}
	left, right := fn(root)
	half := uint64(1) << depth
	leftDest := dest
	rightDest := []L(nil)
	if half < uint64(len(dest)) { // split destination if it spans the left and right sub-trees
		leftDest = leftDest[:half]
		rightDest = rightDest[half:]
	}
	collectSubtree[L](fn, left, depth-1, leftDest)
	collectSubtree[L](fn, right, depth-1, rightDest)
}

func BytesList(fn PairPreimageFn, root eth.Bytes32, depth uint8) []byte {
	subRoot, lengthMixin := fn(root)
	if ([24]byte)(lengthMixin[:24]) != ([24]byte{}) {
		panic(fmt.Errorf("expected uint64 length, but got %x", lengthMixin))
	}
	byteLength := binary.LittleEndian.Uint64(lengthMixin[24:])
	length := byteLength / 32
	if capacity := uint64(1) << depth; length > capacity {
		panic(fmt.Errorf("length %d larger than capacity %d", length, capacity))
	}
	result := make([]byte, byteLength)
	collectBytesSubtree(fn, subRoot, depth, result)
	return result
}

// collectBytesSubtree recursively traverses down the binary merkle tree,
// and puts the leaf bytes data into the dest slice.
func collectBytesSubtree(fn PairPreimageFn, root eth.Bytes32, depth uint8, dest []byte) {
	if len(dest) == 0 { // no need to expand unused sub-trees
		return
	}
	if depth == 0 {
		if len(dest) > 32 {
			panic("expected no more than one node")
		}
		copy(dest, root[:])
		return
	}
	left, right := fn(root)
	half := uint64(32) << depth
	leftDest := dest
	rightDest := []byte(nil)
	if half < uint64(len(dest)) { // split destination if it spans the left and right sub-trees
		leftDest = leftDest[:half]
		rightDest = rightDest[half:]
	}
	collectBytesSubtree(fn, left, depth-1, leftDest)
	collectBytesSubtree(fn, right, depth-1, rightDest)
}

func OracleSSZ(p *preimage.OracleClient) PairPreimageFn {
	return func(root eth.Bytes32) (left, right eth.Bytes32) {
		data := p.Get(preimage.Sha256Key(root))
		if len(data) != 64 {
			panic(fmt.Errorf("expected root %s to have 64 byte SHA2-256 preimage, but got %x", root, data))
		}
		copy(left[:], data[:32])
		copy(right[:], data[32:])
		return
	}
}
