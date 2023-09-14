package ibci

import (
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
)

// IBCIApplication is an abci.Application wrapper around the IBC core module.
var _ abci.Application = (*IBCIApplication)(nil)

type IBCIApplication struct {
	ibc *ibckeeper.Keeper
}

func NewIBCIApplication() *IBCIApplication {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		ibcexported.StoreKey,
		capabilitytypes.StoreKey,
	)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)
	capabilityKeeper := capabilitykeeper.NewKeeper(
		marshaler,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)
	scopedIBCKeeper := capabilityKeeper.ScopeToModule(ibcexported.ModuleName)

	ibcKeeper := ibckeeper.NewKeeper(
		marshaler,
		keys[ibcexported.StoreKey],
		paramstypes.Subspace{}, // Not used
		&historicalStorage{},   // TODO: Need ability to fetch historical headers
		&dummyUpgradeKeeper{},  // Upgrades not needed
		scopedIBCKeeper,
	)

	return &IBCIApplication{
		ibc: ibcKeeper,
	}
}

// Implemented

func (IBCIApplication) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {
	return abci.ResponseDeliverTx{Code: abci.CodeTypeOK}
}

func (IBCIApplication) CheckTx(req abci.RequestCheckTx) abci.ResponseCheckTx {
	return abci.ResponseCheckTx{Code: abci.CodeTypeOK}
}

// Not implemented or required

func (IBCIApplication) Info(req abci.RequestInfo) abci.ResponseInfo {
	return abci.ResponseInfo{}
}

func (IBCIApplication) Commit() abci.ResponseCommit {
	return abci.ResponseCommit{}
}

func (IBCIApplication) Query(req abci.RequestQuery) abci.ResponseQuery {
	return abci.ResponseQuery{}
}

func (IBCIApplication) InitChain(req abci.RequestInitChain) abci.ResponseInitChain {
	return abci.ResponseInitChain{}
}

func (IBCIApplication) BeginBlock(req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return abci.ResponseBeginBlock{}
}

func (IBCIApplication) EndBlock(req abci.RequestEndBlock) abci.ResponseEndBlock {
	return abci.ResponseEndBlock{}
}

func (IBCIApplication) ListSnapshots(req abci.RequestListSnapshots) abci.ResponseListSnapshots {
	return abci.ResponseListSnapshots{}
}

func (IBCIApplication) OfferSnapshot(req abci.RequestOfferSnapshot) abci.ResponseOfferSnapshot {
	return abci.ResponseOfferSnapshot{}
}

func (IBCIApplication) LoadSnapshotChunk(req abci.RequestLoadSnapshotChunk) abci.ResponseLoadSnapshotChunk {
	return abci.ResponseLoadSnapshotChunk{}
}

func (IBCIApplication) ApplySnapshotChunk(req abci.RequestApplySnapshotChunk) abci.ResponseApplySnapshotChunk {
	return abci.ResponseApplySnapshotChunk{}
}

func (IBCIApplication) PrepareProposal(req abci.RequestPrepareProposal) abci.ResponsePrepareProposal {
	return abci.ResponsePrepareProposal{}
}

func (IBCIApplication) ProcessProposal(req abci.RequestProcessProposal) abci.ResponseProcessProposal {
	return abci.ResponseProcessProposal{}
}
