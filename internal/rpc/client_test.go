package rpc

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

// MockSubscription is a mock implementation of ethereum.Subscription
type MockSubscription struct {
	errCh chan error
}

// Unsubscribe mocks the ethereum.Subscription method Unsubscribe
func (m *MockSubscription) Unsubscribe() {
	close(m.errCh)
}

// Err mocks the ethereum.Subscription method Err
func (m *MockSubscription) Err() <-chan error {
	return m.errCh
}

// MockRPCSubscriber is  a mock implementation of RPCSubscriber
type MockRPCSubscriber struct {
	headers chan *types.Header
	err     error
}

// SubscribeNewHead mocks the RPCSubscriber SubscribeNewHead method
func (m *MockRPCSubscriber) SubscribeNewHead(ctx context.Context, outCh chan *types.Header) (ethereum.Subscription, error) {
	if m.err != nil {
		return nil, m.err
	}

	go func() {
		for header := range m.headers {
			outCh <- header
		}
	}()
	return &MockSubscription{
		errCh: make(chan error),
	}, nil
}

// Close mocks the RPCSubscriber method Close
func (m *MockRPCSubscriber) Close() {}

func TestSubscribeToLatestBlock(t *testing.T) {
	tests := []struct {
		name          string
		headers       []*types.Header
		expected      []*big.Int
		expectedError bool
	}{
		{
			name: "single header",
			headers: []*types.Header{{
				Number: big.NewInt(1),
			}},
			expected:      []*big.Int{big.NewInt(1)},
			expectedError: false,
		},
		{
			name: "multiple headers",
			headers: []*types.Header{
				{
					Number: big.NewInt(2),
				},
				{
					Number: big.NewInt(3),
				},
				{
					Number: big.NewInt(4),
				},
			},
			expected:      []*big.Int{big.NewInt(2), big.NewInt(3), big.NewInt(4)},
			expectedError: false,
		},
		{
			name:          "context canceled",
			headers:       []*types.Header{},
			expected:      []*big.Int{},
			expectedError: true,
		},
	}
	for _, tt := range tests {
		outCh := make(chan *big.Int, 1)
		errCh := make(chan error, 1)

		headerCh := make(chan *types.Header, len(tt.headers))
		for _, header := range tt.headers {
			headerCh <- header
		}
		close(headerCh)
		mockSubscriber := &MockRPCSubscriber{
			headers: headerCh,
		}
		client := &RPCClient{
			WSClient:   mockSubscriber,
			HTTPClient: &MockRPCSubscriber{},
		}

		ctx, cancel := context.WithCancel(context.Background())
		if tt.expectedError {
			cancel()
		}
		defer cancel()

		go client.SubscribeToLatestBlockNumber(ctx, outCh, errCh)

		var received []*big.Int
		for range tt.expected {
			select {
			case bn := <-outCh:
				received = append(received, bn)
			case err := <-errCh:
				if !tt.expectedError {
					t.Fatalf("unexpected error %v", err)
				}
				return
			case <-time.After(500 * time.Millisecond):
				t.Fatalf("test timed out after 500ms")
				return
			}
		}

		if tt.expectedError {
			select {
			case err := <-errCh:
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("received unexpected error %v", err)
				}
			case <-time.After(500 * time.Microsecond):
				t.Fatalf("test timed out after 500ms")
			}
		} else {
			for i, exp := range tt.expected {
				if exp.Cmp(received[i]) != 0 {
					t.Fatalf("expected block number %v but got %v", exp, received[i])
				}
			}
		}

	}
}
