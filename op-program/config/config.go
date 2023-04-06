package config

import (
	"errors"

	opnode "github.com/ethereum-optimism/optimism/op-node"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/sources"
	"github.com/ethereum-optimism/optimism/op-program/flags"
	"github.com/urfave/cli"
)

var (
	ErrMissingRollupConfig = errors.New("missing rollup config")
	ErrL1AndL2Inconsistent = errors.New("l1 and l2 options must be specified together or both omitted")
)

type Config struct {
	Rollup     *rollup.Config
	L2URL      string
	L1URL      string
	L1TrustRPC bool
	L1RPCKind  sources.RPCProviderKind
}

func (c *Config) Check() error {
	if c.Rollup == nil {
		return ErrMissingRollupConfig
	}
	if err := c.Rollup.Check(); err != nil {
		return err
	}
	if (c.L1URL != "") != (c.L2URL != "") {
		return ErrL1AndL2Inconsistent
	}
	return nil
}

func (c *Config) FetchingEnabled() bool {
	return c.L1URL != "" && c.L2URL != ""
}

// NewConfig creates a Config with all optional values set to the CLI default value
func NewConfig(rollupCfg *rollup.Config) *Config {
	return &Config{
		Rollup:    rollupCfg,
		L1RPCKind: sources.RPCKindBasic,
	}
}

func NewConfigFromCLI(ctx *cli.Context) (*Config, error) {
	if err := flags.CheckRequired(ctx); err != nil {
		return nil, err
	}
	rollupCfg, err := opnode.NewRollupConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &Config{
		Rollup:     rollupCfg,
		L2URL:      ctx.GlobalString(flags.L2NodeAddr.Name),
		L1URL:      ctx.GlobalString(flags.L1NodeAddr.Name),
		L1TrustRPC: ctx.GlobalBool(flags.L1TrustRPC.Name),
		L1RPCKind:  sources.RPCProviderKind(ctx.GlobalString(flags.L1RPCProviderKind.Name)),
	}, nil
}
