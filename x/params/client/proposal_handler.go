package client

import (
	govclient "github.com/kamleshesporg/mrmintchain/x/gov/client"
	"github.com/kamleshesporg/mrmintchain/x/params/client/cli"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewSubmitParamChangeProposalTxCmd)
