package rest

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"net/http"

	"dgamingfoundation/hackathon-hub/x/hh"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	txbuilder "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"

	"github.com/gorilla/mux"
)

const (
	restName = "name"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
	// GetNFTokenData(TokenID) TokenData // Получить информацию о токене
	r.HandleFunc(fmt.Sprintf("/%s/nft/{%s}", storeName, restName), getNFTHandler(cdc, cliCtx, storeName)).Methods("GET")
	// GetNFTokensOnSaleList() []TokenData // Возвращает список продающихся токенов с ценами
	r.HandleFunc(fmt.Sprintf("/%s/nft/list/{%s}/", storeName, restName), getNFTOnSaleListHandler(cdc, cliCtx, storeName)).Methods("GET")

	// TransferNFTokenToZone(ZoneID, TokenID) TransferID // Передаёт токен на соседнуюю зону (напр. зону выпуска токенов), но не выставляет на продажу
	r.HandleFunc(fmt.Sprintf("/%s/nft/transfer", storeName), transferNFTokenToZone(cdc, cliCtx)).Methods("POST")

	// GetTransferStatus(TransferID) Status возвращает статус трансфера - в процессе, прилетел, ошибка
	r.HandleFunc(fmt.Sprintf("/%s/nft/transfer/{%s}", storeName, restName), getTransferStatus(cdc, cliCtx, storeName)).Methods("GET")

	// BuyNFToken(TokenID) Status // Меняет владельца токена, меняет статус токена на непродаваемый, переводит деньги (с комиссией) бывшему владельцу токена
	r.HandleFunc(fmt.Sprintf("/%s/nft/buy", storeName), buyNFToken(cdc, cliCtx)).Methods("POST")
	// PutNFTokenOnTheMarket(TokenID, Price) Status // Меняет статус токена на продаваемый, устанавливает цену
	r.HandleFunc(fmt.Sprintf("/%s/nft/sell", storeName), putNFTokenOnTheMarket(cdc, cliCtx)).Methods("POST")
}

func getNFTHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/NFToken/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func getNFTOnSaleListHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/NFTokens/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

func getTransferStatus(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/Transfer/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cdc, res, cliCtx.Indent)
	}
}

type putOnMarketNFTReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Owner   string       `json:"owner"`
	Token   hh.BaseNFT   `json:"token"`
	Price   string       `json:"price"`

	// User data
	Name     string `json:"name"`
	Password string `json:"password"`
}

func putNFTokenOnTheMarket(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req putOnMarketNFTReq

		priceInCoins, err := sdk.ParseCoins(req.Price)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		nftToken := hh.NFT{req.Token, true, priceInCoins}

		runPostFunction(w, r, cdc, cliCtx, req.BaseReq, &req, req.Name, req.Password, req.Owner, func(addr sdk.AccAddress) sdk.Msg {
			return hh.NewMsgPutNFTokenOnTheMarket(nftToken, addr)
		})
	}
}

type transferNFTReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Owner     string       `json:"owner"`
	Recipient string       `json:"recipient"`
	NFTokenID string       `json:"token_id"`
	ZoneID    string       `json:"zone_id"`

	// User data
	Name     string `json:"name"`
	Password string `json:"password"`
}

type buyNFTReq struct {
	BaseReq   rest.BaseReq `json:"base_req"`
	Owner     string       `json:"owner"`
	NFTokenID string       `json:"token_id"`
	Price   string       `json:"price"`

	// User data
	Name     string `json:"name"`
	Password string `json:"password"`
}

func transferNFTokenToZone(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req transferNFTReq

		recipient, err := sdk.AccAddressFromBech32(req.Recipient)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		runPostFunction(w, r, cdc, cliCtx, req.BaseReq, &req, req.Name, req.Password, req.Owner, func(sender sdk.AccAddress) sdk.Msg {
			return hh.NewMsgTransferNFTokenToZone(req.NFTokenID, req.ZoneID, sender, recipient)
		})
	}
}

func buyNFToken(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req buyNFTReq

		priceInCoins, err := sdk.ParseCoins(req.Price)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		runPostFunction(w, r, cdc, cliCtx, req.BaseReq, &req, req.Name, req.Password, req.Owner, func(addr sdk.AccAddress) sdk.Msg {
			return hh.NewMsgBuyNFToken(req.NFTokenID, priceInCoins, addr)
		})
	}
}

func runPostFunction(w http.ResponseWriter, r *http.Request, cdc *codec.Codec,
	cliCtx context.CLIContext, baseReq rest.BaseReq, reqRef interface{},
	name, password, owner string,
	postFunc func(address sdk.AccAddress) sdk.Msg) {

	if !rest.ReadRESTReq(w, r, cdc, reqRef) {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
		return
	}

	baseReq = baseReq.Sanitize()
	if !baseReq.ValidateBasic(w) {
		return
	}

	addr, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// create the message
	msg := postFunc(addr)
	err = msg.ValidateBasic()
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, baseReq.GasAdjustment, flags.DefaultGasAdjustment)
	if !ok {
		return
	}

	_, gas, err := flags.ParseGas(baseReq.Gas)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	txBldr := txbuilder.NewTxBuilder(
		utils.GetTxEncoder(cdc), baseReq.AccountNumber, baseReq.Sequence, gas, gasAdj,
		baseReq.Simulate, baseReq.ChainID, baseReq.Memo, baseReq.Fees, baseReq.GasPrices,
	)

	msgBytes, err := txBldr.BuildAndSign(name, password, []sdk.Msg{msg})
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = cliCtx.BroadcastTxCommit(msgBytes)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	rest.PostProcessResponse(w, cdc, http.StatusOK, true)
}
