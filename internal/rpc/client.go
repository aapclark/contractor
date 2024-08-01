package rpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/aapclark/go-indexer/m/v2/internal/config"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type RpcClient struct {
	Url        string
	WsUrl      string
	wsClient   *ethclient.Client
	httpClient *ethclient.Client
}

func NewClient(cfg config.RpcConfig) (*RpcClient, error) {
	c, err := ethclient.Dial(cfg.Url)
	if err != nil {
		return nil, err
	}
	ws, err := ethclient.Dial(cfg.StreamUrl)
	if err != nil {
		return nil, err
	}

	client := &RpcClient{
		Url:        cfg.Url,
		httpClient: c,
		WsUrl:      cfg.StreamUrl,
		wsClient:   ws,
	}

	return client, nil
}

func (c RpcClient) SubscribeToLatestBlockNumber(ctx context.Context, outCh chan *big.Int, errCh chan error) {
	headerCh := make(chan *types.Header)

	s, err := c.wsClient.SubscribeNewHead(ctx, headerCh)
	if err != nil {
		errCh <- err
		return
	}

	for {
		select {
		case header := <-headerCh:
			outCh <- header.Number
		case err := <-s.Err():
			errCh <- err
			return
		case <-ctx.Done():
			fmt.Println("canceled")
			// errCh <- ctx.Err()
			c.close()
			return
		}
	}
}

func (c *RpcClient) close() {
	c.httpClient.Close()
	c.wsClient.Close()
}
