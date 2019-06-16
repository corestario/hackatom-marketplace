package hh

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/ibc/02-client/tendermint"
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
	ibck.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	return cdc
}

func makeAcc() sdk.AccAddress {
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])
	return sdk.AccAddress(pub.Address())
}

func TestPutTwoNFTOnMarket(t *testing.T) {
	stKey := sdk.NewKVStoreKey(StoreKey)
	ti := setupTestInput(stKey, "1")
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
	stKey := sdk.NewKVStoreKey(StoreKey)
	ti := setupTestInput(stKey,"1")
	k := NewKeeper(nil, ibck.Keeper{}, stKey, ti.cdc)

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

func TestIBC(t *testing.T)  {
	stKey := sdk.NewKVStoreKey(StoreKey)
	ti1 := setupTestInput(stKey,"ch1")
	ti2 := setupTestInput(stKey,"ch2")
	ibcKeeper1:=ibck.NewKeeper(ti1.cdc,stKey)
	ibcKeeper2:=ibck.NewKeeper(ti2.cdc,stKey)


	clientID1:="clientID1"
	chainID1:="chainID1"


	connID:="some conn"
	cp1:="cp1"
	cp2:="cp2"
	id:="123"



	var err error
	err = ibcKeeper1.CreateClient(ti1.ctx,clientID1,tendermint.ConsensusState{
		ChainID:chainID1,
	})
	if err!=nil {
		t.Fatal(err)
	}

	err=ibcKeeper1.OpenConnection(ti1.ctx,clientID1, cp1, clientID1, cp2)
	if err!=nil {
		t.Fatal(err)
	}
	err=ibcKeeper1.OpenChannel(ti1.ctx, ModuleName, connID, id, cp1, cp2)
	if err!=nil {
		t.Fatal(err)
	}








	k1 := NewKeeper(nil, ibcKeeper1, stKey, ti1.cdc)
	k2 := NewKeeper(nil, ibcKeeper2, stKey, ti2.cdc)
	_=k2
	acc1:=makeAcc()
	acc2:=makeAcc()




	err=k1.TransferNFTokenToZone(
		ti1.ctx,
		NFT{
			BaseNFT:BaseNFT{
					ID:"one",
				},
		},
		"zone1",
		acc1,
		acc2,
		)
	if err!=nil {
		t.Fatal(err)
	}



}

type testInput struct {
	cdc *codec.Codec
	ctx sdk.Context

	stKey *sdk.KVStoreKey

	auth auth.AccountKeeper
	bank bank.Keeper
}

func setupTestInput(key sdk.StoreKey, chainID string) testInput {
	db := dbm.NewMemDB()

	cdc := MakeCodec()
	key2 := sdk.NewKVStoreKey("test")

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
	ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(key2, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"+randomSuffix}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(
		cdc, authCapKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount,
	)
	ak.SetParams(ctx, auth.DefaultParams())
	ctx := sdk.NewContext(ms, abci.Header{ChainID: chainID}, false, log.NewNopLogger())

	bankKeeper := bank.NewBaseKeeper(ak, pk.Subspace(types.DefaultParamspace), types.DefaultCodespace)
	bankKeeper.SetSendEnabled(ctx, true)

	return testInput{cdc: cdc, ctx: ctx, stKey: stKey, auth: ak, bank: bankKeeper}
}