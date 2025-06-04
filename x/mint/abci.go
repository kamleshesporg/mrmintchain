// package mint

// import (
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/kamleshesporg/mrmintchain/x/mint/keeper"
// 	abci "github.com/tendermint/tendermint/abci/types"
// )

// // // BeginBlocker mints new tokens for the previous block.
// // func BeginBlocker(ctx sdk.Context, k keeper.Keeper, ic types.InflationCalculationFn) {
// // 	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

// // 	// fetch stored minter & params
// // 	minter := k.GetMinter(ctx)
// // 	params := k.GetParams(ctx)

// // 	// recalculate inflation rate
// // 	totalStakingSupply := k.StakingTokenSupply(ctx)
// // 	bondedRatio := k.BondedRatio(ctx)
// // 	minter.Inflation = ic(ctx, minter, params, bondedRatio)
// // 	minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalStakingSupply)
// // 	k.SetMinter(ctx, minter)

// // 	// mint coins, update supply
// // 	mintedCoin := minter.BlockProvision(params)
// // 	mintedCoins := sdk.NewCoins(mintedCoin)

// // 	err := k.MintCoins(ctx, mintedCoins)
// // 	if err != nil {
// // 		panic(err)
// // 	}

// // 	// send the minted coins to the fee collector account
// // 	err = k.AddCollectedFees(ctx, mintedCoins)
// // 	if err != nil {
// // 		panic(err)
// // 	}

// // 	if mintedCoin.Amount.IsInt64() {
// // 		defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
// // 	}

// // 	ctx.EventManager().EmitEvent(
// // 		sdk.NewEvent(
// // 			types.EventTypeMint,
// // 			sdk.NewAttribute(types.AttributeKeyBondedRatio, bondedRatio.String()),
// // 			sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
// // 			sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
// // 			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
// // 		),
// // 	)
// // }

// func BeginBlocker(ctx sdk.Context, k keeper.Keeper, _ abci.RequestBeginBlock) {
// 	height := ctx.BlockHeight()

// 	rewardPerBlock := k.GetRewardPerBlock(ctx)
// 	totalMinted := k.GetTotalMintedTokens(ctx)
// 	maxRewardTokens := k.GetMaxRewardTokens(ctx)
// 	halvingInterval := k.GetHalvingInterval(ctx)

// 	if totalMinted >= maxRewardTokens {
// 		return // stop minting rewards
// 	}

// 	// Halving every halvingInterval blocks
// 	if height%halvingInterval == 0 && height != 0 {
// 		rewardPerBlock = rewardPerBlock / 2
// 		k.SetRewardPerBlock(ctx, rewardPerBlock)
// 	}

// 	// Mint and distribute rewards without exceeding max
// 	mintAmount := rewardPerBlock
// 	if totalMinted+mintAmount > maxRewardTokens {
// 		mintAmount = maxRewardTokens - totalMinted
// 	}

// 	coins := sdk.NewCoins(sdk.NewCoin("aphoton", sdk.NewInt(mintAmount)))
// 	err := k.MintCoins(ctx, coins)
// 	if err != nil {
// 		// handle error
// 		return
// 	}

// 	k.DistributeRewards(ctx, coins)
// 	k.SetTotalMintedTokens(ctx, totalMinted+mintAmount)
// }

package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	"github.com/kamleshesporg/mrmintchain/x/mint/keeper"
)

func BeginBlocker(ctx types.Context, k keeper.Keeper) {
	blockHeight := ctx.BlockHeight()

	maxRewardSupply := k.GetMaxRewardSupply(ctx)
	totalMinted := k.GetTotalMinted(ctx)
	rewardPerBlock := k.GetRewardPerBlock(ctx)
	halvingInterval := k.GetHalvingInterval(ctx)

	// Skip if already maxed out
	if totalMinted.GTE(maxRewardSupply) || rewardPerBlock.IsZero() {
		return
	}

	// Halving logic
	if blockHeight > 0 && blockHeight%halvingInterval == 0 {
		newReward := rewardPerBlock.QuoRaw(2)
		if newReward.IsZero() {
			newReward = sdk.NewInt(1) // minimum reward to avoid zero
		}
		k.SetRewardPerBlock(ctx, newReward)
	}

	// Final minting amount (adjust for remaining supply cap)
	remaining := maxRewardSupply.Sub(totalMinted)
	mintAmount := rewardPerBlock
	if rewardPerBlock.GT(remaining) {
		mintAmount = remaining
	}

	// Mint new coins
	mintCoin := sdk.NewCoin(k.GetParams(ctx).MintDenom, mintAmount)
	err := k.MintCoins(ctx, sdk.NewCoins(mintCoin))
	if err != nil {
		panic(err)
	}

	// Send minted coins to fee collector or specific address
	err = k.AddCollectedFees(ctx, sdk.NewCoins(mintCoin))
	if err != nil {
		panic(err)
	}

	// Update total minted
	k.SetTotalMinted(ctx, totalMinted.Add(mintAmount))
}
