package hh

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName

// --------------------------------------------------------------------------
//
// MsgTransferNFTokenToZone
//
// --------------------------------------------------------------------------

// MsgTransferNFTokenToZone.
type MsgTransferNFTokenToZone struct {
	Sender    sdk.AccAddress
	Recipient sdk.AccAddress
	TokenID   string
	ZoneID    string
}

// NewMsgTransferNFTokenToZone is a constructor function for MsgTransferNFTokenToZone
func NewMsgTransferNFTokenToZone(tokenID, zoneID string, sender, recipient sdk.AccAddress) MsgTransferNFTokenToZone {
	return MsgTransferNFTokenToZone{
		TokenID:   tokenID,
		ZoneID:    zoneID,
		Sender:    sender,
		Recipient: recipient,
	}
}

// Route should return the name of the module
func (msg MsgTransferNFTokenToZone) Route() string { return RouterKey }

// Type should return the action
func (msg MsgTransferNFTokenToZone) Type() string { return "transfer_token_to_zone" }

// ValidateBasic runs stateless checks on the message
func (msg MsgTransferNFTokenToZone) ValidateBasic() sdk.Error {
	if len(msg.TokenID) == 0 || len(msg.ZoneID) == 0 {
		return sdk.ErrUnknownRequest("TokenID and/or ZoneID cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgTransferNFTokenToZone) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgTransferNFTokenToZone) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// --------------------------------------------------------------------------
//
// MsgPutNFTokenOnTheMarket
//
// --------------------------------------------------------------------------

// MsgPutNFTokenOnTheMarket.
type MsgPutNFTokenOnTheMarket struct {
	Sender sdk.AccAddress
	Token  NFT
}

// NewMsgPutNFTokenOnTheMarket is a constructor function for MsgPutNFTokenOnTheMarket
func NewMsgPutNFTokenOnTheMarket(token NFT, sender sdk.AccAddress) MsgPutNFTokenOnTheMarket {
	return MsgPutNFTokenOnTheMarket{
		Token:  token,
		Sender: sender,
	}
}

// Route should return the name of the module
func (msg MsgPutNFTokenOnTheMarket) Route() string { return RouterKey }

// Type should return the action
func (msg MsgPutNFTokenOnTheMarket) Type() string { return "put_token_on_the_market" }

// ValidateBasic runs stateless checks on the message
func (msg MsgPutNFTokenOnTheMarket) ValidateBasic() sdk.Error {
	if len(msg.Token.ID) == 0 {
		return sdk.ErrUnknownRequest("TokenID cannot be empty")
	}
	if !msg.Token.Price.IsAllPositive() {
		return sdk.ErrUnknownRequest("Token price should be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgPutNFTokenOnTheMarket) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgPutNFTokenOnTheMarket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// --------------------------------------------------------------------------
//
// MsgBuyNFToken
//
// --------------------------------------------------------------------------

// MsgBuyNFToken.
type MsgBuyNFToken struct {
	Sender    sdk.AccAddress
	NFTokenID string
	Price     sdk.Coins
}

// NewMsgBuyNFToken is a constructor function for MsgBuyNFToken
func NewMsgBuyNFToken(tokenID string, price sdk.Coins, sender sdk.AccAddress) MsgBuyNFToken {
	return MsgBuyNFToken{
		NFTokenID: tokenID,
		Price:     price,
		Sender:    sender,
	}
}

// Route should return the name of the module
func (msg MsgBuyNFToken) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBuyNFToken) Type() string { return "buy_token" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBuyNFToken) ValidateBasic() sdk.Error {
	if len(msg.NFTokenID) == 0 {
		return sdk.ErrUnknownRequest("TokenID cannot be empty")
	}
	if !msg.Price.IsAllPositive() {
		return sdk.ErrUnknownRequest("Token price should be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgBuyNFToken) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgBuyNFToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
