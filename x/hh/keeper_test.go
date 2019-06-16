package hh

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/02-client/tendermint"
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
	ti := setupTestInput()

	k := ti.keeper

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
	k := ti.keeper

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
	k := ti.keeper

	sellerAccount := makeAcc()
	acc := k.accountKeeper.NewAccountWithAddress(ti.ctx, sellerAccount)
	k.accountKeeper.SetAccount(ti.ctx, acc)

	if !k.coinKeeper.GetCoins(ti.ctx, sellerAccount).IsEqual(sdk.NewCoins()) {
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
			ID:          "1234",
			Owner:       sellerAccount,
			Name:        "dog",
			Description: "a wet dog",
			Image:       "some.gif",
			TokenURI:    ".ws",
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
	accBuyer := k.accountKeeper.NewAccountWithAddress(ti.ctx, buyerAccount)
	k.accountKeeper.SetAccount(ti.ctx, accBuyer)

	initialCoin := sdk.NewInt64Coin(denomination, 10000)
	initialCoins := sdk.NewCoins(initialCoin)
	k.coinKeeper.SetCoins(ti.ctx, buyerAccount, initialCoins)

	if !k.coinKeeper.GetCoins(ti.ctx, buyerAccount).IsEqual(initialCoins) {
		t.Fatal("sellerAccount should have", initialCoins.String())
	}

	err = k.BuyNFToken(ti.ctx, nftToSell.BaseNFT.ID, buyerAccount)
	if err != nil {
		t.Fatal(err)
	}

	if !k.coinKeeper.GetCoins(ti.ctx, sellerAccount).IsEqual(price) {
		t.Fatal("sellerAccount should get nft price")
	}
	if !k.coinKeeper.GetCoins(ti.ctx, buyerAccount).IsEqual(sdk.NewCoins(initialCoin.Sub(priceCoin))) {
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

func TestIBC(t *testing.T) {
	ti1 := setupTestInput()
	ti2 := setupTestInput()

	clientID1 := "clientID1"

	connID := "some conn"
	cp1 := "cp1"
	cp2 := "cp2"
	id := "123"

	var err error
	err = ti1.keeper.ibcKeeper.CreateClient(ti1.ctx, clientID1, tendermint.ConsensusState{
		ChainID: ti1.chainID,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = ti1.keeper.ibcKeeper.OpenConnection(ti1.ctx, connID, cp1, clientID1, cp2)
	if err != nil {
		t.Fatal(err)
	}
	err = ti1.keeper.ibcKeeper.OpenChannel(ti1.ctx, ModuleName, connID, id, cp1, cp2)
	if err != nil {
		t.Fatal(err)
	}

	//acc1:=makeAcc()
	//acc2:=makeAcc()

	packet := NewSendTokenPacket(&BaseNFT{
		ID: "one",
	})

	err = ti1.keeper.ibcKeeper.Send(ti1.ctx, connID, id, packet)
	if err != nil {
		t.Fatal(err)
	}

	err = ti2.keeper.ibcKeeper.CreateClient(ti2.ctx, clientID1, tendermint.ConsensusState{
		ChainID: ti2.chainID,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = ti2.keeper.ibcKeeper.OpenConnection(ti2.ctx, connID, cp2, clientID1, cp2)
	if err != nil {
		t.Fatal(err)
	}
	err = ti2.keeper.ibcKeeper.OpenChannel(ti2.ctx, ModuleName, connID, id, cp2, cp1)
	if err != nil {
		t.Fatal(err)
	}

	pkt := &SendTokenPacket{}

	err = ti2.keeper.ibcKeeper.Receive(ti2.ctx, connID, id, pkt, nil)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(pkt)

}

type testInput struct {
	cdc *codec.Codec
	ctx sdk.Context

	stKey *sdk.KVStoreKey

	keeper Keeper

	chainID string
}

func setupTestInput() testInput {
	db := dbm.NewMemDB()

	cdc := MakeCodec()

	var randomSuffixBytes [16]byte
	rand.Read(randomSuffixBytes[:])

	randomSuffix := string(randomSuffixBytes[:])

	authCapKey := sdk.NewKVStoreKey("authCapKey" + randomSuffix)
	fckCapKey := sdk.NewKVStoreKey("fckCapKey" + randomSuffix)
	stKey := sdk.NewKVStoreKey("storeKey" + randomSuffix)
	ibcKey := sdk.NewKVStoreKey("ibckey" + randomSuffix)
	feeKey := sdk.NewKVStoreKey("feekey" + randomSuffix)
	storeKey := sdk.NewKVStoreKey("storeKeyKeeper" + randomSuffix)
	keyParams := sdk.NewKVStoreKey("params" + randomSuffix)
	tkeyParams := sdk.NewTransientStoreKey("transient_params" + randomSuffix)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(feeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(stKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	chainID := "test-chain-id" + randomSuffix

	ctx := sdk.NewContext(ms, abci.Header{ChainID: chainID}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(
		cdc, authCapKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount,
	)
	ak.SetParams(ctx, auth.DefaultParams())

	bankKeeper := bank.NewBaseKeeper(ak, pk.Subspace(types.DefaultParamspace), types.DefaultCodespace)
	bankKeeper.SetSendEnabled(ctx, true)

	ibcKeeper := ibck.NewKeeper(cdc, ibcKey)

	feeCollectionKeeper := auth.NewFeeCollectionKeeper(cdc, feeKey)

	keeper := NewKeeper(bankKeeper,
		ibcKeeper,
		ak,
		feeCollectionKeeper,
		storeKey,
		cdc)

	return testInput{cdc: cdc, ctx: ctx, stKey: stKey, keeper: keeper, chainID: chainID}
}
