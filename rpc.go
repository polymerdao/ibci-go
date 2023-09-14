package ibci

import (
	"fmt"
	"net/http"

	abci "github.com/cometbft/cometbft/abci/types"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	rpc "github.com/cometbft/cometbft/rpc/jsonrpc/server"
	rpcserver "github.com/cometbft/cometbft/rpc/jsonrpc/server"
	rpctypes "github.com/cometbft/cometbft/rpc/jsonrpc/types"
	"github.com/cometbft/cometbft/types"
)

func (i *IBCOutpost) startRPC() error {
	config := rpcserver.DefaultConfig()

	rpcLogger := i.logger.With("module", "rpc-server")
	mux := http.NewServeMux()
	rpcserver.RegisterRPCFuncs(mux, i.routes(), rpcLogger)
	listener, err := rpcserver.Listen(
		i.rpcCfg.ListenAddress,
		config,
	)
	if err != nil {
		return err
	}

	go func() {
		if err := rpcserver.Serve(
			listener,
			mux,
			rpcLogger,
			config,
		); err != nil {
			panic(err)
		}
	}()

	i.rpcListener = listener

	return nil
}

func (i *IBCOutpost) routes() map[string]*rpc.RPCFunc {
	return map[string]*rpc.RPCFunc{
		"block_results":     rpc.NewRPCFunc(i.BlockResults, "height", rpc.Cacheable("height")),
		"broadcast_tx_sync": rpc.NewRPCFunc(i.BroadcastTxSync, "tx"),
	}
}

// BroadcastTxSync returns with the response from CheckTx. Does not wait for
// the transaction result.
// More: https://docs.cometbft.com/main/rpc/#/Tx/broadcast_tx_sync
func (i *IBCOutpost) BroadcastTxSync(ctx *rpctypes.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	resCh := make(chan *abci.ResponseCheckTx, 1)
	err := env.Mempool.CheckTx(tx, func(res *abci.ResponseCheckTx) {
		select {
		case <-ctx.Context().Done():
		case resCh <- res:
		}
	}, mempl.TxInfo{})
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Context().Done():
		return nil, fmt.Errorf("broadcast confirmation not received: %w", ctx.Context().Err())
	case res := <-resCh:
		return &ctypes.ResultBroadcastTx{
			Code:      res.Code,
			Data:      res.Data,
			Log:       res.Log,
			Codespace: res.Codespace,
			Hash:      tx.Hash(),
		}, nil
	}
}

func (i *IBCOutpost) BlockResults(_ *rpctypes.Context, heightPtr *int64) (*ctypes.ResultBlockResults, error) {
	height, err := env.getHeight(env.BlockStore.Height(), heightPtr)
	if err != nil {
		return nil, err
	}

	results, err := env.StateStore.LoadFinalizeBlockResponse(height)
	if err != nil {
		return nil, err
	}

	return &ctypes.ResultBlockResults{
		Height:                height,
		TxsResults:            results.TxResults,
		FinalizeBlockEvents:   results.Events,
		ValidatorUpdates:      results.ValidatorUpdates,
		ConsensusParamUpdates: results.ConsensusParamUpdates,
	}, nil
}
