package keeper

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// "github.com/kamleshesporg/mrmintchain/x/mint/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/kamleshesporg/mrmintchain/x/mint/types"
)

// Keeper of the mint store
type Keeper struct {
	cdc              codec.BinaryCodec
	storeKey         storetypes.StoreKey
	paramSpace       paramtypes.Subspace
	stakingKeeper    types.StakingKeeper
	bankKeeper       types.BankKeeper
	feeCollectorName string
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec, key storetypes.StoreKey, paramSpace paramtypes.Subspace,
	sk types.StakingKeeper, ak types.AccountKeeper, bk types.BankKeeper,
	feeCollectorName string,
) Keeper {
	// ensure mint module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the mint module account has not been set")
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:              cdc,
		storeKey:         key,
		paramSpace:       paramSpace,
		stakingKeeper:    sk,
		bankKeeper:       bk,
		feeCollectorName: feeCollectorName,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// get the minter
func (k Keeper) GetMinter(ctx sdk.Context) (minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.MinterKey)
	if b == nil {
		panic("stored minter should not have been nil")
	}

	k.cdc.MustUnmarshal(b, &minter)
	return
}

// set the minter
func (k Keeper) SetMinter(ctx sdk.Context, minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&minter)
	store.Set(types.MinterKey, b)
}

// GetParams returns the total set of minting parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of minting parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// StakingTokenSupply implements an alias call to the underlying staking keeper's
// StakingTokenSupply to be used in BeginBlocker.
func (k Keeper) StakingTokenSupply(ctx sdk.Context) math.Int {
	return k.stakingKeeper.StakingTokenSupply(ctx)
}

// BondedRatio implements an alias call to the underlying staking keeper's
// BondedRatio to be used in BeginBlocker.
func (k Keeper) BondedRatio(ctx sdk.Context) sdk.Dec {
	return k.stakingKeeper.BondedRatio(ctx)
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx sdk.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}
	return k.bankKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

// AddCollectedFees implements an alias call to the underlying supply keeper's
// AddCollectedFees to be used in BeginBlocker.
func (k Keeper) AddCollectedFees(ctx sdk.Context, fees sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, fees)
}

// func (k Keeper) GetRewardPerBlock(ctx sdk.Context) sdk.Int {
// 	bz := k.store(ctx).Get(types.KeyRewardPerBlock)
// 	if bz == nil {
// 		return sdk.NewInt(5) // default value
// 	}
// 	dkmath.NewIntFromBigEndian(bz)
// }

// func (k Keeper) SetRewardPerBlock(ctx sdk.Context, reward sdk.Int) {
// 	bz := reward.BigInt().Bytes()
// 	k.store(ctx).Set(types.KeyRewardPerBlock, bz)
// }

/** new_chanegs **/
func (k Keeper) GetRewardPerBlock(ctx sdk.Context) sdkmath.Int {
	bz := k.store(ctx).Get(types.KeyRewardPerBlock)
	if bz == nil {
		return sdk.NewInt(5) // default value
	}

	return sdkmath.NewIntFromBigInt(new(big.Int).SetBytes(bz))
}

func (k Keeper) SetRewardPerBlock(ctx sdk.Context, reward sdk.Int) {
	bz := reward.BigInt().Bytes()
	k.store(ctx).Set(types.KeyRewardPerBlock, bz)
}
func (k Keeper) GetTotalMinted(ctx sdk.Context) sdk.Int {
	bz := k.store(ctx).Get(types.KeyTotalMinted)
	if bz == nil {
		return sdk.ZeroInt()
	}
	return sdkmath.NewIntFromBigInt(new(big.Int).SetBytes(bz))
}

func (k Keeper) SetTotalMinted(ctx sdk.Context, total sdk.Int) {
	bz := total.BigInt().Bytes()
	k.store(ctx).Set(types.KeyTotalMinted, bz)
}
func (k Keeper) GetHalvingInterval(ctx sdk.Context) int64 {
	bz := k.store(ctx).Get(types.KeyHalvingInterval)
	fmt.Printf("1111 GetHalvingInterval =>%d", bz)
	if bz == nil {
		return 20 // default interval in blocks
	}
	fmt.Printf("2222 GetHalvingInterval =>%d", bz)
	return int64(binary.BigEndian.Uint64(bz))
}

func (k Keeper) SetHalvingInterval(ctx sdk.Context, interval int64) {
	bz := sdk.Uint64ToBigEndian(uint64(interval))
	k.store(ctx).Set(types.KeyHalvingInterval, bz)
}
func (k Keeper) GetMaxRewardSupply(ctx sdk.Context) sdk.Int {
	bz := k.store(ctx).Get(types.KeyMaxRewardSupply)
	if bz == nil {
		return sdk.NewInt(10000)
	}
	return sdkmath.NewIntFromBigInt(new(big.Int).SetBytes(bz))
}

func (k Keeper) SetMaxRewardSupply(ctx sdk.Context, max sdk.Int) {
	bz := max.BigInt().Bytes()
	k.store(ctx).Set(types.KeyMaxRewardSupply, bz)
}
func (k Keeper) store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

/** new_changes_end */
