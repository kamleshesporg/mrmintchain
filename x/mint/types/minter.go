package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// type Minter struct {
// 	Inflation        sdk.Dec `json:"inflation" yaml:"inflation"`
// 	AnnualProvisions sdk.Dec `json:"annual_provisions" yaml:"annual_provisions"`
// 	RewardPerBlock   sdk.Int `json:"reward_per_block" yaml:"reward_per_block"`
// 	TotalMinted      sdk.Int `json:"total_minted" yaml:"total_minted"`
// }

// NewMinter returns a new Minter object with the given inflation, annual provisions,
// reward per block, and total minted tokens.
func NewMinter(inflation, annualProvisions sdk.Dec, rewardPerBlock sdk.Dec, totalMinted sdk.Dec) Minter {
	return Minter{
		Inflation:        inflation,
		AnnualProvisions: annualProvisions,
		RewardPerBlock:   rewardPerBlock,
		TotalMinted:      totalMinted,
	}
}

// InitialMinter returns an initial Minter object with a given inflation value.
func InitialMinter(inflation sdk.Dec) Minter {
	return NewMinter(
		inflation,
		sdk.NewDec(0),
		sdk.NewDec(0),
		sdk.NewDec(0),
	)
}

// DefaultInitialMinter returns a default initial Minter object with inflation 13%.
func DefaultInitialMinter() Minter {
	return InitialMinter(
		sdk.NewDecWithPrec(0, 2),
	)
}

// ValidateMinter validates the Minter fields.
func ValidateMinter(minter Minter) error {
	if minter.Inflation.IsNegative() {
		return fmt.Errorf("mint parameter Inflation should be positive, got %s", minter.Inflation.String())
	}
	if minter.RewardPerBlock.IsNegative() {
		return fmt.Errorf("reward_per_block cannot be negative, got %s", minter.RewardPerBlock.String())
	}
	if minter.TotalMinted.IsNegative() {
		return fmt.Errorf("total_minted cannot be negative, got %s", minter.TotalMinted.String())
	}
	return nil
}

// NextInflationRate returns the next inflation rate based on the bonded ratio.
func (m Minter) NextInflationRate(params Params, bondedRatio sdk.Dec) sdk.Dec {
	inflationRateChangePerYear := sdk.OneDec().
		Sub(bondedRatio.Quo(params.GoalBonded)).
		Mul(params.InflationRateChange)
	inflationRateChange := inflationRateChangePerYear.Quo(sdk.NewDec(int64(params.BlocksPerYear)))

	inflation := m.Inflation.Add(inflationRateChange)
	if inflation.GT(params.InflationMax) {
		inflation = params.InflationMax
	}
	if inflation.LT(params.InflationMin) {
		inflation = params.InflationMin
	}

	return inflation
}

// NextAnnualProvisions calculates annual provisions from inflation and total supply.
func (m Minter) NextAnnualProvisions(_ Params, totalSupply math.Int) sdk.Dec {
	return m.Inflation.MulInt(totalSupply)
}

// BlockProvision returns block rewards and respects max reward cap.
func (m Minter) BlockProvision(params Params) sdk.Coin {
	if params.MaxRewardTokens.IsPositive() && m.TotalMinted.GTE(params.MaxRewardTokens) {
		// Reward limit reached
		return sdk.NewCoin(params.MintDenom, sdk.NewInt(0))
	}

	// Return the fixed reward per block (halving logic can be added here later if needed)
	return sdk.NewCoin(params.MintDenom, sdk.NewInt(0))
}
