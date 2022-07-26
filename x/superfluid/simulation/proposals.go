package simulation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/osmosis-labs/osmosis/v10/simulation/simtypes"
	gammtypes "github.com/osmosis-labs/osmosis/v10/x/gamm/types"
	"github.com/osmosis-labs/osmosis/v10/x/superfluid/keeper"
	"github.com/osmosis-labs/osmosis/v10/x/superfluid/types"
)

// SimulateSetSuperfluidAssetsProposal generates random superfluid asset set proposal content.
func SimulateSetSuperfluidAssetsProposal(k keeper.Keeper, sim *simtypes.SimCtx, ctx sdk.Context) (*govtypes.MsgSubmitProposal, error) {
	pools, err := sim.GAMMKeeper().GetPoolsAndPoke(ctx)
	if err != nil {
		return &govtypes.MsgSubmitProposal{}, err
	}

	if len(pools) == 0 {
		return &govtypes.MsgSubmitProposal{}, err
	}

	r := sim.GetSeededRand("select random seed")
	poolIndex := r.Intn(len(pools))
	pool := pools[poolIndex]

	content := &types.SetSuperfluidAssetsProposal{
		Title:       "set superfluid assets",
		Description: "set superfluid assets description",
		Assets: []types.SuperfluidAsset{
			{
				Denom:     gammtypes.GetPoolShareDenom(pool.GetId()),
				AssetType: types.SuperfluidAssetTypeLPShare,
			},
		},
	}
	params := sim.GovKeeper().GetDepositParams(ctx)
	depositRequired := params.MinDeposit
	acc, err := sim.RandomSimAccountWithMinCoins(ctx, depositRequired)
	if err != nil {
		return &govtypes.MsgSubmitProposal{}, err
	}

	prop, err := govtypes.NewMsgSubmitProposal(content, depositRequired, acc.Address)
	if err != nil {
		return &govtypes.MsgSubmitProposal{}, err
	}

	return prop, nil
}

// func SimulateVoteYesProposal(k keeper.Keeper, sim *simtypes.SimCtx, ctx sdk.Context) (*govtypes.MsgVote, error) {
// 	val := sim.StakingKeeper().GetBondedValidatorsByPower(ctx)
// 	if len(val) == 0 {
// 		return &govtypes.MsgVote{}, nil
// 	}

// 	option := govtypes.OptionYes
// 	propID, _ := sim.GovKeeper().GetProposalID(ctx)
// 	prop, ok := sim.GovKeeper().GetProposal(ctx, propID-1)

// 	if !ok {
// 		return &govtypes.MsgVote{}, nil
// 	}
// 	sim.GovKeeper().ActivateVotingPeriod(ctx, prop)

// 	valAddr, _ := sdk.ValAddressFromBech32(val[0].OperatorAddress)
// 	accAddr, _ := sdk.AccAddressFromHex(hex.EncodeToString(valAddr.Bytes()))
// 	msg := govtypes.NewMsgVote(accAddr, propID-1, option)
// 	return msg, nil
// }

// // SimulateRemoveSuperfluidAssetsProposal generates random superfluid asset removal proposal content.
// func SimulateRemoveSuperfluidAssetsProposal(k keeper.Keeper, gk types.GammKeeper) simtypes.ContentSimulatorFn {
// 	return func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) simtypes.Content {
// 		assets := k.GetAllSuperfluidAssets(ctx)

// 		if len(assets) == 0 {
// 			return nil
// 		}

// 		assetIndex := r.Intn(len(assets))
// 		asset := assets[assetIndex]

// 		return &types.RemoveSuperfluidAssetsProposal{
// 			Title:                 "remove superfluid assets",
// 			Description:           "remove superfluid assets description",
// 			SuperfluidAssetDenoms: []string{asset.Denom},
// 		}
// 	}
// }
