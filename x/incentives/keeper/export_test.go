package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/osmosis-labs/osmosis/v11/x/incentives/types"
)

// AddGaugeRefByKey appends the provided gauge ID into an array associated with the provided key.
func (k Keeper) AddGaugeRefByKey(ctx sdk.Context, key []byte, gaugeID uint64) error {
	return k.addGaugeRefByKey(ctx, key, gaugeID)
}

func (k Keeper) DeleteGaugeRefByKey(ctx sdk.Context, key []byte, guageID uint64) error {
	return k.deleteGaugeRefByKey(ctx, key, guageID)
}

func (k Keeper) GetGaugeRefs(ctx sdk.Context, key []byte) []uint64 {
	return k.getGaugeRefs(ctx, key)
}

func (k Keeper) GetAllGaugeIDsByDenom(ctx sdk.Context, denom string) []uint64 {
	return k.getAllGaugeIDsByDenom(ctx, denom)
}

func (k Keeper) MoveUpcomingGaugeToActiveGauge(ctx sdk.Context, gauge types.Gauge) error {
	return k.moveUpcomingGaugeToActiveGauge(ctx, gauge)
}

func (k Keeper) MoveActiveGaugeToFinishedGauge(ctx sdk.Context, gauge types.Gauge) error {
	return k.moveActiveGaugeToFinishedGauge(ctx, gauge)
}

// ChargeFeeIfSufficientFeeDenomBalance see chargeFeeIfSufficientFeeDenomBalance spec.
func (k Keeper) ChargeFeeIfSufficientFeeDenomBalance(ctx sdk.Context, address sdk.AccAddress, fee sdk.Int, gaugeCoins sdk.Coins) error {
	return k.chargeFeeIfSufficientFeeDenomBalance(ctx, address, fee, gaugeCoins)
}

func (k Keeper) GetTimeKey(timestamp time.Time) []byte {
	return getTimeKey(timestamp)
}

func (k Keeper) CombineKeys(keys ...[]byte) []byte {
	return combineKeys(keys...)
}
