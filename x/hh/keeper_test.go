package hh

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"math/rand"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibck "github.com/cosmos/cosmos-sdk/x/ibc/keeper"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/cosmos/cosmos-sdk/types/module"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

var ModuleBasics module.BasicManager

const denomination = "usd"

func init() {

	ModuleBasics = module.NewBasicManager(
		AppModule{},
	)
}

func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	ibck.RegisterCodec(cdc)

	ModuleBasics.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	return cdc
}

func makeAcc() sdk.AccAddress {
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])
	return sdk.AccAddress(pub.Address())
}

func TestPutTwoNFTOnMarket(t *testing.T) {
	ti := setupTestInput()
	k := NewKeeper(ti.bank, ibck.Keeper{}, ti.auth, auth.FeeCollectionKeeper{}, ti.stKey, ti.cdc)

	account := makeAcc()
	price := sdk.Coins{sdk.Coin{
		denomination,
		sdk.NewInt(100),
	}}
	someToken := NFT{
		BaseNFT{
			ID: "1234",
		},
		false,
		price,
	}
	k.setNFTOwner(ti.ctx, someToken.BaseNFT.ID, account)
	//put first NFT
	err := k.PutNFTokenOnTheMarket(ti.ctx, someToken, account)
	if err != nil {
		t.Fatal(err)
	}
	nftList := k.GetNFTokens(ti.ctx)
	if len(nftList) != 1 {
		t.Fatal("incorrect length")
	}

	newToken := someToken
	newToken.ID = newToken.ID + "1"
	k.setNFTOwner(ti.ctx, newToken.ID, account)
	err = k.PutNFTokenOnTheMarket(ti.ctx, newToken, account)
	if err != nil {
		t.Fatal(err)
	}
	nftList = k.GetNFTokens(ti.ctx)
	if len(nftList) != 2 {
		t.Fatal("incorrect length")
	}
}

func TestPutSameNFTOnMarket(t *testing.T) {
	ti := setupTestInput()
	k := NewKeeper(ti.bank, ibck.Keeper{}, ti.auth, auth.FeeCollectionKeeper{}, ti.stKey, ti.cdc)

	price := sdk.Coins{sdk.Coin{
		denomination,
		sdk.NewInt(100),
	}}
	someToken := NFT{
		BaseNFT{
			ID: "1234",
		},
		false,
		price,
	}
	account := makeAcc()
	k.setNFTOwner(ti.ctx, someToken.ID, account)
	//put first NFT
	err := k.PutNFTokenOnTheMarket(ti.ctx, someToken, account)
	if err != nil {
		t.Fatal(err)
	}

	nftList := k.GetNFTokens(ti.ctx)
	if len(nftList) != 1 {
		t.Fatal("incorrect length")
	}

	err = k.PutNFTokenOnTheMarket(ti.ctx, someToken, sdk.AccAddress{})
	if err == nil {
		t.FailNow()
	}

	nftList = k.GetNFTokens(ti.ctx)
	if len(nftList) != 1 {
		t.Fatal("incorrect length")
	}
}

