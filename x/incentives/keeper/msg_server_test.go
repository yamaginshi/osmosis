package keeper_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/osmosis-labs/osmosis/v10/x/incentives/keeper"
	"github.com/osmosis-labs/osmosis/v10/x/incentives/types"

	incentiveskeeper "github.com/osmosis-labs/osmosis/v10/x/incentives/keeper"
	lockuptypes "github.com/osmosis-labs/osmosis/v10/x/lockup/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ = suite.TestingSuite(nil)

func (suite *KeeperTestSuite) TestCreateGaugeFee() {

	tests := []struct {
		name                 string
		accountBalanceToFund sdk.Coins
		gaugeAddition        sdk.Coins
		expectedEndBalance   sdk.Coins
		isPerpetual          bool
		isModuleAccount      bool
		expectErr            bool
	}{
		{
			name:                 "user creates a non-perpetual gauge and fills gauge with all remaining tokens in their wallet",
			accountBalanceToFund: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(60000000))),
			gaugeAddition:        sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000))),
			//expectedEndBalance:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0))),
		},
		{
			name:                 "user creates a non-perpetual gauge and fills gauge with some remaining tokens in their wallet",
			accountBalanceToFund: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(70000000))),
			gaugeAddition:        sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000))),
			//expectedEndBalance:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000))),
		},
		{
			name:                 "user with multiple denoms creates a non-perpetual gauge and fills gauge with some remaining tokens in their wallet",
			accountBalanceToFund: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(70000000)), sdk.NewCoin("foo", sdk.NewInt(70000000))),
			gaugeAddition:        sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000))),
			//expectedEndBalance:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000)), sdk.NewCoin("foo", sdk.NewInt(70000000))),
		},
		{
			name:                 "module account creates a perpetual gauge and fills gauge with some remaining tokens in the module",
			accountBalanceToFund: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(70000000)), sdk.NewCoin("foo", sdk.NewInt(70000000))),
			gaugeAddition:        sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000))),
			//expectedEndBalance:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000)), sdk.NewCoin("foo", sdk.NewInt(70000000))),
			isPerpetual:     true,
			isModuleAccount: true,
		},
		{
			name:                 "user with multiple denoms creates a perpetual gauge and fills gauge with some remaining tokens in their wallet",
			accountBalanceToFund: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(70000000)), sdk.NewCoin("foo", sdk.NewInt(70000000))),
			gaugeAddition:        sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000))),
			//expectedEndBalance:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000)), sdk.NewCoin("foo", sdk.NewInt(70000000))),
			isPerpetual: true,
		},
		{
			name:                 "user tries to create a non-perpetual gauge but does not have enough funds to pay for the create gauge fee",
			accountBalanceToFund: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(40000000))),
			gaugeAddition:        sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000))),
			//expectedEndBalance:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(40000000))),
			expectErr: true,
		},
		{
			name:                 "user tries to create a non-perpetual gauge but does not have the correct fee denom",
			accountBalanceToFund: sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(60000000))),
			gaugeAddition:        sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(10000000))),
			//expectedEndBalance:   sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(60000000))),
			expectErr: true,
		},
		// TODO: This is unexpected behavior
		// We need validation to not charge fee if user doesn't have enough funds
		// {
		// 	name:                 "one user tries to create a gauge, has enough funds to pay for the create gauge fee but not enough to fill the gauge",
		// 	accountBalanceToFund: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(60000000))),
		// 	gaugeAddition:        sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(30000000))),
		// 	expectedEndBalance:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(60000000))),
		// 	expectErr:            true,
		// },

	}

	for _, tc := range tests {
		suite.SetupTest()

		testAccountPubkey := secp256k1.GenPrivKeyFromSecret([]byte("acc")).PubKey()
		testAccountAddress := sdk.AccAddress(testAccountPubkey.Address())

		ctx := suite.Ctx
		bankKeeper := suite.App.BankKeeper
		msgServer := keeper.NewMsgServerImpl(suite.App.IncentivesKeeper)

		suite.FundAcc(testAccountAddress, tc.accountBalanceToFund)

		if tc.isModuleAccount {
			modAcc := authtypes.NewModuleAccount(authtypes.NewBaseAccount(testAccountAddress, testAccountPubkey, 1, 0),
				"module",
				"permission",
			)
			suite.App.AccountKeeper.SetModuleAccount(ctx, modAcc)
		}

		suite.SetupManyLocks(1, defaultLiquidTokens, defaultLPTokens, defaultLockDuration)
		distrTo := lockuptypes.QueryCondition{
			LockQueryType: lockuptypes.ByDuration,
			Denom:         defaultLPDenom,
			Duration:      defaultLockDuration,
		}

		msg := &types.MsgCreateGauge{
			IsPerpetual:       tc.isPerpetual,
			Owner:             testAccountAddress.String(),
			DistributeTo:      distrTo,
			Coins:             tc.gaugeAddition,
			StartTime:         time.Now(),
			NumEpochsPaidOver: 1,
		}

		// System under test.

		_, err := msgServer.CreateGauge(sdk.WrapSDKContext(ctx), msg)

		if tc.expectErr {
			suite.Require().Error(err)
		} else {
			suite.Require().NoError(err)
		}

		balanceAmount := bankKeeper.GetAllBalances(suite.Ctx, testAccountAddress)
		suite.Require().Equal(tc.expectedEndBalance.String(), balanceAmount.String(), "test: %v", tc.name)

		if tc.expectErr {
			suite.Require().Equal(tc.accountBalanceToFund.String(), balanceAmount.String(), "test: %v", tc.name)
		} else {
			fee := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(incentiveskeeper.CreateGaugeFee)))
			fmt.Printf("accountBalanceToFund %v \n", tc.accountBalanceToFund.String())
			fmt.Printf("gaugeAddition %v \n", tc.gaugeAddition.String())
			accountBalance := tc.accountBalanceToFund.Sub(tc.gaugeAddition)
			finalAccountBalance := accountBalance.Sub(fee)
			fmt.Printf("accountBalance %v \n", accountBalance.String())
			fmt.Printf("fee %v \n", fee.String())
			suite.Require().Equal(finalAccountBalance.String(), balanceAmount.String(), "test: %v", tc.name)
		}

	}
}

