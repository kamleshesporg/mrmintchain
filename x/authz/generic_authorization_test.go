package authz_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/kamleshesporg/mrmintchain/x/authz"
)

func TestGenericAuthorization(t *testing.T) {
	t.Log("verify ValidateBasic returns nil for service msg")
	a := authz.NewGenericAuthorization(banktypes.SendAuthorization{}.MsgTypeURL())
	require.NoError(t, a.ValidateBasic())
	require.Equal(t, banktypes.SendAuthorization{}.MsgTypeURL(), a.Msg)
}
