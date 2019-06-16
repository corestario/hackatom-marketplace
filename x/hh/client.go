package hh

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func GetCmdTokenInfo(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "token-info [id]",
		Short: "See token data by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			tokenID := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/NFToken/%s", queryRoute, tokenID), nil)
			if err != nil {
				fmt.Printf("could not find tokenID - %s: %v\n", tokenID, err)
				return nil
			}

			var out NFT
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdListTokens(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "token-list",
		Short: "See token list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/NFTokens", queryRoute), nil)
			if err != nil {
				fmt.Printf("could not get token list: %v", err)
				return nil
			}

			var out QueryResNFTokens
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdTransferInfo(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "token-transfer [id]",
		Short: "token transfer by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			transferID := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/Transfer/%s", queryRoute, transferID), nil)
			if err != nil {
				fmt.Printf("could not find tokenID - %s: %v\n", transferID, err)
				return nil
			}

			var out Transfer
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

func GetCmdTransferToken(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "transfer-token [tokenID] [zoneID]",
		Short: "bid for existing name or claim new name",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			msg := NewMsgTransferTokenToZone(cliCtx.GetFromAddress(), args[0], args[2])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

const (
	restName = "name"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec, storeName string) {
}

func getNFTHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/NFToken/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getNFTOnSaleListHandler(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/NFTokens/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func getTransferStatus(cdc *codec.Codec, cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		paramType := vars[restName]

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/Transfer/%s", storeName, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

type putOnMarketNFTReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Owner   string       `json:"owner"`
	Token   BaseNFT      `json:"token"`
	Price   sdk.Coin     `json:"price"`

	// User data
	Name     string `json:"name"`
	Password string `json:"password"`
}

func putNFTokenOnTheMarket(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req putOnMarketNFTReq
		//err:=json.NewDecoder(r.Body).Decode(&req)
		//fmt.Println("err",err)

		//fmt.Println("nftToken",nftToken)

		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		fmt.Println("---------------------", 2)
		fmt.Println("reqRef", req)
		baseReq := req.BaseReq.Sanitize()

		fmt.Println("[[[[[[[[[[[[[reqRef2", req)
		fmt.Println("[[[[[[[[[[[[[baseReq", baseReq)
		fmt.Println("reqRef---------------", baseReq.ChainID)
		if !baseReq.ValidateBasic(w) {
			return
		}

		fmt.Println("---------------------", 3)

		addr, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fmt.Println("---------------------", 4)

		nftToken := NFT{req.Token, true, sdk.Coins{req.Price}}
		// create the message
		msg := NewMsgPutNFTokenOnTheMarket(nftToken, addr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fmt.Println("---------------------", 5)

		gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, baseReq.GasAdjustment, flags.DefaultGasAdjustment)
		if !ok {
			return
		}

		_, gas, err := flags.ParseGas(baseReq.Gas)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fmt.Println("---------------------", 6)

		txBldr := authtxb.NewTxBuilder(
			utils.GetTxEncoder(cdc), baseReq.AccountNumber, baseReq.Sequence, gas, gasAdj,
			baseReq.Simulate, baseReq.ChainID, baseReq.Memo, baseReq.Fees, baseReq.GasPrices,
		)

		msgBytes, err := txBldr.BuildAndSign(req.Name, req.Password, []sdk.Msg{msg})
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fmt.Println("---------------------", 7)

		_, err = cliCtx.BroadcastTxCommit(msgBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, true)
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
	Price     string       `json:"price"`

	// User data
	Name     string `json:"name"`
	Password string `json:"password"`
}

func transferNFTokenToZone(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req transferNFTReq

		runPostFunction(w, r, cdc, cliCtx, req.BaseReq, &req, req.Name, req.Password, req.Owner, func(sender sdk.AccAddress) sdk.Msg {
			return NewMsgTransferTokenToZone(sender, req.NFTokenID, req.ZoneID)
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

		fmt.Println("---------------------", 1)
		if !rest.ReadRESTReq(w, r, cdc, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		fmt.Println("---------------------", 2)
		fmt.Println("reqRef", req)
		baseReq := req.BaseReq.Sanitize()

		fmt.Println("[[[[[[[[[[[[[reqRef2", req)
		fmt.Println("reqRef---------------", baseReq.ChainID)
		if !baseReq.ValidateBasic(w) {
			return
		}

		fmt.Println("---------------------", 3)

		addr, err := sdk.AccAddressFromBech32(req.Owner)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fmt.Println("---------------------", 4)

		// create the message
		msg := NewMsgBuyNFToken(req.NFTokenID, priceInCoins, addr)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fmt.Println("---------------------", 5)

		gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, baseReq.GasAdjustment, flags.DefaultGasAdjustment)
		if !ok {
			return
		}

		_, gas, err := flags.ParseGas(baseReq.Gas)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fmt.Println("---------------------", 6)

		txBldr := authtxb.NewTxBuilder(
			utils.GetTxEncoder(cdc), baseReq.AccountNumber, baseReq.Sequence, gas, gasAdj,
			baseReq.Simulate, baseReq.ChainID, baseReq.Memo, baseReq.Fees, baseReq.GasPrices,
		)

		msgBytes, err := txBldr.BuildAndSign(req.Name, req.Password, []sdk.Msg{msg})
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		fmt.Println("---------------------", 7)

		_, err = cliCtx.BroadcastTxCommit(msgBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, true)
	}
}

func runPostFunction(w http.ResponseWriter, r *http.Request, cdc *codec.Codec,
	cliCtx context.CLIContext, baseReq rest.BaseReq, reqRef interface{},
	name, password, owner string,
	postFunc func(address sdk.AccAddress) sdk.Msg) {

	fmt.Println("---------------------", 1)
	if !rest.ReadRESTReq(w, r, cdc, reqRef) {
		rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
		return
	}

	fmt.Println("---------------------", 2)
	fmt.Println("reqRef", reqRef)
	baseReq = baseReq.Sanitize()

	fmt.Println("[[[[[[[[[[[[[reqRef2", reqRef)
	fmt.Println("reqRef---------------", baseReq.ChainID)
	if !baseReq.ValidateBasic(w) {
		return
	}

	fmt.Println("---------------------", 3)

	addr, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("---------------------", 4)

	// create the message
	msg := postFunc(addr)
	err = msg.ValidateBasic()
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("---------------------", 5)

	gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, baseReq.GasAdjustment, flags.DefaultGasAdjustment)
	if !ok {
		return
	}

	_, gas, err := flags.ParseGas(baseReq.Gas)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("---------------------", 6)

	txBldr := authtxb.NewTxBuilder(
		utils.GetTxEncoder(cdc), baseReq.AccountNumber, baseReq.Sequence, gas, gasAdj,
		baseReq.Simulate, baseReq.ChainID, baseReq.Memo, baseReq.Fees, baseReq.GasPrices,
	)

	msgBytes, err := txBldr.BuildAndSign(name, password, []sdk.Msg{msg})
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("---------------------", 7)

	_, err = cliCtx.BroadcastTxCommit(msgBytes)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	rest.PostProcessResponse(w, cliCtx, true)
}
