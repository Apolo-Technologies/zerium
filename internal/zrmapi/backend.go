// Copyright 2017 The zerium Authors
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

// Package ethapi implements the general Zerium API functions.
package zrmapi

import (
	"context"
	"math/big"

	"github.com/apolo-technologies/zerium/accounts"
	"github.com/apolo-technologies/zerium/zrmcom"
	"github.com/apolo-technologies/zerium/core"
	"github.com/apolo-technologies/zerium/core/state"
	"github.com/apolo-technologies/zerium/core/types"
	"github.com/apolo-technologies/zerium/core/vm"
	"github.com/apolo-technologies/zerium/zrm/downloader"
	"github.com/apolo-technologies/zerium/zrmdb"
	"github.com/apolo-technologies/zerium/event"
	"github.com/apolo-technologies/zerium/params"
	"github.com/apolo-technologies/zerium/rpc"
)

// Backend interface provides the zrmcom API services (that are provided by
// both full and light clients) with access to necessary functions.
type Backend interface {
	// general Zerium API
	Downloader() *downloader.Downloader
	ProtocolVersion() int
	SuggestPrice(ctx context.Context) (*big.Int, error)
	ChainDb() zrmdb.Database
	EventMux() *event.TypeMux
	AccountManager() *accounts.Manager
	// BlockChain API
	SetHead(number uint64)
	HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error)
	BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error)
	StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error)
	GetBlock(ctx context.Context, blockHash zrmcom.Hash) (*types.Block, error)
	GetReceipts(ctx context.Context, blockHash zrmcom.Hash) (types.Receipts, error)
	GetTd(blockHash zrmcom.Hash) *big.Int
	GetZvm(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.Zvm, func() error, error)
	SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription
	SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription
	SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription

	// TxPool API
	SendTx(ctx context.Context, signedTx *types.Transaction) error
	GetPoolTransactions() (types.Transactions, error)
	GetPoolTransaction(txHash zrmcom.Hash) *types.Transaction
	GetPoolNonce(ctx context.Context, addr zrmcom.Address) (uint64, error)
	Stats() (pending int, queued int)
	TxPoolContent() (map[zrmcom.Address]types.Transactions, map[zrmcom.Address]types.Transactions)
	SubscribeTxPreEvent(chan<- core.TxPreEvent) event.Subscription

	ChainConfig() *params.ChainConfig
	CurrentBlock() *types.Block
}

func GetAPIs(apiBackend Backend) []rpc.API {
	nonceLock := new(AddrLocker)
	return []rpc.API{
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicZeriumAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicBlockChainAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicTransactionPoolAPI(apiBackend, nonceLock),
			Public:    true,
		}, {
			Namespace: "txpool",
			Version:   "1.0",
			Service:   NewPublicTxPoolAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(apiBackend),
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicAccountAPI(apiBackend.AccountManager()),
			Public:    true,
		}, {
			Namespace: "personal",
			Version:   "1.0",
			Service:   NewPrivateAccountAPI(apiBackend, nonceLock),
			Public:    false,
		},
	}
}
