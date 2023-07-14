package op_e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ethereum-optimism/optimism/op-e2e/e2eutils"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"
	"github.com/onsi/gomega/gexec"
)

type ExternalClient struct {
	Name    string
	BinPath string
	DataDir string
	JWTPath string
	ChainID uint64
	Genesis *core.Genesis

	cmd   *exec.Cmd
	ports ports
}

type ports struct {
	HTTPEndpoint     uint `json:"HTTPEndpoint"`
	WSEndpoint       uint `json:"WSEndpoint"`
	HTTPAuthEndpoint uint `json:"HTTPAuthEndpoint"`
	WSAuthEndpoint   uint `json:"WSAuthEndpoint"`
}

func (ec *ExternalClient) Run() error {
	genesisPath := filepath.Join(ec.DataDir, "genesis.json")
	o, err := os.Create(genesisPath)
	if err != nil {
		return fmt.Errorf("create genesis file: %w", err)
	}
	if err = json.NewEncoder(o).Encode(ec.Genesis); err != nil {
		return fmt.Errorf("write genesis file: %w", err)
	}

	portsFile := filepath.Join(ec.DataDir, "ports.json")
	cmd := exec.Command(
		ec.BinPath,
		"--datadir", ec.DataDir,
		"--jwt", ec.JWTPath,
		"--chainid", strconv.FormatUint(ec.ChainID, 10),
		"--genesis", genesisPath,
		"--ports", portsFile,
	)
	ec.cmd = cmd
	cmd.Stdout = gexec.NewPrefixedWriter(ec.Name, os.Stdout)
	cmd.Stderr = gexec.NewPrefixedWriter(ec.Name, os.Stderr)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting command: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = e2eutils.WaitFor(ctx, time.Second, func() (bool, error) {
		data, err := os.ReadFile(portsFile)
		if err != nil {
			log.Warn("Ports file not available", "file", portsFile, "err", err)
			return false, nil // Retry later
		}
		err = json.Unmarshal(data, &ec.ports)
		if err != nil {
			log.Warn("Parse ports file", "file", portsFile, "err", err)
			return false, nil // Retry - it may have been partially written
		}
		return true, nil
	})
	if err != nil {
		return fmt.Errorf("get port info: %w", err)
	}
	log.Info("Started external client", "bin", ec.BinPath, "ports", ec.ports)
	return nil
}

func (ec *ExternalClient) Close() {
	err := ec.cmd.Process.Kill()
	if err != nil {
		panic(err) // TODO: Be better...
	}
	_ = ec.cmd.Wait() // Probably should have a timeout and stuff...
}

func (ec *ExternalClient) WSEndpoint() string {
	return fmt.Sprintf("http://127.0.0.1:%d", ec.ports.WSEndpoint)
}

func (ec *ExternalClient) HTTPEndpoint() string {
	return fmt.Sprintf("http://127.0.0.1:%d", ec.ports.HTTPEndpoint)
}

func (ec *ExternalClient) WSAuthEndpoint() string {
	return fmt.Sprintf("http://127.0.0.1:%d", ec.ports.WSAuthEndpoint)
}

func (ec *ExternalClient) HTTPAuthEndpoint() string {
	return fmt.Sprintf("http://127.0.0.1:%d", ec.ports.HTTPAuthEndpoint)
}

var _ EthInstance = (*ExternalClient)(nil)
