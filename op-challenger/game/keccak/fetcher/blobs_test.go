package fetcher

import (
	"fmt"
	"testing"

	"github.com/ethereum-optimism/optimism/op-challenger/game/keccak/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/trie/testutil"
	"github.com/stretchr/testify/require"
)

func TestLeafStart(t *testing.T) {
	tests := []struct {
		leafIdx         uint64
		expectedBlob    uint64
		expectedElement uint64
		expectedOffset  uint64
	}{
		{0, 0, 0, 1},
		{1, 0, 5, 14},
		{2, 0, 10, 27},
		{754, 0, 4086, 7},

		{755, 1, 0, 1},
		{755 + 1, 1, 5, 14},
		{755 + 2, 1, 10, 27},
		{755 + 754, 1, 4086, 7},

		{2 * 755, 2, 0, 1},
		{2*755 + 1, 2, 5, 14},
		{2*755 + 2, 2, 10, 27},
		{2*755 + 754, 2, 4086, 7},

		{200 * 755, 200, 0, 1},
		{200*755 + 1, 200, 5, 14},
		{200*755 + 2, 200, 10, 27},
		{200*755 + 754, 200, 4086, 7},
	}
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("leaf-%v", test.leafIdx), func(t *testing.T) {
			blob, element, offset := LeafStart(test.leafIdx)
			require.Equal(t, test.expectedBlob, blob, "incorrect blob")
			require.Equal(t, test.expectedElement, element, "incorrect field element")
			require.Equal(t, test.expectedOffset, offset, "incorrect offset")
		})
	}
}

func TestLeafEnd(t *testing.T) {
	tests := []struct {
		leafIdx         uint64
		expectedBlob    uint64
		expectedElement uint64
		expectedOffset  uint64
	}{
		{0, 0, 5, 13},
		{1, 0, 10, 26},
		{753, 0, 4086, 6},
		{754, 0, 4091, 19},

		{755, 1, 5, 13},
		{755 + 1, 1, 10, 26},
		{755 + 753, 1, 4086, 6},
		{755 + 754, 1, 4091, 19},

		{2 * 755, 2, 5, 13},
		{2*755 + 1, 2, 10, 26},
		{2*755 + 753, 2, 4086, 6},
		{2*755 + 754, 2, 4091, 19},

		{200 * 755, 200, 5, 13},
		{200*755 + 1, 200, 10, 26},
		{200*755 + 753, 200, 4086, 6},
		{200*755 + 754, 200, 4091, 19},
	}
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("leaf-%v", test.leafIdx), func(t *testing.T) {
			blob, element, offset := LeafEnd(test.leafIdx)
			require.Equal(t, test.expectedBlob, blob, "incorrect blob")
			require.Equal(t, test.expectedElement, element, "incorrect field element")
			require.Equal(t, test.expectedOffset, offset, "incorrect offset")
		})
	}
}

func TestLeaf(t *testing.T) {
	leaves := make([]types.Leaf, leavesPerBlob+10)
	for i := range leaves {
		leaves[i] = types.Leaf{
			Input:           [types.BlockSize]byte(testutil.RandBytes(types.BlockSize)),
			Index:           uint64(i),
			StateCommitment: common.Hash{0xaa, byte(i), 0xcc},
		}
	}
	blobs := Encode(leaves)
	metadata := types.LargePreimageMetaData{
		ClaimedSize: uint32(leafSize) * uint32(len(leaves)),
	}
	thang, err := NewMagicBlobThang(metadata, blobs)
	require.NoError(t, err)
	for i, expected := range leaves {
		actual, err := thang.Leaf(uint64(i))
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	}
}