func (suite *KeeperTestSuite) TestAddToGaugeFee() {

	tests := []struct {
		name                 string
		accountBalanceToFund sdk.Coins
		gaugeAddition        sdk.Coins
		expectedEndBalance   sdk.Coins
		nonexistentGauge     bool
		isPerpetual          bool
		isModuleAccount      bool
		expectErr            bool
	}{
		{
			name:                 "user adds to a perpetual gauge and fills gauge with all remaining tokens in their wallet",
			accountBalanceToFund: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(60000000))),
			gaugeAddition:        sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000))),
			isPerpetual:          true,
			//expectedEndBalance:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0))),
		},
		{
			name:                 "user adds to a perpetual gauge and fills gauge with some remaining tokens in their wallet",
			accountBalanceToFund: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(70000000))),
			gaugeAddition:        sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000))),
			isPerpetual:          true,
			//expectedEndBalance:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000))),
		},
		{
			name:                 "user with multiple denoms adds to a perpetual gauge and fills gauge with some remaining tokens in their wallet",
			accountBalanceToFund: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(70000000)), sdk.NewCoin("foo", sdk.NewInt(70000000))),
			gaugeAddition:        sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000))),
			isPerpetual:          true,
			//expectedEndBalance:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000)), sdk.NewCoin("foo", sdk.NewInt(70000000))),
		},
		{
			name:                 "module account adds to a perpetual gauge and fills gauge with some remaining tokens in the module account",
			accountBalanceToFund: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(70000000)), sdk.NewCoin("foo", sdk.NewInt(70000000))),
			gaugeAddition:        sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000))),
			//expectedEndBalance:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000)), sdk.NewCoin("foo", sdk.NewInt(70000000))),
			isPerpetual:     true,
			isModuleAccount: true,
		},
		{
			name:                 "user with multiple denoms adds to a non-perpetual gauge and fills gauge with some remaining tokens in their wallet",
			accountBalanceToFund: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(70000000)), sdk.NewCoin("foo", sdk.NewInt(70000000))),
			gaugeAddition:        sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000))),
			//expectedEndBalance:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000)), sdk.NewCoin("foo", sdk.NewInt(70000000))),
		},
		{
			name:                 "user tries to add to a non-perpetual gauge but does not have enough funds to pay for the add to gauge fee",
			accountBalanceToFund: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(40000000))),
			gaugeAddition:        sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000))),
			//expectedEndBalance:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(40000000))),
			expectErr: true,
		},
		{
			name:                 "user tries to add to a non-perpetual gauge but does not have the correct fee denom",
			accountBalanceToFund: sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(60000000))),
			gaugeAddition:        sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(10000000))),
			//expectedEndBalance:   sdk.NewCoins(sdk.NewCoin("foo", sdk.NewInt(60000000))),
			expectErr: true,
		},
		{
			name:                 "user tries to add to a perpetual gauge but the gauge does not exist",
			accountBalanceToFund: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(60000000))),
			gaugeAddition:        sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10000000))),
			nonexistentGauge:     true,
			isPerpetual:          true,
			expectErr:            true,
			//expectedEndBalance:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0))),
		},
		// TODO: This is unexpected behavior
		// We need validation to not charge fee if user doesn't have enough funds
		// {
		// 	name:                 "one user tries to create a gauge, has enough funds to pay for the create gauge fee but not enough to fill the gauge",
		// 	accountBalanceToFund: sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(60000000))),
		// 	gaugeAddition:        sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(30000000))),
		// 	expectedEndBalance:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(60000000))),
		// 	expectErr:            true,
		// },

	}

	for _, tc := range tests {
		suite.SetupTest()

		testAccountPubkey := secp256k1.GenPrivKeyFromSecret([]byte("acc")).PubKey()
		testAccountAddress := sdk.AccAddress(testAccountPubkey.Address())

		ctx := suite.Ctx
		incentivesKeepers := suite.App.IncentivesKeeper

		suite.FundAcc(testAccountAddress, tc.accountBalanceToFund)

		if tc.isModuleAccount {
			modAcc := authtypes.NewModuleAccount(authtypes.NewBaseAccount(testAccountAddress, testAccountPubkey, 1, 0),
				"module",
				"permission",
			)
			suite.App.AccountKeeper.SetModuleAccount(ctx, modAcc)
		}

		// suite.SetupManyLocks(1, defaultLiquidTokens, defaultLPTokens, defaultLockDuration)
		// distrTo := lockuptypes.QueryCondition{
		// 	LockQueryType: lockuptypes.ByDuration,
		// 	Denom:         defaultLPDenom,
		// 	Duration:      defaultLockDuration,
		// }

		// System under test.
		coins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(500000000)))
		gaugeID, _, _, _ := suite.SetupNewGauge(true, coins)
		if tc.nonexistentGauge {
			gaugeID = gaugeID + 50
		}

		//_, err := incentivesKeepers.CreateGauge(ctx, tc.isPerpetual, testAccountAddress, tc.gaugeAddition, distrTo, time.Time{}, 1)
		err := incentivesKeepers.AddToGaugeRewards(ctx, testAccountAddress, tc.gaugeAddition, gaugeID)

		if tc.expectErr {
			suite.Require().Error(err)
		} else {
			suite.Require().NoError(err)
		}

		bal := suite.App.BankKeeper.GetAllBalances(suite.Ctx, testAccountAddress)
		suite.Require().Equal(tc.expectedEndBalance.String(), bal.String(), "test: %v", tc.name)

		if tc.expectErr {
			suite.Require().Equal(tc.accountBalanceToFund.String(), bal.String(), "test: %v", tc.name)
		} else {
			finalAccountBalalance := tc.accountBalanceToFund.Sub(tc.gaugeAddition.Add(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(incentiveskeeper.AddToGaugeFee))))
			suite.Require().Equal(finalAccountBalalance.String(), bal.String(), "test: %v", tc.name)
		}

	}
}
