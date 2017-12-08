// Copyright 2016 The zerium Authors
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

// Package les implements the Light Zerium Subprotocol.
package les

import (
	"fmt"
	"sync"
	"time"

	"github.com/apolo-technologies/zerium/accounts"
	"github.com/apolo-technologies/zerium/common"
	"github.com/apolo-technologies/zerium/common/hexutil"
	"github.com/apolo-technologies/zerium/consensus"
	"github.com/apolo-technologies/zerium/core"
	"github.com/apolo-technologies/zerium/core/bloombits"
	"github.com/apolo-technologies/zerium/core/types"
	"github.com/apolo-technologies/zerium/zrm"
	"github.com/apolo-technologies/zerium/zrm/downloader"
	"github.com/apolo-technologies/zerium/zrm/filters"
	"github.com/apolo-technologies/zerium/zrm/gasprice"
	"github.com/apolo-technologies/zerium/zrmdb"
	"github.com/apolo-technologies/zerium/event"
	"github.com/apolo-technologies/zerium/my/zrmapi"
	"github.com/apolo-technologies/zerium/light"
	"github.com/apolo-technologies/zerium/log"
	"github.com/apolo-technologies/zerium/node"
	"github.com/apolo-technologies/zerium/p2p"
	"github.com/apolo-technologies/zerium/p2p/discv5"
	"github.com/apolo-technologies/zerium/params"
	rpc "github.com/apolo-technologies/zerium/rpc"
)

type LightZerium struct {
	odr         *LesOdr
	relay       *LesTxRelay
	chainConfig *params.ChainConfig
	// Channel for shutting down the service
	shutdownChan chan bool
	// Handlers
	peers           *peerSet
	txPool          *light.TxPool
	blockchain      *light.LightChain
	protocolManager *ProtocolManager
	serverPool      *serverPool
	reqDist         *requestDistributor
	retriever       *retrieveManager
	// DB interfaces
	chainDb zrmdb.Database // Block chain database

	bloomRequests                              chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer, chtIndexer, bloomTrieIndexer *core.ChainIndexer

	ApiBackend *LesApiBackend

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	networkId     uint64
	netRPCService *zrmapi.PublicNetAPI

	wg sync.WaitGroup
}

func New(ctx *node.ServiceContext, config *zrm.Config) (*LightZerium, error) {
	chainDb, err := zrm.CreateDB(ctx, config, "lightchaindata")
	if err != nil {
		return nil, err
	}
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)
	if _, isCompat := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !isCompat {
		return nil, genesisErr
	}
	log.Info("[LES]: Initialised chain configuration", "config", chainConfig)

	peers := newPeerSet()
	quitSync := make(chan struct{})

	lzrm := &LightZerium{
		chainConfig:      chainConfig,
		chainDb:          chainDb,
		eventMux:         ctx.EventMux,
		peers:            peers,
		reqDist:          newRequestDistributor(peers, quitSync),
		accountManager:   ctx.AccountManager,
		engine:           zrm.CreateConsensusEngine(ctx, config, chainConfig, chainDb),
		shutdownChan:     make(chan bool),
		networkId:        config.NetworkId,
		bloomRequests:    make(chan chan *bloombits.Retrieval),
		bloomIndexer:     zrm.NewBloomIndexer(chainDb, light.BloomTrieFrequency),
		chtIndexer:       light.NewChtIndexer(chainDb, true),
		bloomTrieIndexer: light.NewBloomTrieIndexer(chainDb, true),
	}

	lzrm.relay = NewLesTxRelay(peers, lzrm.reqDist)
	lzrm.serverPool = newServerPool(chainDb, quitSync, &lzrm.wg)
	lzrm.retriever = newRetrieveManager(peers, lzrm.reqDist, lzrm.serverPool)
	lzrm.odr = NewLesOdr(chainDb, lzrm.chtIndexer, lzrm.bloomTrieIndexer, lzrm.bloomIndexer, lzrm.retriever)
	if lzrm.blockchain, err = light.NewLightChain(lzrm.odr, lzrm.chainConfig, lzrm.engine); err != nil {
		return nil, err
	}
	lzrm.bloomIndexer.Start(lzrm.blockchain)
	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		lzrm.blockchain.SetHead(compat.RewindTo)
		core.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}

	lzrm.txPool = light.NewTxPool(lzrm.chainConfig, lzrm.blockchain, lzrm.relay)
	if lzrm.protocolManager, err = NewProtocolManager(lzrm.chainConfig, true, ClientProtocolVersions, config.NetworkId, lzrm.eventMux, lzrm.engine, lzrm.peers, lzrm.blockchain, nil, chainDb, lzrm.odr, lzrm.relay, quitSync, &lzrm.wg); err != nil {
		return nil, err
	}
	lzrm.ApiBackend = &LesApiBackend{lzrm, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.GasPrice
	}
	lzrm.ApiBackend.gpo = gasprice.NewOracle(lzrm.ApiBackend, gpoParams)
	return lzrm, nil
}

