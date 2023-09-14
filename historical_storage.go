package ibci

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Only used for fetching historical self consensus states.
type historicalStorage struct{}

func (h *historicalStorage) GetHistoricalInfo(ctx sdk.Context, height int64) (stakingtypes.HistoricalInfo, bool) {
	// TODO: Impl historical storage.
	return stakingtypes.HistoricalInfo{}, false
}

func (h *historicalStorage) UnbondingTime(ctx sdk.Context) time.Duration {
	return 0
}
