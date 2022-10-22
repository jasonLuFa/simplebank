package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	db "github.com/jasonLuFa/simplebank/db/sqlc"
)

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context){
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner: req.Owner,
		Currency: req.Currency,
		Balance: "0",
	}

	account,err:= server.store.CreateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,account)
}

type getAccountRequest struct{
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context){
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	account,err:= server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err ==sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound,errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,account)
}

type listAccountsRequest struct{
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req listAccountsRequest
	if err:= ctx.ShouldBindQuery(&req); err != nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit: req.PageSize, 
		Offset: (req.PageID -1) * req.PageSize,
	}

	account,err := server.store.ListAccounts(ctx,arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,account)

}

type deleteAccountRequest struct{
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(ctx *gin.Context) {
	var req deleteAccountRequest
	if err:= ctx.ShouldBindUri(&req); err != nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	err := server.store.DeleteAccount(ctx,req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,gin.H{
		"success": "isdelete",
	})
}

// 這邊 balance 因為是 string ，所以要另外判斷是否轉成換換成 float64 時，是否為 < 0
type updateAccountRequest struct{
	Balance string `json:"balance" binding:"required"`
}

func (server *Server) updateAccount(ctx *gin.Context) {
	var req updateAccountRequest
	id,err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	if err:= ctx.ShouldBindJSON(&req); err != nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID:id,
		Balance: req.Balance,
	}

	account,err := server.store.UpdateAccount(ctx,arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,account)
}