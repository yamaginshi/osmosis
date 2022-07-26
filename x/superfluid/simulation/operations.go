package simulation

import (
	"errors"

	"github.com/osmosis-labs/osmosis/v10/simulation/simtypes"
	"github.com/osmosis-labs/osmosis/v10/x/superfluid/keeper"
	"github.com/osmosis-labs/osmosis/v10/x/superfluid/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Simulation operation weights constants.
const (
	DefaultWeightMsgSuperfluidDelegate          int = 100
	DefaultWeightMsgSuperfluidUndelegate        int = 50
	DefaultWeightMsgSuperfluidRedelegate        int = 50
	DefaultWeightSetSuperfluidAssetsProposal    int = 100
	DefaultWeightRemoveSuperfluidAssetsProposal int = 2

	OpWeightMsgSuperfluidDelegate   = "op_weight_msg_superfluid_delegate"
	OpWeightMsgSuperfluidUndelegate = "op_weight_msg_superfluid_undelegate"
	OpWeightMsgSuperfluidRedelegate = "op_weight_msg_superfluid_redelegate"
)

// // WeightedOperations returns all the operations from the module with their respective weights.
// func WeightedOperations(
// 	appParams simtypes.AppParams, cdc codec.JSONCodec, ak stakingtypes.AccountKeeper,
// 	bk stakingtypes.BankKeeper, sk types.StakingKeeper, lk types.LockupKeeper, k keeper.Keeper,
// ) simulation.WeightedOperations {
// 	var (
// 		weightMsgSuperfluidDelegate   int
// 		weightMsgSuperfluidUndelegate int
// 		// weightMsgSuperfluidRedelegate int
// 	)

// 	appParams.GetOrGenerate(cdc, OpWeightMsgSuperfluidDelegate, &weightMsgSuperfluidDelegate, nil,
// 		func(_ *rand.Rand) {
// 			weightMsgSuperfluidDelegate = DefaultWeightMsgSuperfluidDelegate
// 		},
// 	)

// 	appParams.GetOrGenerate(cdc, OpWeightMsgSuperfluidUndelegate, &weightMsgSuperfluidUndelegate, nil,
// 		func(_ *rand.Rand) {
// 			weightMsgSuperfluidUndelegate = DefaultWeightMsgSuperfluidUndelegate
// 		},
// 	)

// 	// appParams.GetOrGenerate(cdc, OpWeightMsgSuperfluidRedelegate, &weightMsgSuperfluidRedelegate, nil,
// 	// 	func(_ *rand.Rand) {
// 	// 		weightMsgSuperfluidRedelegate = DefaultWeightMsgSuperfluidRedelegate
// 	// 	},
// 	// )

// 	return simulation.WeightedOperations{
// 		simulation.NewWeightedOperation(
// 			weightMsgSuperfluidDelegate,
// 			SimulateMsgSuperfluidDelegate(ak, bk, sk, lk, k),
// 		),
// 		simulation.NewWeightedOperation(
// 			weightMsgSuperfluidUndelegate,
// 			SimulateMsgSuperfluidUndelegate(ak, bk, lk, k),
// 		),
// 		// simulation.NewWeightedOperation(
// 		// 	weightMsgSuperfluidRedelegate,
// 		// 	SimulateMsgSuperfluidRedelegate(ak, bk, sk, lk, k),
// 		// ),
// 	}
// }

// SimulateMsgSuperfluidDelegate generates a MsgSuperfluidDelegate with random values.
func SimulateMsgSuperfluidDelegate(k keeper.Keeper, sim *simtypes.SimCtx, ctx sdk.Context) (*types.MsgSuperfluidDelegate, error) {
	// select random validator
	validator, err := sim.RandomValidator(ctx)
	if err != nil {
		return nil, err
	}

	// select random lockup
	lock, _, err := sim.RandomLockAndAccount(ctx)
	if err != nil {
		return nil, err
	}

	multiplier := k.GetOsmoEquivalentMultiplier(ctx, lock.Coins[0].Denom)
	if multiplier.IsZero() {
		return nil, errors.New("osmo multiplier is zero")
	}

	if !k.GetLockIdIntermediaryAccountConnection(ctx, lock.ID).Empty() {
		return nil, errors.New("lock is already used for superfluid staking")
	}

	return &types.MsgSuperfluidDelegate{
		Sender:  lock.Owner,
		LockId:  lock.ID,
		ValAddr: validator.OperatorAddress,
	}, nil
}

func SimulateMsgSuperfluidUndelegate(k keeper.Keeper, sim *simtypes.SimCtx, ctx sdk.Context) (*types.MsgSuperfluidUndelegate, error) {
	lock, simAccount, err := sim.RandomLockAndAccount(ctx)
	if err != nil {
		return nil, err
	}

	if k.GetLockIdIntermediaryAccountConnection(ctx, lock.ID).Empty() {
		return nil, errors.New("lock is not used for superfluid staking")
	}

	return &types.MsgSuperfluidUndelegate{
		Sender: simAccount.Address.String(),
		LockId: lock.ID,
	}, nil
}

// func SimulateMsgSuperfluidRedelegate(ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper, sk types.StakingKeeper, lk types.LockupKeeper, k keeper.Keeper) simtypes.Operation {
// 	return func(
// 		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
// 	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
// 		simAccount, _ := simtypes.RandomAcc(r, accs)

// 		// select random validator
// 		validator := RandomValidator(ctx, r, sk)
// 		if validator == nil {
// 			return simtypes.NoOpMsg(
// 				types.ModuleName, types.TypeMsgSuperfluidRedelegate, "No validator"), nil, nil
// 		}

// 		lock, simAccount := RandomLockAndAccount(ctx, r, lk, accs)
// 		if lock == nil {
// 			return simtypes.NoOpMsg(
// 				types.ModuleName, types.TypeMsgSuperfluidRedelegate, "Account have no period lock"), nil, nil
// 		}

// 		if k.GetLockIdIntermediaryAccountConnection(ctx, lock.ID).Empty() {
// 			return simtypes.NoOpMsg(
// 				types.ModuleName, types.TypeMsgSuperfluidRedelegate, "Lock is not used for superfluid staking"), nil, nil
// 		}

// 		msg := types.MsgSuperfluidRedelegate{
// 			Sender:     lock.Owner,
// 			LockId:     lock.ID,
// 			NewValAddr: validator.OperatorAddress,
// 		}

// 		txGen := simappparams.MakeTestEncodingConfig().TxConfig
// 		return osmosimtypes.GenAndDeliverTxWithRandFees(
// 			r, app, txGen, &msg, nil, ctx, simAccount, ak, bk, types.ModuleName)
// 	}
// }
