package ibci

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

type dummyUpgradeKeeper struct{}

func (d *dummyUpgradeKeeper) ClearIBCState(ctx sdk.Context, lastHeight int64) {}

func (d *dummyUpgradeKeeper) GetUpgradePlan(ctx sdk.Context) (plan upgradetypes.Plan, havePlan bool) {
	return upgradetypes.Plan{}, false
}

func (d *dummyUpgradeKeeper) GetUpgradedClient(ctx sdk.Context, height int64) ([]byte, bool) {
	return nil, false
}

func (d *dummyUpgradeKeeper) SetUpgradedClient(ctx sdk.Context, planHeight int64, bz []byte) error {
	return nil
}

func (d *dummyUpgradeKeeper) GetUpgradedConsensusState(ctx sdk.Context, lastHeight int64) ([]byte, bool) {
	return nil, false
}

func (d *dummyUpgradeKeeper) SetUpgradedConsensusState(ctx sdk.Context, planHeight int64, bz []byte) error {
	return nil
}

func (d *dummyUpgradeKeeper) ScheduleUpgrade(ctx sdk.Context, plan upgradetypes.Plan) error {
	return nil
}
