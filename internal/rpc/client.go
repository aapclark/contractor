package rpc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TODO: rename this as more methods are needed than just subscribing
type RPCSubscriber interface {
	SubscribeNewHead(context.Context, chan *types.Header) (ethereum.Subscription, error)
	Close()
}

type RPCClient struct {
	wsClient   RPCSubscriber
	httpClient RPCSubscriber
}

func NewRPCClient(httpClient, wsClient RPCSubscriber) (*RPCClient, error) {
	client := &RPCClient{
		httpClient: httpClient,
		wsClient:   wsClient,
	}

	return client, nil
}

func NewEthClient(url string) (*ethclient.Client, error) {
	c, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c RPCClient) SubscribeToLatestBlockNumber(ctx context.Context, outCh chan *big.Int, errCh chan error) {
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
			// TODO: determine whether to send context error over channel
			// errCh <- ctx.Err()
			c.close()
			return
		}
	}
}

func (c *RPCClient) close() {
	c.httpClient.Close()
	c.wsClient.Close()
}
