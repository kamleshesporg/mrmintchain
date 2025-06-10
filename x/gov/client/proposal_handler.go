// package client

// import (
// 	"github.com/spf13/cobra"
// )

// // function to create the cli handler
// type CLIHandlerFn func() *cobra.Command

// // ProposalHandler wraps CLIHandlerFn
// type ProposalHandler struct {
// 	CLIHandler CLIHandlerFn
// }

// // NewProposalHandler creates a new ProposalHandler object
// func NewProposalHandler(cliHandler CLIHandlerFn) ProposalHandler {
// 	return ProposalHandler{
// 		CLIHandler: cliHandler,
// 	}
// }

package client

import (
	"github.com/spf13/cobra"
)

// ProposalHandler wraps a CLI handler function
type ProposalHandler struct {
	CLIHandler func() *cobra.Command
}

// NewProposalHandler creates a new ProposalHandler
func NewProposalHandler(cliHandler func() *cobra.Command) ProposalHandler {
	return ProposalHandler{
		CLIHandler: cliHandler,
	}
}
