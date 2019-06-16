package hh

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "hh" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgTransferTokenToZone:
			return handleMsgTransferTokenToZone(ctx, keeper, msg)
		case MsgPutNFTokenOnTheMarket:
			return handleMsgPutNFTokenOnTheMarket(ctx, keeper, msg)
		case MsgBuyNFToken:
			return handleMsgBuyNFToken(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized hh Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgPutNFTokenOnTheMarket(ctx sdk.Context, keeper Keeper, msg MsgPutNFTokenOnTheMarket) sdk.Result {
	err := keeper.PutNFTokenOnTheMarket(ctx, msg.Token, msg.Sender)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}
	return sdk.Result{}
}

func handleMsgBuyNFToken(ctx sdk.Context, keeper Keeper, msg MsgBuyNFToken) sdk.Result {
	err := keeper.BuyNFToken(ctx, msg.NFTokenID, msg.Sender)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}
	return sdk.Result{}
}
