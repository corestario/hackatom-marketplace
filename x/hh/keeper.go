package hh

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"strings"
)

// StoreKey to be used when creating the KVStore
const StoreKey = "hh"

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the hh Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

func (k Keeper) TransferNFTokenToZone(ctx sdk.Context, nfToken NFT, zoneID string, sender, recipient sdk.AccAddress) error {
	if !k.getNFTOwner(ctx, nfToken.ID).Empty() {
		return errors.New("token has already existed")
	}
	k.StoreNFT(ctx, nfToken, recipient)
	return k.PutNFTokenOnTheMarket(ctx, nfToken.BaseNFT, nfToken.Price, recipient)
}

func (k Keeper) GetTransfer(ctx sdk.Context, transferID string) (*Transfer, error) {
	// TODO: implement.
	return nil, nil
}

func (k Keeper) PutNFTokenOnTheMarket(ctx sdk.Context, token BaseNFT, price sdk.Coins, sender sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)
	if sender.Empty() {
		return errors.New("empty sender")
	}

	if store.Has(composePutNFTToMarketKey(token.ID)) {
		return errors.New("nft has already existed on market")
	}

	if k.getNFTOwner(ctx, token.ID).Equals(sender) == false {
		return errors.New("you are not owner of the nft")
	}

	nftOnSale := NFT{token, true, price}
	nftOnSaleBin := k.cdc.MustMarshalBinaryBare(nftOnSale)

	store.Set(composePutNFTToMarketKey(token.ID), nftOnSaleBin)
	return nil
}

func (k Keeper) setNFTOwner(ctx sdk.Context, NFTTokenID string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(composePutNFTOwnerKey(NFTTokenID), owner.Bytes())
}

func (k Keeper) getNFTOwner(ctx sdk.Context, NFTTokenID string) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	return sdk.AccAddress(store.Get(composePutNFTOwnerKey(NFTTokenID)))
}

func (k Keeper) BuyNFToken(ctx sdk.Context, nfTokenID string, buyer sdk.AccAddress) error {
	// TODO: implement.
	return nil
}

func (k Keeper) StoreNFT(ctx sdk.Context, nft NFT, owner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	if store.Has([]byte(nft.ID)) {
		return
	}

	store.Set([]byte(nft.ID), k.cdc.MustMarshalBinaryBare(nft))
	store.Set(composePutNFTOwnerKey(nft.ID), owner.Bytes())
}

func (k Keeper) GetNFToken(ctx sdk.Context, tokenID string) *NFT {
	store := ctx.KVStore(k.storeKey)
	nftBin := store.Get(composePutNFTToMarketKey(tokenID))

	nftOnSale := new(NFT)
	err := k.cdc.UnmarshalBinaryBare(nftBin, nftOnSale)
	if err != nil {
		ctx.Logger().Error("error while GetNFToken ", "tokenID", tokenID)
		return nil
	}

	return nftOnSale
}

func (k Keeper) GetNFTokens(ctx sdk.Context) sdk.Iterator {
	// TODO: implement.
	return nil
}

func (k Keeper) GetNFTokensOnSaleList(ctx sdk.Context) []NFT {
	it := k.GetNFTIterator(ctx)
	var nftList []NFT
	for {
		if it.Valid() == false {
			break
		}

		var price sdk.Coins
		err := json.Unmarshal(it.Value(), &price)
		if err != nil {
			fmt.Println("json.Unmarshal err", err)
		}
		nftList = append(nftList, NFT{
			BaseNFT: BaseNFT{
				ID: strings.TrimPrefix(string(it.Key()), markerPrefix),
			},
			Price: price,
		})
		it.Next()
	}
	return nftList
}

func (k Keeper) GetNFTIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte(markerPrefix))
}

const markerPrefix = "market_"
const ownerSuffix = "_owner"

func composePutNFTToMarketKey(tokenID string) []byte {
	return []byte(markerPrefix + tokenID)
}

func composePutNFTOwnerKey(tokenID string) []byte {
	return []byte(tokenID + ownerSuffix)
}
