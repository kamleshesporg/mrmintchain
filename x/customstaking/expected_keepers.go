package customstaking

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StakingKeeper defines the expected staking keeper
type StakingKeeper interface {
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (stakingtypes.Validator, bool)
	BeginUnbonding(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec) (time.Time, error)
	Delegate(
		ctx sdk.Context,
		delAddr sdk.AccAddress,
		bondAmt sdk.Int,
		bondStatus stakingtypes.BondStatus,
		validator stakingtypes.Validator,
		subtractAccount bool,
	) (newShares sdk.Dec, err error)
	// Add other required methods as needed
}
