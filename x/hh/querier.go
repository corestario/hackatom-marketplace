package hh

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Query endpoints supported by the hh Querier.
const (
	QueryNFToken  = "NFToken"
	QueryNFTokens = "NFTokens"
	QueryTransfer = "Transfer"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryNFToken:
			return queryNFToken(ctx, path[1:], req, keeper)
		case QueryNFTokens:
			return queryNFTokens(ctx, req, keeper)
		case QueryTransfer:
			return queryTransfer(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown hh query endpoint")
		}
	}
}

// nolint: unparam
func queryNFToken(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	tokenID := path[0]

	token := keeper.GetNFToken(ctx, tokenID)
	if token == nil {
		return nil, sdk.ErrUnknownRequest("token does not exist")
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, token)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

func queryNFTokens(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	tokens := QueryResNFTokens(keeper.GetNFTokens(ctx))

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, tokens)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}

type QueryResNFTokens []NFT

func (m QueryResNFTokens) String() string {
	var out []string
	for _, token := range m {
		out = append(out, token.String())
	}
	return strings.Join(out, "\n")
}

func queryTransfer(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	transferID := path[0]
	transfer, err1 := keeper.GetTransfer(ctx, transferID)
	if err1 != nil {
		return nil, sdk.ErrUnknownRequest("token does not exist")
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, transfer)
	if err2 != nil {
		panic("could not marshal result to JSON")
	}

	return bz, nil
}
