package rpc

import (
	"context"
	"math/big"

	"github.com/aapclark/go-indexer/m/v2/internal/config"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RpcClient struct {
	Url    string
	client *ethclient.Client
}

func NewClient(cfg config.RpcConfig) (*RpcClient, error) {
	c, err := ethclient.Dial(cfg.Url)
	if err != nil {
		return nil, err
	}
	client := &RpcClient{
		Url:    cfg.Url,
		client: c,
	}
	return client, nil
}

func (c RpcClient) SubscribeLatestBlockNumber(ch chan *big.Int) error {
	var err error
	ctx := context.Background()
	hCh := make(chan *types.Header)

	_, err = c.client.SubscribeNewHead(ctx, hCh)
	for header := range hCh {
		ch <- header.Number
	}

	return err
}
