package hh

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgTransferNFTokenToZone{}, "hh/TransferNFTokenToZone", nil)
	cdc.RegisterConcrete(MsgPutNFTokenOnTheMarket{}, "hh/PutNFTokenOnTheMarket", nil)
	cdc.RegisterConcrete(MsgBuyNFToken{}, "hh/BuyNFToken", nil)
}

var ModuleCdc = codec.New()
