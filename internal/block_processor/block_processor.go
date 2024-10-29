package block_processor

import (
	"context"
	"fmt"
	"math/big"

	"github.com/aapclark/go-indexer/m/v2/internal/rpc"
)

/*
struct:
- store last processed block
- store latest block from stream
- store earliest block from stream (this is used to determine if the http or stream should be used for the log query)
- use rpc

- live (ws) query method
- historic (http) query method
- helper func to determine when to use historic query vs. live query
- health check method
*/

type BlockProcessor struct {
	CurrentBlock         *big.Int // ? store block or just number
	firstStreamedBlock   *big.Int
	LatestProcessedBlock *big.Int // ? store block or just number
	// ? blocks with errors (for retry)

	// unmarshaled abi?
	// logs channel

	RPC *rpc.RPCClient
}

func NewBlockProcessor(rpc *rpc.RPCClient, startBlock big.Int) *BlockProcessor {
	return &BlockProcessor{
		LatestProcessedBlock: &startBlock,
		RPC:                  rpc,
	}
}

// updateCurrentBlock sets the BlockProcessor instance's CurrentBlock to the provided value
// if the firstStreamedBlock field has not been set, this function will do so
func (b *BlockProcessor) updateCurrentBlock(bn *big.Int) error {
	if b.CurrentBlock.Cmp(bn) > 0 {
		return fmt.Errorf("blockprocessor.updateCurrentBlock: module's current block is %v but tried to update to %v which is lower", b.CurrentBlock, bn)
	}
	b.CurrentBlock = bn
	// handle recording of firstStreamedBlock by assigning it only if previous value is nil
	if b.firstStreamedBlock == nil {
		b.firstStreamedBlock = bn
	}
	return nil
}

/*
shouldMakeHistoricQuery returns a boolean when the module state is such that a historic query should be made to retrieve logs.
If the latest processed block is lesser than the first streamed block, an historic query will be needed.
*/
func (b *BlockProcessor) shouldMakeHistoricQuery() bool {
	return b.firstStreamedBlock.Cmp(b.LatestProcessedBlock) > 1
}

// TrackLatestBlock subscribes the BlockProcessor to the underlying RPC module's latest block number feed
func (b *BlockProcessor) TrackLatestBlock(ctx context.Context) {
	bCh := make(chan *big.Int)
	eCh := make(chan error)

	b.RPC.SubscribeToLatestBlockNumber(ctx, bCh, eCh)

	for {
		select {
		case bn := <-bCh:
			b.updateCurrentBlock(bn)
		case <-eCh:
			// TODO: handle this error in some way
			return
		case <-ctx.Done():
			return
		}
	}
}
