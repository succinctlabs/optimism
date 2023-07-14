package op_e2e

import (
	"flag"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/log"
)

var enableParallelTesting bool = true

// Init testing to enable test flags
var _ = func() bool {
	testing.Init()
	return true
}()

var verboseGethNodes bool
var externalClientPath string

func init() {
	flag.BoolVar(&verboseGethNodes, "gethlogs", true, "Enable logs on geth nodes")
	flag.StringVar(&externalClientPath, "l2.client", "", "Run tests with an external L2 execution client")
	flag.Parse()
	if os.Getenv("OP_E2E_DISABLE_PARALLEL") == "true" {
		enableParallelTesting = false
	}
}

func InitParallel(t *testing.T) {
	t.Helper()
	if enableParallelTesting {
		t.Parallel()
	}
	if !verboseGethNodes {
		log.Root().SetHandler(log.DiscardHandler())
	}
}
