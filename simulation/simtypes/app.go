package simtypes

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	gammtypes "github.com/osmosis-labs/osmosis/v10/x/gamm/types"

	lockuptypes "github.com/osmosis-labs/osmosis/v10/x/lockup/types"
)

type App interface {
	GetBaseApp() *baseapp.BaseApp
	AppCodec() codec.Codec
	GetAccountKeeper() AccountKeeper
	GetBankKeeper() BankKeeper
	GetStakingKeeper() StakingKeeper
	GetLockupKeeper() LockupKeeper
	GetGAMMKeeper() GAMMKeeper
	GetGovKeeper() GovKeeper
}

type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	// GetAllAccounts(ctx sdk.Context) []authtypes.AccountI
}

type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// TODO: Revisit
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

type StakingKeeper interface {
	GetAllValidators(ctx sdk.Context) (validators []stakingtypes.Validator)
	GetValidatorDelegations(ctx sdk.Context, valAddr sdk.ValAddress) (delegations []stakingtypes.Delegation)
	GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator
}

type LockupKeeper interface {
	GetPeriodLocks(ctx sdk.Context) ([]lockuptypes.PeriodLock, error)
	GetAccountPeriodLocks(ctx sdk.Context, addr sdk.AccAddress) []lockuptypes.PeriodLock
}

type GAMMKeeper interface {
	GetPoolsAndPoke(ctx sdk.Context) ([]gammtypes.PoolI, error)
}

type GovKeeper interface {
	GetProposal(ctx sdk.Context, proposalID uint64) (govtypes.Proposal, bool)
	GetProposals(ctx sdk.Context) (proposals govtypes.Proposals)
	GetProposalID(ctx sdk.Context) (uint64, error)
	GetDepositParams(ctx sdk.Context) govtypes.DepositParams
	ActivateVotingPeriod(ctx sdk.Context, proposal govtypes.Proposal)
}