func TestPutAndBuyNFT(t *testing.T) {
	ti := setupTestInput()
	k := NewKeeper(ti.bank, ibck.Keeper{}, ti.auth, auth.FeeCollectionKeeper{}, ti.stKey, ti.cdc)

	sellerAccount := makeAcc()
	acc := ti.auth.NewAccountWithAddress(ti.ctx, sellerAccount)
	ti.auth.SetAccount(ti.ctx, acc)
	
	if !ti.bank.GetCoins(ti.ctx, sellerAccount).IsEqual(sdk.NewCoins()) {
		t.Fatal("sellerAccount should be empty")
	}

	nftList := k.GetNFTokens(ti.ctx)
	if len(nftList) != 0 {
		t.Fatal("nft list should be empty")
	}

	priceCoin := sdk.Coin{denomination, sdk.NewInt(100)}
	price := sdk.Coins{priceCoin}
	
	nftToSell := NFT{
		BaseNFT{
			ID: "1234",
			Owner: sellerAccount,
			Name: "dog",
			Description: "a wet dog",
			Image: "some.gif",
			TokenURI: ".ws",
		},
		false,
		price,
	}
	
	
	k.setNFTOwner(ti.ctx, nftToSell.BaseNFT.ID, sellerAccount)

	//put first NFT
	err := k.PutNFTokenOnTheMarket(ti.ctx, nftToSell, sellerAccount)
	if err != nil {
		t.Fatal(err)
	}

	nftList = k.GetNFTokens(ti.ctx)
	if len(nftList) != 1 {
		t.Fatal("incorrect length")
	}

	storedNFT := nftList[0]
	if storedNFT.OnSale != true {
		t.Fatal("the stored token has wrong OnSale")
	}
	if storedNFT.Owner.String() != sellerAccount.String() {
		t.Fatal("the stored token has wrong Owner")
	}
	if !storedNFT.Price.IsEqual(price) {
		t.Fatal("the stored token has wrong Price")
	}
	if storedNFT.ID != nftToSell.ID {
		t.Fatal("the stored token has wrong ID")
	}
	if storedNFT.Name != nftToSell.Name {
		t.Fatal("the stored token has wrong Name")
	}
	if storedNFT.Description != nftToSell.Description {
		t.Fatal("the stored token has wrong Description")
	}
	if storedNFT.TokenURI != nftToSell.TokenURI {
		t.Fatal("the stored token has wrong TokenURI")
	}

	buyerAccount := makeAcc()
	accBuyer := ti.auth.NewAccountWithAddress(ti.ctx, buyerAccount)
	ti.auth.SetAccount(ti.ctx, accBuyer)

	initialCoin := sdk.NewInt64Coin(denomination, 10000)
	initialCoins := sdk.NewCoins(initialCoin)
	ti.bank.SetCoins(ti.ctx, buyerAccount, initialCoins)

	if !ti.bank.GetCoins(ti.ctx, buyerAccount).IsEqual(initialCoins) {
		t.Fatal("sellerAccount should have", initialCoins.String())
	}

	err = k.BuyNFToken(ti.ctx, nftToSell.BaseNFT.ID, buyerAccount)
	if err != nil {
		t.Fatal(err)
	}


	if !ti.bank.GetCoins(ti.ctx, sellerAccount).IsEqual(price) {
		t.Fatal("sellerAccount should get nft price")
	}
	if !ti.bank.GetCoins(ti.ctx, buyerAccount).IsEqual(sdk.NewCoins(initialCoin.Sub(priceCoin))) {
		t.Fatal("buyerAccount should be lost nft price")
	}

	nftList = k.GetNFTokens(ti.ctx)
	if len(nftList) != 1 {
		t.Fatal("incorrect length")
	}
	storedNFT = nftList[0]

	if storedNFT.OnSale != false {
		t.Fatal("the stored token has wrong OnSale")
	}
	if storedNFT.Owner.String() != buyerAccount.String() {
		t.Fatal("the stored token has wrong Owner")
	}
	if !storedNFT.Price.IsEqual(price) {
		t.Fatal("the stored token has wrong Price")
	}
	if storedNFT.ID != nftToSell.ID {
		t.Fatal("the stored token has wrong ID")
	}
	if storedNFT.Name != nftToSell.Name {
		t.Fatal("the stored token has wrong Name")
	}
	if storedNFT.Description != nftToSell.Description {
		t.Fatal("the stored token has wrong Description")
	}
	if storedNFT.TokenURI != nftToSell.TokenURI {
		t.Fatal("the stored token has wrong TokenURI")
	}
}

type testInput struct {
	cdc *codec.Codec
	ctx sdk.Context

	stKey *sdk.KVStoreKey

	auth auth.AccountKeeper
	bank bank.Keeper
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()

	cdc := MakeCodec()

	var randomSuffixBytes [16]byte
	rand.Read(randomSuffixBytes[:])

	randomSuffix := string(randomSuffixBytes[:])

	authCapKey := sdk.NewKVStoreKey("authCapKey"+randomSuffix)
	fckCapKey := sdk.NewKVStoreKey("fckCapKey"+randomSuffix)
	stKey := sdk.NewKVStoreKey("storeKey"+randomSuffix)
	keyParams := sdk.NewKVStoreKey("params"+randomSuffix)
	tkeyParams := sdk.NewTransientStoreKey("transient_params"+randomSuffix)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(stKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"+randomSuffix}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(
		cdc, authCapKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount,
	)
	ak.SetParams(ctx, auth.DefaultParams())

	bankKeeper := bank.NewBaseKeeper(ak, pk.Subspace(types.DefaultParamspace), types.DefaultCodespace)
	bankKeeper.SetSendEnabled(ctx, true)

	return testInput{cdc: cdc, ctx: ctx, stKey: stKey, auth: ak, bank: bankKeeper}
}