package fetcher

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-challenger/game/keccak/types"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

const (
	bytesPerFieldElement       = uint64(32)
	usableBytesPerFieldElement = bytesPerFieldElement - 1
	fieldElementsPerBlob       = uint64(4096)
	leafSize                   = uint64(types.BlockSize + len(common.Hash{}))
	leavesPerBlob              = usableBytesPerFieldElement * fieldElementsPerBlob / leafSize
)

var (
	ErrInvalidBlobCount = errors.New("invalid blob count")
	ErrInvalidLeafIndex = errors.New("invalid leaf index")
)

type MagicBlobThang struct {
	leafCount uint64
	blobs     []*eth.Blob
}

func Encode(leaves []types.Leaf) []*eth.Blob {
	blobs := make([]*eth.Blob, 0, uint64(len(leaves))/leavesPerBlob+1)
	for idx := uint64(0); idx < uint64(len(leaves)); idx += leavesPerBlob {
		end := min(uint64(len(leaves)), idx+leavesPerBlob)
		blobs = append(blobs, encodeToBlob(leaves[idx:end]))
	}
	return blobs
}

func encodeToBlob(leaves []types.Leaf) *eth.Blob {
	if uint64(len(leaves)) > leavesPerBlob {
		panic("too many leaves")
	}
	data := make([]byte, 0, uint64(len(leaves))*leafSize)
	for _, leaf := range leaves {
		data = append(data, leaf.Input[:]...)
		data = append(data, leaf.StateCommitment[:]...)
	}
	dataIdx := 0
	var blob eth.Blob
	for fieldIdx := uint64(0); fieldIdx < fieldElementsPerBlob; fieldIdx++ {
		elementData := data[dataIdx:min(len(data), dataIdx+int(usableBytesPerFieldElement))]
		copy(blob[fieldIdx*bytesPerFieldElement+1:], elementData)
		dataIdx += len(elementData)
		if dataIdx >= len(data) {
			break
		}
	}
	return &blob
}

func NewMagicBlobThang(metadata types.LargePreimageMetaData, blobs []*eth.Blob) (*MagicBlobThang, error) {
	size := metadata.ClaimedSize
	leafCount := uint64(size / types.BlockSize)
	if size%types.BlockSize == 0 {
		// The input data fully fills the leaves so padding gets added as an additional leaf
		leafCount++
	}
	if uint64(len(blobs)) != leafCount/leavesPerBlob+1 {
		return nil, fmt.Errorf("%w expeted %v but was %v", ErrInvalidBlobCount, leafCount/leavesPerBlob+1, len(blobs))
	}
	return &MagicBlobThang{
		leafCount: leafCount,
		blobs:     blobs,
	}, nil
}

func (b *MagicBlobThang) StateCommitments() []common.Hash {
	commitments := make([]common.Hash, b.leafCount)
	for leafIdx := range commitments {
		_, commitment := b.readLeaf(uint64(leafIdx))
		commitments[leafIdx] = commitment
	}
	return commitments
}

func (b *MagicBlobThang) Reader() io.Reader {
	readers := make([]io.Reader, b.leafCount)
	for leafIdx := range readers {
		input, _ := b.readLeaf(uint64(leafIdx))
		readers[leafIdx] = bytes.NewReader(input[:])
	}
	return io.MultiReader(readers...)
}

func (b *MagicBlobThang) DataForLeaf(idx uint64) types.LeafProofData {
	startBlobIdx, startElementIdx, _ := LeafStart(idx)
	endBlobIdx, endElementIdx, _ := LeafEnd(idx)
	if startBlobIdx != endBlobIdx {
		panic("leafs should not span multiple blobs")
	}
	commitment, err := b.blobs[startBlobIdx].ComputeKZGCommitment()
	if err != nil {
		panic(err)
	}
	blob := b.blobs[startBlobIdx].KZGBlob()
	points := make([]kzg4844.Point, b.leafCount)
	proofs := make([]kzg4844.Proof, b.leafCount)
	claims := make([]kzg4844.Claim, b.leafCount)
	for elementIdx := startElementIdx; elementIdx <= endElementIdx; elementIdx++ {
		var point kzg4844.Point
		new(big.Int).SetUint64(elementIdx).FillBytes(point[:])
		kzgProof, claim, err := kzg4844.ComputeProof(*blob, point)
		if err != nil {
			panic(err)
		}
		points[elementIdx] = point
		proofs[elementIdx] = kzgProof
		claims[elementIdx] = claim
	}
	return types.LeafProofData{
		Commitment: commitment,
		Points:     points,
		Proofs:     proofs,
		Claims:     claims,
	}
}

func (b *MagicBlobThang) Leaf(idx uint64) (types.Leaf, error) {
	if idx > b.leafCount {
		return types.Leaf{}, ErrInvalidLeafIndex
	}
	input, commitment := b.readLeaf(idx)
	return types.Leaf{
		Input:           input,
		Index:           idx,
		StateCommitment: commitment,
	}, nil
}

func (b *MagicBlobThang) readLeaf(leafIdx uint64) ([types.BlockSize]byte, common.Hash) {
	blobIdx, elementIdx, offset := LeafStart(leafIdx)
	leafData := b.ReadFrom(blobIdx, elementIdx, offset, leafSize)
	if uint64(len(leafData)) != leafSize {
		panic(fmt.Errorf("read incorrect leaf data length expected %v but was %v", leafSize, len(leafData)))
	}
	input := [types.BlockSize]byte(leafData)
	commitment := common.BytesToHash(leafData[types.BlockSize:])
	return input, commitment
}

func (b *MagicBlobThang) ReadFrom(blobIdx uint64, elementIdx uint64, offset uint64, length uint64) []byte {
	data := make([]byte, 0, length)
	for uint64(len(data)) < length {
		blob := b.blobs[blobIdx]
		element := blob[elementIdx*bytesPerFieldElement : (elementIdx+1)*bytesPerFieldElement]
		section := element[offset:min(uint64(len(element)), offset+length-uint64(len(data)))]
		data = append(data, section...)
		offset = 1
		elementIdx++
		if elementIdx >= fieldElementsPerBlob {
			elementIdx = 0
			blobIdx++
		}
	}
	return data
}

func LeafStart(leafIdx uint64) (blobIdx uint64, fieldElementIdx uint64, offset uint64) {
	blobIdx = leafIdx / leavesPerBlob
	leafInBlob := leafIdx % leavesPerBlob
	firstByteInBlob := leafInBlob * leafSize
	fieldElementIdx = firstByteInBlob / usableBytesPerFieldElement
	firstByteInField := firstByteInBlob % usableBytesPerFieldElement
	offset = firstByteInField + 1 // Account for first byte of field element being padded with a 0
	return
}

func LeafEnd(leafIdx uint64) (blobIdx uint64, fieldElementIdx uint64, offset uint64) {
	blobIdx = leafIdx / leavesPerBlob
	leafInBlob := leafIdx % leavesPerBlob
	firstByteInBlob := (leafInBlob + 1) * leafSize
	fieldElementIdx = firstByteInBlob / usableBytesPerFieldElement
	firstByteInField := firstByteInBlob % usableBytesPerFieldElement
	offset = firstByteInField
	return
}
