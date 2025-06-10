package client

import (
	govclient "github.com/kamleshesporg/mrmintchain/x/gov/client"
	"github.com/kamleshesporg/mrmintchain/x/upgrade/client/cli"
)

var (
	LegacyProposalHandler       = govclient.NewProposalHandler(cli.NewCmdSubmitLegacyUpgradeProposal)
	LegacyCancelProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitLegacyCancelUpgradeProposal)
)
