package hh

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	NFTRecords []NFT `json:"nft_records"`
	AuthData auth.GenesisState   `json:"auth"`
	BankData bank.GenesisState   `json:"bank"`
	Accounts []*auth.BaseAccount `json:"accounts"`
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

func InitGenesis(ctx sdk.Context, keeper Keeper, state GenesisState) []abci.ValidatorUpdate {
	for _, record := range state.NFTRecords {
		keeper.StoreNFT(ctx, record, sdk.AccAddress{})
	}

	for _, acc := range state.Accounts {
		acc.AccountNumber = keeper.accountKeeper.GetNextAccountNumber(ctx)
		keeper.accountKeeper.SetAccount(ctx, acc)
	}

	keeper.coinKeeper.SetSendEnabled(ctx, true)
	keeper.accountKeeper.SetParams(ctx, auth.DefaultParams())

	auth.InitGenesis(ctx, keeper.accountKeeper, keeper.feeCollectionKeeper, state.AuthData)
	bank.InitGenesis(ctx, keeper.coinKeeper, state.BankData)

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
