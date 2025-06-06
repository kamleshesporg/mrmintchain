package main_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client/flags"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	"github.com/kamleshesporg/mrmintchain/app"
	mrmintchaind "github.com/kamleshesporg/mrmintchain/cmd/mrmintchaind"
)

func TestInitCmd(t *testing.T) {
	rootCmd, _ := mrmintchaind.NewRootCmd()
	rootCmd.SetArgs([]string{
		"init",            // Test the init cmd
		"mrmintchaintest", // Moniker
		fmt.Sprintf("--%s=%s", cli.FlagOverwrite, "true"), // Overwrite genesis.json, in case it already exists
		fmt.Sprintf("--%s=%s", flags.FlagChainID, "mrmintchain_9000-1"),
	})

	err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome)
	require.NoError(t, err)
}
