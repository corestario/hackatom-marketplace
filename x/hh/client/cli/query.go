package cli

import (
	"fmt"

	"dgamingfoundation/hackathon-hub/x/hh"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetCmdTokenInfo(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "token-info [id]",
		Short: "See token data",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			tokenID := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/NFToken/%s", queryRoute, tokenID), nil)
			if err != nil {
				fmt.Printf("could not find tokenID - %s: %v\n", tokenID, err)
				return nil
			}

			var out hh.NFT
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdListTokens(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "list-tokens",
		Short: "See token data",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/NFTokens", queryRoute), nil)
			if err != nil {
				fmt.Printf("could not get token list: %v", err)
				return nil
			}

			var out hh.QueryResNFTokens
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdTransferInfo(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "token-info [id]",
		Short: "See token data",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			transferID := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/Transfer/%s", queryRoute, transferID), nil)
			if err != nil {
				fmt.Printf("could not find tokenID - %s: %v\n", transferID, err)
				return nil
			}

			var out hh.Transfer
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}
