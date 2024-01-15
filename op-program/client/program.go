package client

import (
	"errors"
	"io"
	"os"

	"github.com/ethereum/go-ethereum/log"

	preimage "github.com/ethereum-optimism/optimism/op-preimage"
	cldr "github.com/ethereum-optimism/optimism/op-program/client/driver"
	oppio "github.com/ethereum-optimism/optimism/op-program/io"
)

// Main executes the client program in a detached context and exits the current process.
// The client runtime environment must be preset before calling this function.
func Main(logger log.Logger) {
	log.Info("Starting fault proof program client")
	preimageOracle := CreatePreimageChannel()
	preimageHinter := CreateHinterChannel()
	if err := RunProgram(logger, preimageOracle, preimageHinter); errors.Is(err, cldr.ErrClaimNotValid) {
		log.Error("Claim is invalid", "err", err)
		os.Exit(1)
	} else if err != nil {
		log.Error("Program failed", "err", err)
		os.Exit(2)
	} else {
		log.Info("Claim successfully verified")
		os.Exit(0)
	}
}

// RunProgram executes the Program, while attached to an IO based pre-image oracle, to be served by a host.
func RunProgram(logger log.Logger, preimageOracle io.ReadWriter, preimageHinter io.ReadWriter) error {
	pClient := preimage.NewOracleClient(preimageOracle)
	hClient := preimage.NewHintWriter(preimageHinter)

	bootInfo := NewBootstrapClient(pClient).BootInfo()
	game := bootInfo.Boot(pClient)
	return game.Run(logger, pClient, hClient)
}

func CreateHinterChannel() oppio.FileChannel {
	r := os.NewFile(HClientRFd, "preimage-hint-read")
	w := os.NewFile(HClientWFd, "preimage-hint-write")
	return oppio.NewReadWritePair(r, w)
}

// CreatePreimageChannel returns a FileChannel for the preimage oracle in a detached context
func CreatePreimageChannel() oppio.FileChannel {
	r := os.NewFile(PClientRFd, "preimage-oracle-read")
	w := os.NewFile(PClientWFd, "preimage-oracle-write")
	return oppio.NewReadWritePair(r, w)
}
