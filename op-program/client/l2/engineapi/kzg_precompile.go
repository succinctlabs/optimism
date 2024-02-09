// This file contains code of the upstream go-ethereum kzgPointEvaluation implementation.
// Modifications have been made, primarily to substitute kzg4844.VerifyProof with a preimage oracle call.
//
// Original copyright disclaimer, applicable only to this file:
// -------------------------------------------------------------------
// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package engineapi

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"

	preimage "github.com/ethereum-optimism/optimism/op-preimage"
	"github.com/ethereum-optimism/optimism/op-program/client/l1"
	"github.com/ethereum-optimism/optimism/op-service/eth"
)

// TODO: change the below function in op-geth to allow overriding of precompiles:
// func (evm *EVM) precompile(addr common.Address) (PrecompiledContract, bool)
// TODO: and then conditionally insert it into the execution-engine:
// this execution-engine is used outside of the op-program too,
// and we do not always want to use this precompile substitute.

// OracleKZGPointEvaluation implements the EIP-4844 point evaluation precompile,
// using the preimage-oracle to perform the evaluation.
type OracleKZGPointEvaluation struct {
	oracle preimage.Oracle
	hint   preimage.Hinter
}

// RequiredGas estimates the gas required for running the point evaluation precompile.
func (b *OracleKZGPointEvaluation) RequiredGas(input []byte) uint64 {
	return params.BlobTxPointEvaluationPrecompileGas
}

const (
	blobVerifyInputLength     = 192 // Max input length for the point evaluation precompile.
	blobPrecompileReturnValue = "000000000000000000000000000000000000000000000000000000000000100073eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001"
)

var (
	errBlobVerifyInvalidInputLength = errors.New("invalid input length")
	errBlobVerifyMismatchedVersion  = errors.New("mismatched versioned hash")
	errBlobVerifyKZGProof           = errors.New("error verifying kzg proof")
)

// Run executes the point evaluation precompile.
func (b *OracleKZGPointEvaluation) Run(input []byte) ([]byte, error) {
	// Modification note: the L1 precompile behavior may change, but not in incompatible ways.
	// We want to enforce the subset that represents the EVM behavior activated in L2.
	// Below is a copy of the Cancun behavior. L1 might expand on that at a later point.

	if len(input) != blobVerifyInputLength {
		return nil, errBlobVerifyInvalidInputLength
	}
	// versioned hash: first 32 bytes
	var versionedHash common.Hash
	copy(versionedHash[:], input[:])

	var (
		point kzg4844.Point
		claim kzg4844.Claim
	)
	// Evaluation point: next 32 bytes
	copy(point[:], input[32:])
	// Expected output: next 32 bytes
	copy(claim[:], input[64:])

	// input kzg point: next 48 bytes
	var commitment kzg4844.Commitment
	copy(commitment[:], input[96:])
	if eth.KZGToVersionedHash(commitment) != versionedHash {
		return nil, errBlobVerifyMismatchedVersion
	}

	// Proof: next 48 bytes
	var proof kzg4844.Proof
	copy(proof[:], input[144:])

	// Modification note: below replaces the kzg4844.VerifyProof call
	// ------------------------------------------------------------------
	// Now the custom OP-Stack Fault-proof part:
	// 1) emit all the data we need to perform the point-evaluation-call
	b.hint.Hint(l1.KZGPointEvaluationHint(input))

	// 2) commit to all the input data
	key := preimage.KZGPointEvaluationKey(crypto.Keccak256Hash(input[:]))

	// 3) get back a 1 (valid) or 0 (invalid)
	result := b.oracle.Get(key)
	// anything else unexpected is simply invalid oracle behavior
	if len(result) != 1 || result[0] > 1 {
		panic(fmt.Errorf("unexpected preimage oracle KZGPointEvaluation behavior, got result: %x", result))
	}
	// 4) check the result
	if result[0] == 0 {
		return nil, fmt.Errorf("%w: invalid KZG point evaluation", errBlobVerifyKZGProof)
	}
	// ------------------------------------------------------------------

	return common.Hex2Bytes(blobPrecompileReturnValue), nil
}
