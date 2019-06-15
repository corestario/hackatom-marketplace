package hh

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	NFTRecords []NFT `json:"nft_records"`
}

func NewGenesisState(nftRecords []NFT) GenesisState {
	return GenesisState{NFTRecords: nftRecords}
}

func ValidateGenesis(data GenesisState) error {
	for range data.NFTRecords {
		// everything is fine
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		NFTRecords: []NFT{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.NFTRecords {
		keeper.StoreNFT(ctx, record, sdk.AccAddress{})
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []NFT
	iterator := k.GetNFTIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		id := string(iterator.Key())
		var nft *NFT
		nft = k.GetNFToken(ctx, id)

		if nft != nil {
			records = append(records, *nft)
		}
	}
	return GenesisState{NFTRecords: records}
}
