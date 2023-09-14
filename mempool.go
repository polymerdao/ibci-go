package ibci

import (
	cfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/mempool"
	mempoolv1 "github.com/cometbft/cometbft/mempool/v1"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/proxy"
	"github.com/cometbft/cometbft/state"
	"github.com/cometbft/cometbft/version"
)

const txIndexerStatus = "off"

func createMempoolAndMempoolReactor(
	config *cfg.MempoolConfig,
	proxyApp proxy.AppConnMempool,
	lastBlockHeight int64,
	logger log.Logger,
) (mempool.Mempool, p2p.Reactor) {
	mp := mempoolv0.NewTxMempool(
		logger,
		config,
		proxyApp,
		lastBlockHeight,
	)

	reactor := mempoolv1.NewReactor(
		config,
		mp,
	)

	return mp, reactor
}

func createSwitch(
	config *cfg.P2PConfig,
	transport p2p.Transport,
	mempoolReactor p2p.Reactor,
	nodeInfo p2p.NodeInfo,
	nodeKey *p2p.NodeKey,
	p2pLogger log.Logger,
) *p2p.Switch {
	sw := p2p.NewSwitch(
		config,
		transport,
	)
	sw.SetLogger(p2pLogger)
	sw.AddReactor("MEMPOOL", mempoolReactor)

	sw.SetNodeInfo(nodeInfo)
	sw.SetNodeKey(nodeKey)

	p2pLogger.Info("P2P Node ID", "ID", nodeKey.ID())
	return sw
}

func makeNodeInfo(
	moniker string, // human readable description of node
	p2pCfg *cfg.P2PConfig,
	rpcCfg *cfg.RPCConfig,
	nodeKey *p2p.NodeKey,
	chainID string,
) (p2p.DefaultNodeInfo, error) {
	nodeInfo := p2p.DefaultNodeInfo{
		ProtocolVersion: p2p.NewProtocolVersion(
			version.P2PProtocol, // global
			state.InitStateVersion.Consensus.Block,
			state.InitStateVersion.Consensus.App,
		),
		DefaultNodeID: nodeKey.ID(),
		Network:       chainID,
		Version:       version.TMCoreSemVer,
		Channels: []byte{
			mempool.MempoolChannel, // Only a mempool channel
		},
		Moniker: moniker,
		Other: p2p.DefaultNodeInfoOther{
			TxIndex:    txIndexerStatus,
			RPCAddress: rpcCfg.ListenAddress,
		},
	}

	lAddr := p2pCfg.ExternalAddress

	if lAddr == "" {
		lAddr = p2pCfg.ListenAddress
	}

	nodeInfo.ListenAddr = lAddr

	err := nodeInfo.Validate()
	return nodeInfo, err
}
