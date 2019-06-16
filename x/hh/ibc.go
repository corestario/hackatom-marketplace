package hh

import (
	"encoding/json"
	"fmt"
	commitment "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ConnectionID         = "market_connection"
	CounterpartyID       = "hub"
	CounterpartyClientID = "me"
)

// Handle a message to create nft
func handleMsgTransferTokenToZone(
	ctx sdk.Context,
	keeper Keeper,
	msg MsgTransferTokenToZone,
) sdk.Result {

	token := keeper.GetNFToken(ctx, msg.TokenID)

	packet := NewSendTokenPacket(&token.BaseNFT)
	if err := keeper.ibcKeeper.Send(ctx, msg.ConnectionID, msg.CounterpartyClientID, packet); err != nil {
		return sdk.Result{Code: sdk.CodeUnknownRequest, Log: err.Error()}
	}

	keeper.DeleteNFT(ctx, msg.TokenID)

	return sdk.Result{}
}

// --------------------------------------------------------------------------
//
// TransferTokenToHub
//
// --------------------------------------------------------------------------

type MsgTransferTokenToZone struct {
	Owner   sdk.AccAddress
	TokenID string

	//fixme fill it
	ConnectionID         string
	CounterpartyID       string
	CounterpartyClientID string
}

// NewMsgCreateNFT is a constructor function for MsgCreateNFT
func NewMsgTransferTokenToZone(owner sdk.AccAddress, tokenID string, zoneID string) MsgTransferTokenToZone {
	return MsgTransferTokenToZone{
		Owner:                owner,
		TokenID:              tokenID,
		CounterpartyClientID: zoneID,
	}
}

// Route should return the name of the module
func (msg MsgTransferTokenToZone) Route() string { return RouterKey }

// Type should return the action
func (msg MsgTransferTokenToZone) Type() string { return "transfer_token_to_hub" }

// ValidateBasic runs stateless checks on the message
func (msg MsgTransferTokenToZone) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgTransferTokenToZone) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners defines whose signature is required
func (msg MsgTransferTokenToZone) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

type SendTokenPacket struct {
	Token *BaseNFT `json:"token"`
}

func NewSendTokenPacket(token *BaseNFT) *SendTokenPacket {
	return &SendTokenPacket{Token: token}
}

func (m *SendTokenPacket) Timeout() uint64 {
	return math.MaxUint64
}

func (m *SendTokenPacket) Commit() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("failed to marshal SellTokenPacket packet: %v", err))
	}

	return data
}

type ProofPacket struct {
}

func (p *ProofPacket) GetKey() []byte {
	return []byte{}
}

func (p *ProofPacket) Verify(commitment.Root, []byte) error {
	return nil
}
