// Copyright 2021 Evmos Foundation
// This file is part of Evmos' mrmintchain library.
//
// The mrmintchain library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The mrmintchain library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the mrmintchain library. If not, see https://github.com/kamleshesporg/mrmintchain/blob/main/LICENSE
package feemarket

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/kamleshesporg/mrmintchain/x/feemarket/keeper"
	"github.com/kamleshesporg/mrmintchain/x/feemarket/types"
)

// InitGenesis initializes genesis state based on exported genesis
func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	data types.GenesisState,
) []abci.ValidatorUpdate {
	err := k.SetParams(ctx, data.Params)
	if err != nil {
		panic(errorsmod.Wrap(err, "could not set parameters at genesis"))
	}

	k.SetBlockGasWanted(ctx, data.BlockGas)

	return []abci.ValidatorUpdate{}
}

// ExportGenesis exports genesis state of the fee market module
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:   k.GetParams(ctx),
		BlockGas: k.GetBlockGasWanted(ctx),
	}
}
