package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState defines the mint module's genesis state.
// type GenesisState struct {
// 	Minter Minter `json:"minter" yaml:"minter"`
// 	Params Params `json:"params" yaml:"params"`
// }

// InflationCalculationFn defines the function required to calculate inflation rate during
// BeginBlock. It receives the minter and params stored in the keeper, along with the current
// bondedRatio and returns the newly calculated inflation rate.
// It can be used to specify a custom inflation calculation logic, instead of relying on the
// default logic provided by the sdk.
type InflationCalculationFn func(ctx sdk.Context, minter Minter, params Params, bondedRatio sdk.Dec) sdk.Dec

// DefaultInflationCalculationFn is the default function used to calculate inflation.
func DefaultInflationCalculationFn(_ sdk.Context, minter Minter, params Params, bondedRatio sdk.Dec) sdk.Dec {
	return minter.NextInflationRate(params, bondedRatio)
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(minter Minter, params Params) *GenesisState {
	return &GenesisState{
		Minter: &minter,
		Params: &params,
	}
}

// // DefaultGenesisState creates a default GenesisState object
// func DefaultGenesisState() *GenesisState {
// 	return &GenesisState{
// 		Minter: DefaultInitialMinter(),
// 		Params: DefaultParams(),
// 	}
// }

func DefaultGenesisState() *GenesisState {
	m := DefaultInitialMinter()
	p := DefaultParams()
	return &GenesisState{
		Minter: &m,
		Params: &p,
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants hold.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	// This validates Inflation, RewardPerBlock, TotalMinted, etc.
	return ValidateMinter(*data.Minter)
}
