package authz

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kamleshesporg/mrmintchain/x/authz/keeper"
)

// BeginBlocker is called at the beginning of every block
func BeginBlocker(ctx sdk.Context, keeper keeper.Keeper) {
	// delete all the mature grants
	if err := keeper.DequeueAndDeleteExpiredGrants(ctx); err != nil {
		panic(err)
	}
}
