package ibci

import (
	"net"
	"os"

	abcicli "github.com/cometbft/cometbft/abci/client"
	"github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/libs/log"
	cmtsync "github.com/cometbft/cometbft/libs/sync"
)

const initialHeight = 0

type IBCOutpost struct {
	logger log.Logger

	rpcCfg   *config.RPCConfig
	p2pCfg   *config.P2PConfig
	memplCfg *config.MempoolConfig

	rpcListener net.Listener
}

func NewIBCOutpost() *IBCOutpost {
	return &IBCOutpost{
		logger:   log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
		rpcCfg:   config.DefaultRPCConfig(),
		p2pCfg:   config.DefaultP2PConfig(),
		memplCfg: config.DefaultMempoolConfig(),
	}
}

func (i *IBCOutpost) Start() error {
	if err := i.startRPC(); err != nil {
		return err
	}

	app := NewIBCIApplication()
	proxyApp := abcicli.NewLocalClient(new(cmtsync.Mutex), app)

	mp, mpReactor := createMempoolAndMempoolReactor(
		i.memplCfg,
		proxyApp,
		initialHeight, // TODO: Normally this would be pulled from some state store or genesis
		i.logger,
	)

	nodeKey := p2p.LoadOrGenNodeKey("") // Generates a node key in mem

	sw := createSwitch(
		i.p2pCfg,
		transport,
		mpReactor,
		nodeInfo,
		nodeKey,
		i.logger,
	)

	return nil
}

func (i *IBCOutpost) Stop() error {
	if err := i.rpcListener.Close(); err != nil {
		return err
	}

	return nil
}
