// Copyright 2015 The zerium Authors
// This file is part of the zerium library.
//
// The zerium library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The zerium library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the zerium library. If not, see <http://www.gnu.org/licenses/>.

package zrm

import (
	"context"
	"math/big"

	"github.com/apolo-technologies/zerium/accounts"
	"github.com/apolo-technologies/zerium/common"
	"github.com/apolo-technologies/zerium/common/math"
	"github.com/apolo-technologies/zerium/core"
	"github.com/apolo-technologies/zerium/core/bloombits"
	"github.com/apolo-technologies/zerium/core/state"
	"github.com/apolo-technologies/zerium/core/types"
	"github.com/apolo-technologies/zerium/core/vm"
	"github.com/apolo-technologies/zerium/zrm/downloader"
	"github.com/apolo-technologies/zerium/zrm/gasprice"
	"github.com/apolo-technologies/zerium/zrmdb"
	"github.com/apolo-technologies/zerium/event"
	"github.com/apolo-technologies/zerium/params"
	"github.com/apolo-technologies/zerium/rpc"
)

// zaeapiBackend implements zaeapi.Backend for full nodes
type zaeapiBackend struct {
	zrm *Zerium
	gpo *gasprice.Oracle
}

func (b *zaeapiBackend) ChainConfig() *params.ChainConfig {
	return b.zrm.chainConfig
}

func (b *zaeapiBackend) CurrentBlock() *types.Block {
	return b.zrm.blockchain.CurrentBlock()
}

func (b *zaeapiBackend) SetHead(number uint64) {
	b.zrm.protocolManager.downloader.Cancel()
	b.zrm.blockchain.SetHead(number)
}

func (b *zaeapiBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.zrm.miner.PendingBlock()
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.zrm.blockchain.CurrentBlock().Header(), nil
	}
	return b.zrm.blockchain.GetHeaderByNumber(uint64(blockNr)), nil
}

func (b *zaeapiBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.zrm.miner.PendingBlock()
		return block, nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.zrm.blockchain.CurrentBlock(), nil
	}
	return b.zrm.blockchain.GetBlockByNumber(uint64(blockNr)), nil
}

func (b *zaeapiBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Pending state is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block, state := b.zrm.miner.Pending()
		return state, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	stateDb, err := b.zrm.BlockChain().StateAt(header.Root)
	return stateDb, header, err
}

func (b *zaeapiBackend) GetBlock(ctx context.Context, blockHash common.Hash) (*types.Block, error) {
	return b.zrm.blockchain.GetBlockByHash(blockHash), nil
}

func (b *zaeapiBackend) GetReceipts(ctx context.Context, blockHash common.Hash) (types.Receipts, error) {
	return core.GetBlockReceipts(b.zrm.chainDb, blockHash, core.GetBlockNumber(b.zrm.chainDb, blockHash)), nil
}

func (b *zaeapiBackend) GetTd(blockHash common.Hash) *big.Int {
	return b.zrm.blockchain.GetTdByHash(blockHash)
}

func (b *zaeapiBackend) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	vmError := func() error { return nil }

	context := core.NewEVMContext(msg, header, b.zrm.BlockChain(), nil)
	return vm.NewEVM(context, state, b.zrm.chainConfig, vmCfg), vmError, nil
}

func (b *zaeapiBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.zrm.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *zaeapiBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.zrm.BlockChain().SubscribeChainEvent(ch)
}

func (b *zaeapiBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.zrm.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *zaeapiBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.zrm.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *zaeapiBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.zrm.BlockChain().SubscribeLogsEvent(ch)
}

func (b *zaeapiBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.zrm.txPool.AddLocal(signedTx)
}

func (b *zaeapiBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.zrm.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *zaeapiBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.zrm.txPool.Get(hash)
}

func (b *zaeapiBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.zrm.txPool.State().GetNonce(addr), nil
}

func (b *zaeapiBackend) Stats() (pending int, queued int) {
	return b.zrm.txPool.Stats()
}

func (b *zaeapiBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.zrm.TxPool().Content()
}

func (b *zaeapiBackend) SubscribeTxPreEvent(ch chan<- core.TxPreEvent) event.Subscription {
	return b.zrm.TxPool().SubscribeTxPreEvent(ch)
}

func (b *zaeapiBackend) Downloader() *downloader.Downloader {
	return b.zrm.Downloader()
}

func (b *zaeapiBackend) ProtocolVersion() int {
	return b.zrm.EthVersion()
}

func (b *zaeapiBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *zaeapiBackend) ChainDb() zrmdb.Database {
	return b.zrm.ChainDb()
}

func (b *zaeapiBackend) EventMux() *event.TypeMux {
	return b.zrm.EventMux()
}

func (b *zaeapiBackend) AccountManager() *accounts.Manager {
	return b.zrm.AccountManager()
}

func (b *zaeapiBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.zrm.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *zaeapiBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.zrm.bloomRequests)
	}
}
