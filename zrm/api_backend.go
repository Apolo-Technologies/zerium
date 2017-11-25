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

	"github.com/abt/zerium/accounts"
	"github.com/abt/zerium/common"
	"github.com/abt/zerium/common/math"
	"github.com/abt/zerium/core"
	"github.com/abt/zerium/core/bloombits"
	"github.com/abt/zerium/core/state"
	"github.com/abt/zerium/core/types"
	"github.com/abt/zerium/core/vm"
	"github.com/abt/zerium/zrm/downloader"
	"github.com/abt/zerium/zrm/gasprice"
	"github.com/abt/zerium/zrmdb"
	"github.com/abt/zerium/event"
	"github.com/abt/zerium/params"
	"github.com/abt/zerium/rpc"
)

// ZrmApiBackend implements zrmapi.Backend for full nodes
type ZrmApiBackend struct {
	zrm *Zerium
	gpo *gasprice.Oracle
}

func (b *ZrmApiBackend) ChainConfig() *params.ChainConfig {
	return b.zrm.chainConfig
}

func (b *ZrmApiBackend) CurrentBlock() *types.Block {
	return b.zrm.blockchain.CurrentBlock()
}

func (b *ZrmApiBackend) SetHead(number uint64) {
	b.zrm.protocolManager.downloader.Cancel()
	b.zrm.blockchain.SetHead(number)
}

func (b *ZrmApiBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
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

func (b *ZrmApiBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
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

func (b *ZrmApiBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
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

func (b *ZrmApiBackend) GetBlock(ctx context.Context, blockHash common.Hash) (*types.Block, error) {
	return b.zrm.blockchain.GetBlockByHash(blockHash), nil
}

func (b *ZrmApiBackend) GetReceipts(ctx context.Context, blockHash common.Hash) (types.Receipts, error) {
	return core.GetBlockReceipts(b.zrm.chainDb, blockHash, core.GetBlockNumber(b.zrm.chainDb, blockHash)), nil
}

func (b *ZrmApiBackend) GetTd(blockHash common.Hash) *big.Int {
	return b.zrm.blockchain.GetTdByHash(blockHash)
}

func (b *ZrmApiBackend) GetZVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.ZVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	vmError := func() error { return nil }

	context := core.NewZVMContext(msg, header, b.zrm.BlockChain(), nil)
	return vm.NewZVM(context, state, b.zrm.chainConfig, vmCfg), vmError, nil
}

func (b *ZrmApiBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.zrm.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *ZrmApiBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.zrm.BlockChain().SubscribeChainEvent(ch)
}

func (b *ZrmApiBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.zrm.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *ZrmApiBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.zrm.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *ZrmApiBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.zrm.BlockChain().SubscribeLogsEvent(ch)
}

func (b *ZrmApiBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	return b.zrm.txPool.AddLocal(signedTx)
}

func (b *ZrmApiBackend) GetPoolTransactions() (types.Transactions, error) {
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

func (b *ZrmApiBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.zrm.txPool.Get(hash)
}

func (b *ZrmApiBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.zrm.txPool.State().GetNonce(addr), nil
}

func (b *ZrmApiBackend) Stats() (pending int, queued int) {
	return b.zrm.txPool.Stats()
}

func (b *ZrmApiBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.zrm.TxPool().Content()
}

func (b *ZrmApiBackend) SubscribeTxPreEvent(ch chan<- core.TxPreEvent) event.Subscription {
	return b.zrm.TxPool().SubscribeTxPreEvent(ch)
}

func (b *ZrmApiBackend) Downloader() *downloader.Downloader {
	return b.zrm.Downloader()
}

func (b *ZrmApiBackend) ProtocolVersion() int {
	return b.zrm.EthVersion()
}

func (b *ZrmApiBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *ZrmApiBackend) ChainDb() zrmdb.Database {
	return b.zrm.ChainDb()
}

func (b *ZrmApiBackend) EventMux() *event.TypeMux {
	return b.zrm.EventMux()
}

func (b *ZrmApiBackend) AccountManager() *accounts.Manager {
	return b.zrm.AccountManager()
}

func (b *ZrmApiBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.zrm.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *ZrmApiBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.zrm.bloomRequests)
	}
}
