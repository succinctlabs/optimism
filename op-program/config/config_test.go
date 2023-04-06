package config

import (
	"testing"

	"github.com/ethereum-optimism/optimism/op-node/chaincfg"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfigIsValid(t *testing.T) {
	err := NewConfig(&chaincfg.Goerli).Check()
	require.NoError(t, err)
}

func TestRollupConfig(t *testing.T) {
	t.Run("Required", func(t *testing.T) {
		err := NewConfig(nil).Check()
		require.ErrorIs(t, err, ErrMissingRollupConfig)
	})

	t.Run("Valid", func(t *testing.T) {
		err := NewConfig(&rollup.Config{}).Check()
		require.ErrorIs(t, err, rollup.ErrBlockTimeZero)
	})
}

func TestFetchingArgConsistency(t *testing.T) {
	t.Run("RequireL2WhenL1Set", func(t *testing.T) {
		cfg := NewConfig(&chaincfg.Beta1)
		cfg.L1URL = "https://example.com:1234"
		require.ErrorIs(t, cfg.Check(), ErrL1AndL2Inconsistent)
	})
	t.Run("RequireL1WhenL2Set", func(t *testing.T) {
		cfg := NewConfig(&chaincfg.Beta1)
		cfg.L2URL = "https://example.com:1234"
		require.ErrorIs(t, cfg.Check(), ErrL1AndL2Inconsistent)
	})
	t.Run("AllowNeitherSet", func(t *testing.T) {
		cfg := NewConfig(&chaincfg.Beta1)
		require.NoError(t, cfg.Check())
	})
	t.Run("AllowBothSet", func(t *testing.T) {
		cfg := NewConfig(&chaincfg.Beta1)
		cfg.L1URL = "https://example.com:1234"
		cfg.L2URL = "https://example.com:4678"
		require.NoError(t, cfg.Check())
	})
}

func TestFetchingEnabled(t *testing.T) {
	t.Run("FetchingNotEnabledWhenNoFetcherUrlsSpecified", func(t *testing.T) {
		cfg := NewConfig(&chaincfg.Beta1)
		require.False(t, cfg.FetchingEnabled(), "Should not enable fetching when node URL not supplied")
	})
	t.Run("FetchingNotEnabledWhenNoL1UrlSpecified", func(t *testing.T) {
		cfg := NewConfig(&chaincfg.Beta1)
		cfg.L2URL = "https://example.com:1234"
		require.False(t, cfg.FetchingEnabled(), "Should not enable fetching when node URL not supplied")
	})
	t.Run("FetchingNotEnabledWhenNoL2UrlSpecified", func(t *testing.T) {
		cfg := NewConfig(&chaincfg.Beta1)
		cfg.L1URL = "https://example.com:1234"
		require.False(t, cfg.FetchingEnabled(), "Should not enable fetching when node URL not supplied")
	})
	t.Run("FetchingEnabledWhenBothFetcherUrlsSpecified", func(t *testing.T) {
		cfg := NewConfig(&chaincfg.Beta1)
		cfg.L1URL = "https://example.com:1234"
		cfg.L2URL = "https://example.com:5678"
		require.True(t, cfg.FetchingEnabled(), "Should enable fetching when node URL supplied")
	})
}
