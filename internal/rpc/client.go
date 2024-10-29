package rpc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TODO: add health checks
// TODO: rename this as more methods are needed than just subscribing
type RPCSubscriber interface {
	SubscribeNewHead(context.Context, chan *types.Header) (ethereum.Subscription, error)
	SubscribeFilterLogs(context.Context, ethereum.FilterQuery, chan<- types.Log) (ethereum.Subscription, error)
	Close()
}

type RPCClient struct {
	WSClient   RPCSubscriber
	HTTPClient RPCSubscriber
}

// NewRPCClient returns an instance of NewRPCClient with the provided http and websocket rpc interfaces
func NewRPCClient(httpClient, wsClient RPCSubscriber) *RPCClient {
	client := &RPCClient{
		HTTPClient: httpClient,
		WSClient:   wsClient,
	}

	return client
}

// NewEthClient returns an instance of ethclient.Client connected to the given URL
func NewEthClient(url string) (*ethclient.Client, error) {
	c, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return c, nil
}

/*
SubscribeToLatestBlockNumber subscribes the client to new block data and sends the block number to a channel.
Any errors returned by the underlying subscription are send over an error channel
*/
func (c *RPCClient) SubscribeToLatestBlockNumber(ctx context.Context, outCh chan *big.Int, errCh chan error) {
	headerCh := make(chan *types.Header)
	defer close(headerCh)

	s, err := c.WSClient.SubscribeNewHead(ctx, headerCh)
	if err != nil {
		errCh <- err
		return
	}
	defer s.Unsubscribe()

	for {
		select {
		case header := <-headerCh:
			outCh <- header.Number
		case err := <-s.Err():
			errCh <- err
			return
		case <-ctx.Done():
			errCh <- ctx.Err()
			c.close()
			return
		}
	}
}

func (c *RPCClient) SubscribeToFilteredLogs(ctx context.Context, query ethereum.FilterQuery, outCh chan types.Log, errCh chan error) {
	logsCh := make(chan types.Log)
	s, err := c.WSClient.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		errCh <- err
		return
	}
	defer s.Unsubscribe()

	for {
		select {
		case logs := <-logsCh:
			outCh <- logs
		case err := <-s.Err():
			errCh <- err
			return
		case <-ctx.Done():
			errCh <- ctx.Err()
			c.close()
			return
		}
	}
}

// close closes HTTPClient and WSClient connections
func (c *RPCClient) close() {
	c.HTTPClient.Close()
	c.WSClient.Close()
}
