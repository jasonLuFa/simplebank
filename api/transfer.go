package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/jasonLuFa/simplebank/db/sqlc"
	"github.com/jasonLuFa/simplebank/token"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        string `json:"amount" binding:"required,validAmount"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, valid := server.validateAccountCurrency(ctx,req.FromAccountID,req.Currency)
	if !valid {
		return 
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username{
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}
	
	_, valid = server.validateAccountCurrency(ctx,req.ToAccountID,req.Currency)
	if !valid  {
		return 
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID: req.ToAccountID,
		Amount: req.Amount,
	}

	result, err := server.store.TransferTx(ctx,arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,result)
}

func (server *Server) validateAccountCurrency(ctx *gin.Context, accountID int64, currency string) (db.Account,bool) {
	account, err := server.store.GetAccount(ctx,accountID)
	if err != nil {
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound,errorResponse(err))
			return account,false
		}

		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
	}

	if account.Currency != currency {
		err = fmt.Errorf("account [%d] currency mismatch: db -> %s vs request -> %s",accountID,account.Currency,currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account,false
	}

	return account,true
}