// package customstaking

// import (
// 	"bytes"
// 	"time"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
// 	"github.com/cosmos/cosmos-sdk/x/staking/types"
// )

// // CustomKeeper मूल Keeper को एम्बेड करता है
// type CustomKeeper struct {
// 	*keeper.Keeper // पॉइंटर के रूप में एम्बेड करें (यह महत्वपूर्ण है)

// 	validatorUnbondingTime time.Duration
// 	delegatorUnbondingTime time.Duration
// }

// // NewCustomKeeper Keeper का नया इंस्टेंस बनाता है
// func NewCustomKeeper(
// 	k *keeper.Keeper, // पॉइंटर के रूप में
// 	validatorTime time.Duration,
// 	delegatorTime time.Duration,
// ) *CustomKeeper {
// 	return &CustomKeeper{
// 		Keeper:                 k, // मूल Keeper असाइन करें
// 		validatorUnbondingTime: validatorTime,
// 		delegatorUnbondingTime: delegatorTime,
// 	}
// }

// func (k *CustomKeeper) BeginUnbonding(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec) (time.Time, error) {
// 	// GetValidator सीधे उपलब्ध होगा
// 	validator, found := k.Keeper.GetValidator(ctx, valAddr)
// 	if !found {
// 		return time.Time{}, types.ErrNoValidatorFound
// 	}

// 	// अनबॉन्डिंग टाइम निर्धारित करें
// 	unbondingTime := k.delegatorUnbondingTime
// 	if bytes.Equal(delAddr.Bytes(), validator.GetOperator().Bytes()) {
// 		unbondingTime = k.validatorUnbondingTime
// 	}

// 	return ctx.BlockHeader().Time.Add(unbondingTime), nil
// }

// // Getter methods
// func (k CustomKeeper) ValidatorUnbondingTime(ctx sdk.Context) time.Duration {
// 	return k.validatorUnbondingTime
// }

// func (k CustomKeeper) DelegatorUnbondingTime(ctx sdk.Context) time.Duration {
// 	return k.delegatorUnbondingTime
// }

// v2=============================

// package customstaking

// import (
// 	"fmt"
// 	"time"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
// 	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
// )

// // CustomKeeper wraps the original keeper
// type CustomKeeper struct {
// 	keeper.Keeper // Embed the concrete keeper type

// 	validatorUnbondingTime time.Duration
// 	delegatorUnbondingTime time.Duration
// }

// // NewCustomKeeper creates a new CustomKeeper
// func NewCustomKeeper(
// 	k keeper.Keeper, // Use concrete keeper type
// 	validatorTime time.Duration,
// 	delegatorTime time.Duration,
// ) *CustomKeeper {
// 	return &CustomKeeper{
// 		Keeper:                 k,
// 		validatorUnbondingTime: validatorTime,
// 		delegatorUnbondingTime: delegatorTime,
// 	}
// }

// // BeginUnbonding override
// func (k *CustomKeeper) BeginUnbonding(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec) (time.Time, error) {
// 	fmt.Println("===========================================कस्टम BeginUnbonding कॉल हुआ")
// 	panic("कस्टम BeginUnbonding कॉल हुआ")

// 	validator, found := k.Keeper.GetValidator(ctx, valAddr)
// 	if !found {
// 		return time.Time{}, stakingtypes.ErrNoValidatorFound
// 	}
// 	fmt.Printf("===========================================वैलिडेटर चेक: %v\n", "isValidator")

// 	unbondingTime := k.delegatorUnbondingTime
// 	if delAddr.Equals(validator.GetOperator()) {
// 		unbondingTime = k.validatorUnbondingTime
// 	}
// 	fmt.Printf("===========================================अनबॉन्डिंग टाइम सेट: %v\n", unbondingTime)

// 	return ctx.BlockHeader().Time.Add(unbondingTime), nil
// }

// // Getter methods
// func (k CustomKeeper) ValidatorUnbondingTime(ctx sdk.Context) time.Duration {
// 	return k.validatorUnbondingTime
// }

// func (k CustomKeeper) DelegatorUnbondingTime(ctx sdk.Context) time.Duration {
// 	return k.delegatorUnbondingTime
// }

// v3========================================
package customstaking

import (
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// CustomKeeper wraps the original keeper and implements all methods
type CustomKeeper struct {
	keeper.Keeper // Embed the concrete keeper

	validatorUnbondingTime time.Duration
	delegatorUnbondingTime time.Duration
}

// Interface assertion
var _ StakingKeeper = (*CustomKeeper)(nil)

func NewCustomKeeper(
	k keeper.Keeper,
	validatorTime time.Duration,
	delegatorTime time.Duration,
) *CustomKeeper {
	return &CustomKeeper{
		Keeper:                 k,
		validatorUnbondingTime: validatorTime,
		delegatorUnbondingTime: delegatorTime,
	}
}

// BeginUnbonding override
func (k *CustomKeeper) BeginUnbonding(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec) (time.Time, error) {
	fmt.Println("===========================================कस्टम BeginUnbonding कॉल हुआ")
	panic("कस्टम BeginUnbonding कॉल हुआ")

	validator, found := k.Keeper.GetValidator(ctx, valAddr)
	if !found {
		return time.Time{}, types.ErrNoValidatorFound
	}

	unbondingTime := k.delegatorUnbondingTime
	if delAddr.Equals(validator.GetOperator()) {
		unbondingTime = k.validatorUnbondingTime
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent("custom_unbonding",
			sdk.NewAttribute("time", unbondingTime.String()),
			sdk.NewAttribute("is_validator", strconv.FormatBool(delAddr.Equals(validator.GetOperator()))),
		),
	)

	return ctx.BlockHeader().Time.Add(unbondingTime), nil
}
