package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

func main() {
	var (
		datadir   string
		jwt       string
		chainid   uint64
		genesis   string
		portsFile string
	)
	fmt.Printf("Got argss: %v\n", os.Args)
	flag.StringVar(&datadir, "datadir", "", "Temp directory to store data in")
	flag.StringVar(&jwt, "jwt", "", "File to read jwt auth token from")
	flag.Uint64Var(&chainid, "chainid", 901, "Chain ID")
	flag.StringVar(&genesis, "genesis", "", "Genesis file")
	flag.StringVar(&portsFile, "ports", "", "File to write port information to")
	flag.Parse()

	fmt.Printf("Running Erigon with datadir: %v jwt: %v chainid: %v genesis: %v ports: %v\n",
		datadir, jwt, chainid, genesis, portsFile)
	// TODO: Probably better ways of doing this...
	binpath := filepath.Join(filepath.Dir(os.Args[0]), "..", "erigon")

	cmd := exec.Command(
		binpath,
		"--datadir", datadir,
		"init", genesis,
	)
	sess, err := gexec.Start(
		cmd,
		os.Stdout,
		os.Stderr,
	)
	gomega.RegisterFailHandler(func(message string, callerSkip ...int) {
		panic(message)
	})
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	gomega.Eventually(sess.Err, time.Minute).Should(gbytes.Say("Successfully wrote genesis state"))

	cmd = exec.Command(
		binpath,
		"--chain", "dev",
		"--datadir", datadir,
		"--log.console.verbosity", "dbug",
		"--ws",
		"--mine",
		"--miner.gaslimit", "0",
		"--http.port", "0",
		"--http.addr", "127.0.0.1",
		"--http.api", "eth,debug,net,engine,erigon,web3",
		"--private.api.addr=127.0.0.1:0",
		"--allow-insecure-unlock",
		"--authrpc.addr=127.0.0.1",
		"--nat", "none",
		"--p2p.allowed-ports", "0",
		"--authrpc.port=0",
		"--authrpc.vhosts=*",
		"--authrpc.jwtsecret", jwt,
		"--networkid", strconv.FormatUint(chainid, 10),
		"--torrent.port", "0", // There doesn't seem to be an obvious way to disable torrent listening
	)
	sess, err = gexec.Start(
		cmd,
		os.Stdout,
		os.Stderr,
	)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	var enginePort, httpPort uint
	gomega.Eventually(sess.Err, time.Minute).Should(gbytes.Say("HTTP endpoint opened for Engine API\\s*url=127.0.0.1:"))
	fmt.Fscanf(sess.Err, "%d", &enginePort)
	gomega.Eventually(sess.Err, time.Minute).Should(gbytes.Say("HTTP endpoint opened\\s*url=127.0.0.1:"))
	fmt.Fscanf(sess.Err, "%d", &httpPort)
	gomega.Eventually(sess.Err, time.Minute).Should(gbytes.Say("\\[1/15 Snapshots\\] DONE"))

	data, err := json.Marshal(ports{
		HTTPEndpoint:     httpPort,
		WSEndpoint:       httpPort,
		HTTPAuthEndpoint: enginePort,
		WSAuthEndpoint:   enginePort,
	})
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	err = os.WriteFile(portsFile, data, 0o644)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	sess.Wait(30 * time.Minute)
}

type ports struct {
	HTTPEndpoint     uint `json:"HTTPEndpoint"`
	WSEndpoint       uint `json:"WSEndpoint"`
	HTTPAuthEndpoint uint `json:"HTTPAuthEndpoint"`
	WSAuthEndpoint   uint `json:"WSAuthEndpoint"`
}
