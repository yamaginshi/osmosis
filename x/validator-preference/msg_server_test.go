package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	valPref "github.com/osmosis-labs/osmosis/v12/x/validator-preference"
	"github.com/osmosis-labs/osmosis/v12/x/validator-preference/types"
)

func (suite *KeeperTestSuite) TestSetValidatorSetPreference() {
	suite.SetupTest()

	// setup 3 validators
	valAddrs := suite.SetupMultipleValidators(3)

	type param struct {
		delegator   sdk.AccAddress
		preferences []types.ValidatorPreference
	}

	tests := []struct {
		name       string
		param      param
		expectPass bool
	}{
		{
			name: "creation of new validator set",
			param: param{
				delegator: sdk.AccAddress([]byte("addr1---------------")),
				preferences: []types.ValidatorPreference{
					{
						ValOperAddress: valAddrs[0],
						Weight:         sdk.NewDecWithPrec(5, 1),
					},
					{
						ValOperAddress: valAddrs[1],
						Weight:         sdk.NewDecWithPrec(3, 1),
					},
					{
						ValOperAddress: valAddrs[2],
						Weight:         sdk.NewDecWithPrec(2, 1),
					},
				},
			},
			expectPass: true,
		},
		{
			name: "Update existing validator set",
			param: param{
				delegator: sdk.AccAddress([]byte("addr1---------------")),
				preferences: []types.ValidatorPreference{
					{
						ValOperAddress: valAddrs[0],
						Weight:         sdk.NewDecWithPrec(2, 1),
					},
					{
						ValOperAddress: valAddrs[1],
						Weight:         sdk.NewDecWithPrec(2, 1),
					},
					{
						ValOperAddress: valAddrs[2],
						Weight:         sdk.NewDecWithPrec(6, 1),
					},
				},
			},
			expectPass: true,
		},
		{
			name: "create validator set with unknown validator address",
			param: param{
				delegator: sdk.AccAddress([]byte("addr2---------------")),
				preferences: []types.ValidatorPreference{
					{
						ValOperAddress: "addr1---------------",
						Weight:         sdk.NewDec(1),
					},
					{
						ValOperAddress: valAddrs[1],
						Weight:         sdk.NewDecWithPrec(3, 1),
					},
				},
			},
			expectPass: false,
		},
		{
			name: "create validator set with weights != 1",
			param: param{
				delegator: sdk.AccAddress([]byte("addr3---------------")),
				preferences: []types.ValidatorPreference{
					{
						ValOperAddress: valAddrs[0],
						Weight:         sdk.NewDecWithPrec(5, 1),
					},
					{
						ValOperAddress: valAddrs[1],
						Weight:         sdk.NewDecWithPrec(3, 1),
					},
					{
						ValOperAddress: valAddrs[2],
						Weight:         sdk.NewDecWithPrec(3, 1),
					},
				},
			},
			expectPass: false,
		},
	}

	for _, test := range tests {
		suite.Run(test.name, func() {

			// setup message server
			msgServer := valPref.NewMsgServerImpl(suite.App.ValidatorPreferenceKeeper)
			c := sdk.WrapSDKContext(suite.Ctx)

			// call the create validator set preference
			_, err := msgServer.SetValidatorSetPreference(c, types.NewMsgSetValidatorSetPreference(test.param.delegator, test.param.preferences))
			if test.expectPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}

		})
	}
}
