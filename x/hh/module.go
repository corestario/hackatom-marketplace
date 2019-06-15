package hh

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	abci "github.com/tendermint/tendermint/abci/types"
)

// type check to ensure the interface is properly implemented
var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

const ModuleName = "hh"

// app module Basics object
type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
	return ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// Validation check of the Genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	// once json successfully marshalled, passes along to genesis.go
	return ValidateGenesis(data)
}

type AppModule struct {
	AppModuleBasic
	keeper     Keeper
	coinKeeper bank.Keeper
}

func (am AppModule) RegisterRESTRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// GetNFTokenData(TokenID) TokenData // Получить информацию о токене
	r.HandleFunc(fmt.Sprintf("/%s/nft/{%s}", ModuleName, restName), getNFTHandler(cliCtx.Codec, cliCtx, ModuleName)).Methods("GET")
	// GetNFTokensOnSaleList() []TokenData // Возвращает список продающихся токенов с ценами
	r.HandleFunc(fmt.Sprintf("/%s/nft/list/{%s}/", ModuleName, restName), getNFTOnSaleListHandler(cliCtx.Codec, cliCtx, ModuleName)).Methods("GET")

	// TransferNFTokenToZone(ZoneID, TokenID) TransferID // Передаёт токен на соседнуюю зону (напр. зону выпуска токенов), но не выставляет на продажу
	r.HandleFunc(fmt.Sprintf("/%s/nft/transfer", ModuleName), transferNFTokenToZone(cliCtx.Codec, cliCtx)).Methods("POST")

	// GetTransferStatus(TransferID) Status возвращает статус трансфера - в процессе, прилетел, ошибка
	r.HandleFunc(fmt.Sprintf("/%s/nft/transfer/{%s}", ModuleName, restName), getTransferStatus(cliCtx.Codec, cliCtx, ModuleName)).Methods("GET")

	// BuyNFToken(TokenID) Status // Меняет владельца токена, меняет статус токена на непродаваемый, переводит деньги (с комиссией) бывшему владельцу токена
	r.HandleFunc(fmt.Sprintf("/%s/nft/buy", ModuleName), buyNFToken(cliCtx.Codec, cliCtx)).Methods("POST")
	// PutNFTokenOnTheMarket(TokenID, Price) Status // Меняет статус токена на продаваемый, устанавливает цену
	r.HandleFunc(fmt.Sprintf("/%s/nft/sell", ModuleName), putNFTokenOnTheMarket(cliCtx.Codec, cliCtx)).Methods("POST")
}

func (am AppModule) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	hhTxCmd := &cobra.Command{
		Use:   "hh",
		Short: "hh transactions subcommands",
	}

	hhTxCmd.AddCommand(client.PostCommands(
		GetCmdTransferToken(cdc),
	)...)

	return hhTxCmd
}

func (am AppModule) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	hhQueryCmd := &cobra.Command{
		Use:   "hh",
		Short: "Querying commands for the hh module",
	}
	hhQueryCmd.AddCommand(client.GetCommands(
		GetCmdTokenInfo(ModuleName, cdc),
		GetCmdListTokens(ModuleName, cdc),
		GetCmdTransferInfo(ModuleName, cdc),
	)...)

	return hhQueryCmd
}

func (AppModuleBasic) RegisterRESTRoutes(cliCtx context.CLIContext, r *mux.Router) {
}

func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

func (AppModuleBasic) GetQueryCmd(*codec.Codec) *cobra.Command {
	//panic("implement me")
	return nil
}

// NewAppModule creates a new AppModule Object
func NewAppModule(k Keeper, bankKeeper bank.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
		coinKeeper:     bankKeeper,
	}
}

func (AppModule) Name() string {
	return ModuleName
}

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

func (am AppModule) Route() string {
	return RouterKey
}

func (am AppModule) NewHandler() types.Handler {
	return NewHandler(am.keeper)
}
func (am AppModule) QuerierRoute() string {
	return ModuleName
}

func (am AppModule) NewQuerierHandler() types.Querier {
	return NewQuerier(am.keeper)
}

func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) types.Tags {
	return sdk.EmptyTags()
}

func (am AppModule) EndBlock(types.Context, abci.RequestEndBlock) ([]abci.ValidatorUpdate, types.Tags) {
	return []abci.ValidatorUpdate{}, sdk.EmptyTags()
}

func (am AppModule) InitGenesis(ctx types.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, am.keeper, genesisState)
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}
