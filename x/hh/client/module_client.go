package client

import (
	hhcmd "dgamingfoundation/hackathon-hub/x/hh/client/cli"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"
)

// ModuleClient exports all client functionality from this module
type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{storeKey, cdc}
}

// GetQueryCmd returns the cli query commands for this module
func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	// Group hh queries under a subcommand
	hhQueryCmd := &cobra.Command{
		Use:   "hh",
		Short: "Querying commands for the hh module",
	}

	hhQueryCmd.AddCommand(client.GetCommands(
		hhcmd.GetCmdTokenInfo(mc.storeKey, mc.cdc),
		hhcmd.GetCmdListTokens(mc.storeKey, mc.cdc),
		hhcmd.GetCmdTransferInfo(mc.storeKey, mc.cdc),
	)...)

	return hhQueryCmd
}

// GetTxCmd returns the transaction commands for this module
func (mc ModuleClient) GetTxCmd() *cobra.Command {
	hhTxCmd := &cobra.Command{
		Use:   "hh",
		Short: "hh transactions subcommands",
	}

	hhTxCmd.AddCommand(client.PostCommands(
		hhcmd.GetCmdTransferToken(mc.cdc),
	)...)

	return hhTxCmd
}
