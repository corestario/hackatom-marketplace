package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"dgamingfoundation/hackathon-hub/x/hh"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func GetCmdTransferToken(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "transfer-token [tokenID] [zoneID] [recipient]",
		Short: "bid for existing name or claim new name",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			recipient, err := sdk.AccAddressFromHex(args[2])
			if err != nil {
				return fmt.Errorf("failed to parse recipient address: %v", err)
			}

			msg := hh.NewMsgTransferNFTokenToZone(args[0], args[1], cliCtx.GetFromAddress(), recipient)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
