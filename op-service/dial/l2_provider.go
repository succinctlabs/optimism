package dial

import (
	"context"

	"github.com/ethereum-optimism/optimism/op-service/sources"
	"github.com/ethereum/go-ethereum/ethclient"
)

type L2EndpointProvider interface {
	EthClient(ctx context.Context) (*ethclient.Client, error)
	RollupClient(ctx context.Context) (*sources.RollupClient, error)
}

type StaticL2EndpointProvider struct {
	ethClient    *ethclient.Client
	rollupClient *sources.RollupClient
}

func NewStaticL2EndpointProvider(ethClient *ethclient.Client, rollupClient *sources.RollupClient) *StaticL2EndpointProvider {
	return &StaticL2EndpointProvider{
		ethClient:    ethClient,
		rollupClient: rollupClient,
	}
}

func (p *StaticL2EndpointProvider) EthClient(context.Context) (*ethclient.Client, error) {
	return p.ethClient, nil
}

func (p *StaticL2EndpointProvider) RollupClient(context.Context) (*sources.RollupClient, error) {
	return p.rollupClient, nil
}