func lesTopic(genesisHash common.Hash, protocolVersion uint) discv5.Topic {
	var name string
	switch protocolVersion {
	case lpv1:
		name = "LES"
	case lpv2:
		name = "LES2"
	default:
		panic(nil)
	}
	return discv5.Topic(name + "@" + common.Bytes2Hex(genesisHash.Bytes()[0:8]))
}

type LightDummyAPI struct{}

// Zeriumbase is the address that mining rewards will be send to
func (s *LightDummyAPI) Zeriumbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("not supported")
}

// Coinbase is the address that mining rewards will be send to (alias for Zeriumbase)
func (s *LightDummyAPI) Coinbase() (common.Address, error) {
	return common.Address{}, fmt.Errorf("not supported")
}

// Hashrate returns the POW hashrate
func (s *LightDummyAPI) Hashrate() hexutil.Uint {
	return 0
}

// Mining returns an indication if this node is currently mining.
func (s *LightDummyAPI) Mining() bool {
	return false
}

// APIs returns the collection of RPC services the abt package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *LightZerium) APIs() []rpc.API {
	return append(zrmapi.GetAPIs(s.ApiBackend), []rpc.API{
		{
			Namespace: "zrm",
			Version:   "1.0",
			Service:   &LightDummyAPI{},
			Public:    true,
		}, {
			Namespace: "zrm",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "zrm",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.ApiBackend, true),
			Public:    true,
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)
}

func (s *LightZerium) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *LightZerium) BlockChain() *light.LightChain      { return s.blockchain }
func (s *LightZerium) TxPool() *light.TxPool              { return s.txPool }
func (s *LightZerium) Engine() consensus.Engine           { return s.engine }
func (s *LightZerium) LesVersion() int                    { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *LightZerium) Downloader() *downloader.Downloader { return s.protocolManager.downloader }
func (s *LightZerium) EventMux() *event.TypeMux           { return s.eventMux }

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *LightZerium) Protocols() []p2p.Protocol {
	return s.protocolManager.SubProtocols
}

// Start implements node.Service, starting all internal goroutines needed by the
// Zerium protocol implementation.
func (s *LightZerium) Start(srvr *p2p.Server) error {
	s.startBloomHandlers()
	log.Warn("Light client mode is an experimental feature")
	s.netRPCService = zrmapi.NewPublicNetAPI(srvr, s.networkId)
	// search the topic belonging to the oldest supported protocol because
	// servers always advertise all supported protocols
	protocolVersion := ClientProtocolVersions[len(ClientProtocolVersions)-1]
	s.serverPool.start(srvr, lesTopic(s.blockchain.Genesis().Hash(), protocolVersion))
	s.protocolManager.Start()
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Zerium protocol.
func (s *LightZerium) Stop() error {
	s.odr.Stop()
	if s.bloomIndexer != nil {
		s.bloomIndexer.Close()
	}
	if s.chtIndexer != nil {
		s.chtIndexer.Close()
	}
	if s.bloomTrieIndexer != nil {
		s.bloomTrieIndexer.Close()
	}
	s.blockchain.Stop()
	s.protocolManager.Stop()
	s.txPool.Stop()

	s.eventMux.Stop()

	time.Sleep(time.Millisecond * 200)
	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}
