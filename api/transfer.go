package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/bacnx/simplebank/db/sqlc"
	"github.com/bacnx/simplebank/token"
	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.FromAccountID == req.ToAccountID {
		err := errors.New("from account must be defferent to account")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	isValid, fromOwner := server.validCurrency(ctx, req.FromAccountID, req.Currency)
	if !isValid {
		return
	}

	if fromOwner != payload.Username {
		err := errors.New("from account doesn't belong to authorized account")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	isValid, _ = server.validCurrency(ctx, req.ToAccountID, req.Currency)
	if !isValid {
		return
	}

	transferTxParams := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, transferTxParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validCurrency(ctx *gin.Context, accountID int64, currency string) (bool, string) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false, ""
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false, ""
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false, ""
	}

	return true, account.Owner
}
