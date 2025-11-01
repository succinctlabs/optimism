package sysgo

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/ethereum-optimism/optimism/op-devstack/devtest"
	"github.com/ethereum-optimism/optimism/op-service/logpipe"
)

// SubProcess is a process that can be started and stopped.
type SubProcess struct {
	p   devtest.P
	cmd *exec.Cmd

	stdOutLogs logpipe.LogProcessor
	stdErrLogs logpipe.LogProcessor

	mu     sync.Mutex
	onExit func(error)
}

func NewSubProcess(p devtest.P, stdOutLogs, stdErrLogs logpipe.LogProcessor) *SubProcess {
	return &SubProcess{
		p:          p,
		stdOutLogs: stdOutLogs,
		stdErrLogs: stdErrLogs,
	}
}

func (sp *SubProcess) Start(cmdPath string, args []string, env []string) error {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	if sp.cmd != nil {
		return fmt.Errorf("process is still running (PID: %d)", sp.cmd.Process.Pid)
	}
	cmd := exec.Command(cmdPath, args...)
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdout = sp.stdOutLogs
	cmd.Stderr = sp.stdErrLogs
	if err := cmd.Start(); err != nil {
		return err
	}
	sp.cmd = cmd

	go func(cmd *exec.Cmd) {
		err := cmd.Wait()
		sp.mu.Lock()
		if sp.cmd == cmd {
			sp.cmd = nil
		}
		cb := sp.onExit
		sp.mu.Unlock()

		if cb != nil {
			cb(err)
		}

		if err != nil {
			sp.p.Logger().Error("Sub-process crashed", "err", err)
		} else {
			sp.p.Logger().Info("Sub-process exited")
		}
	}(cmd)

	sp.p.Cleanup(func() {
		err := sp.Stop(true)
		if err != nil {
			sp.p.Logger().Error("Shutdown error", "err", err)
		}
	})
	return nil
}

// Stop waits for the process to stop, interrupting the process if it has not completed and
// interrupt is true.
func (sp *SubProcess) Stop(interrupt bool) error {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	if sp.cmd == nil {
		return nil // already stopped gracefully
	}

	// If not already done, then try an interrupt first as requested.
	if sp.cmd.ProcessState == nil && interrupt {
		sp.p.Logger().Info("Sending interrupt")
		if err := sp.cmd.Process.Signal(os.Interrupt); err != nil {
			return err
		}
	}

	if err := sp.cmd.Wait(); err != nil {
		sp.p.Logger().Warn("Sub-process exited with error", "err", err)
	} else {
		sp.p.Logger().Info("Sub-process gracefully exited")
	}

	sp.cmd = nil
	return nil
}

func (sp *SubProcess) OnExit(fn func(error)) {
	sp.mu.Lock()
	sp.onExit = fn
	sp.mu.Unlock()
}
