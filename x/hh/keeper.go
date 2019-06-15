package hh

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/ibc"
)

// StoreKey to be used when creating the KVStore
const StoreKey = "hh"

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper
	ibcKeeper  ibc.BankKeeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the hh Keeper
func NewKeeper(coinKeeper bank.Keeper, ibcKeeper ibc.BankKeeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		ibcKeeper:  ibcKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

func (k Keeper) TransferNFTokenToZone(ctx sdk.Context, nfToken NFT, zoneID string, sender, recipient sdk.AccAddress) error {
	if !k.getNFTOwner(ctx, nfToken.ID).Empty() {
		return errors.New("call from not the owner")
	}

	//fixme call transfetToZone
	//k.ibcKeeper.TransferNFTokenToZone(nft, msg.ZoneID, msg.Sender, msg.Recipient)

	k.DeleteNFT(ctx, nfToken.ID)
	return nil
}

func (k Keeper) PutNFTokenOnTheMarket(ctx sdk.Context, token NFT, sender sdk.AccAddress) error {
	token.OnSale = true

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

	nftOnSaleBin := k.cdc.MustMarshalBinaryBare(token)

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
	store := ctx.KVStore(k.storeKey)
	if buyer.Empty() {
		return errors.New("empty buyer")
	}

	if store.Has(composePutNFTToMarketKey(nfTokenID)) {
		return errors.New("nft has already existed on market")
	}

	tokenOwner := k.getNFTOwner(ctx, nfTokenID)
	if tokenOwner.Equals(buyer) {
		return errors.New("you are the owner of the nft")
	}

	//get token
	nftBin := store.Get(composePutNFTToMarketKey(nfTokenID))

	nftOnSale := new(NFT)
	err := k.cdc.UnmarshalBinaryBare(nftBin, nftOnSale)
	if err != nil {
		ctx.Logger().Error("error while GetNFToken ", "tokenID", nfTokenID)
		return err
	}

	if !nftOnSale.OnSale {
		return errors.New("you are the owner of the nft")
	}

	//compare money and price
	buyerCoins := k.coinKeeper.GetCoins(ctx, buyer)
	for _, coin := range nftOnSale.Price {
		buyerCoinValue := buyerCoins.AmountOf(coin.Denom)
		if coin.Amount.GT(buyerCoinValue) {
			return errors.New("you dont have enough coins to buy nft")
		}
	}

	//buy
	errResult := k.coinKeeper.SendCoins(ctx, tokenOwner, buyer, nftOnSale.Price)
	if errResult != nil {
		return errResult
	}

	//change owner
	nftOnSale.OnSale = false
	nftOnSale.Owner = buyer

	//store
	nftOnSaleBin := k.cdc.MustMarshalBinaryBare(nftOnSale)
	store.Set(composePutNFTToMarketKey(nfTokenID), nftOnSaleBin)
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

func (k Keeper) DeleteNFT(ctx sdk.Context, tokenID string) {
	store := ctx.KVStore(k.storeKey)

	if !store.Has([]byte(tokenID)) {
		return
	}

	store.Delete([]byte(tokenID))
	store.Delete(composePutNFTOwnerKey(tokenID))
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

func (k Keeper) GetNFTokens(ctx sdk.Context) []NFT {
	it := k.GetNFTIterator(ctx)
	var nftList []NFT
	for {
		if it.Valid() == false {
			break
		}

		var nftoken NFT
		k.cdc.MustUnmarshalBinaryBare(it.Value(), &nftoken)

		fmt.Println(nftoken)

		nftList = append(nftList, nftoken)

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
