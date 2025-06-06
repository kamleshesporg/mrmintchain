package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kamleshesporg/mrmintchain/x/mint/types"
	minttypes "github.com/kamleshesporg/mrmintchain/x/mint/types"
)

// InitGenesis new mint genesis
func (keeper Keeper) InitGenesis(ctx sdk.Context, ak types.AccountKeeper, data *types.GenesisState) {
	// keeper.SetMinter(ctx, data.Minter)
	keeper.SetMinter(ctx, minttypes.Minter{
		Inflation:        data.Minter.Inflation,
		AnnualProvisions: data.Minter.AnnualProvisions,
	})
	// keeper.SetParams(ctx, data.Params)
	keeper.SetParams(ctx, minttypes.Params{
		MintDenom:           data.Params.MintDenom,
		InflationRateChange: data.Params.InflationRateChange,
		InflationMax:        data.Params.InflationMax,
		InflationMin:        data.Params.InflationMin,
		GoalBonded:          data.Params.GoalBonded,
		BlocksPerYear:       data.Params.BlocksPerYear,
		HalvingInterval:     data.Params.HalvingInterval,
		MaxRewardTokens:     data.Params.MaxRewardTokens,
	}) // not cosmos-sdk Params

	// keeper.SetRewardPerBlock(ctx, data.Minter.RewardPerBlock)
	// keeper.SetHalvingInterval(ctx, data.Params.HalvingInterval)
	// keeper.SetMaxRewardSupply(ctx, data.Params.MaxRewardTokens)
	// keeper.SetTotalMinted(ctx, sdk.NewDecWithPrec(0, 0)) // usually starts at zero

	ak.GetModuleAccount(ctx, types.ModuleName)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (keeper Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	minter := keeper.GetMinter(ctx)
	params := keeper.GetParams(ctx)
	return types.NewGenesisState(minter, params)
}
