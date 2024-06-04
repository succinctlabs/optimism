package mipsevm

import (
	"errors"
	"io"
)

type PreimageOracle interface {
	Hint(v []byte)
	GetPreimage(k [32]byte) []byte
}

type Debug struct {
	stack  []uint32
	caller []uint32
}

type InstrumentedState struct {
	state *State

	stdOut io.Writer
	stdErr io.Writer

	lastMemAccess   uint32
	memProofEnabled bool
	memProof        [28 * 32]byte

	preimageOracle PreimageOracle

	// cached pre-image data, including 8 byte length prefix
	lastPreimage []byte
	// key for above preimage
	lastPreimageKey [32]byte
	// offset we last read from, or max uint32 if nothing is read this step
	lastPreimageOffset uint32

	meta         *Metadata
	debugEnabled bool
}

const (
	fdStdin         = 0
	fdStdout        = 1
	fdStderr        = 2
	fdHintRead      = 3
	fdHintWrite     = 4
	fdPreimageRead  = 5
	fdPreimageWrite = 6
)

// See src/syscall/zerrors_linux_mips.go
const (
	MipsEBADF  = 0x9
	MipsEINVAL = 0x16
	MipsEAGAIN = 0xb
)

func NewInstrumentedState(state *State, po PreimageOracle, stdOut, stdErr io.Writer) *InstrumentedState {
	return &InstrumentedState{
		state:          state,
		stdOut:         stdOut,
		stdErr:         stdErr,
		preimageOracle: po,
	}
}

func (m *InstrumentedState) InitDebug(meta *Metadata) error {
	if meta == nil {
		return errors.New("metadata is nil")
	}
	m.debugEnabled = true
	m.meta = meta
	return nil
}

func (m *InstrumentedState) Step(proof bool) (wit *StepWitness, err error) {
	m.memProofEnabled = proof
	m.lastMemAccess = ^uint32(0)
	m.lastPreimageOffset = ^uint32(0)

	preThreadCount := len(m.state.Threads)
	if proof {
		insnProof := m.state.Memory.MerkleProof(m.state.Threads[m.state.CurrentThread].State.PC)
		// TODO: current-thread witness
		wit = &StepWitness{
			State:    m.state.EncodeWitness(),
			MemProof: insnProof[:],
		}
	}
	err = m.mipsStep()
	if err != nil {
		return nil, err
	}

	if proof {
		wit.MemProof = append(wit.MemProof, m.memProof[:]...)
		if m.lastPreimageOffset != ^uint32(0) {
			wit.PreimageOffset = m.lastPreimageOffset
			wit.PreimageKey = m.lastPreimageKey
			wit.PreimageValue = m.lastPreimage
		}
		// if more threads than pre-instruction, add new-thread proof
		if len(m.state.Threads) > preThreadCount {
			// TODO: add proof for last merkle branch, so we can append a new thread state
		}
	}
	return
}

func (m *InstrumentedState) LastPreimage() ([32]byte, []byte, uint32) {
	return m.lastPreimageKey, m.lastPreimage, m.lastPreimageOffset
}
