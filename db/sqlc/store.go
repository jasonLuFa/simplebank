package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jasonLuFa/simplebank/util"
)

// store provides all functions to execute db queries and transactions
type Store struct{
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store{
	return &Store{
		db :db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context,fn func(*Queries) error) error{
	tx, err := store.db.BeginTx(ctx,nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil{
			return fmt.Errorf("tx err: %v, rb err: %v",err,rollbackErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct{
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID int64 `json:"to_account_id"`
	Amount string `json:"amount"`
}

type TransferTxResult struct{
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

var txKey = struct{}{}

func (store *Store) TransferTx(ctx context.Context,arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		// transfer
		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx,CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: 	 arg.ToAccountID,
			Amount: 			 arg.Amount,
		})
		if err != nil {
			return err
		}

		amount := util.StringToFloat64(arg.Amount)
		// fromEntry
		fmt.Println(txName, "create fromEntry")
		result.FromEntry, err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: fmt.Sprintf("%.2f",-amount),
		})
		if err != nil {
			return err
		}

		// toEntry
		fmt.Println(txName, "create toEntry")
		result.ToEntry, err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: fmt.Sprintf("%.2f",amount),
		})
		if err != nil {
			return err
		}

		// update accounts' balance
		// -- fromAccout
		fmt.Println(txName, "get fromAccout")
		fromAccount , err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}
		fmt.Println(txName, "update fromAccout")
		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID: arg.FromAccountID,
			Balance: fmt.Sprintf("%.2f",util.StringToFloat64(fromAccount.Balance) - util.StringToFloat64(arg.Amount)),
		})
		if err != nil {
			return err
		}

		// -- toAccout
		fmt.Println(txName, "get toAccout")
		toAccount , err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}
		fmt.Println(txName, "update toAccout")
		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID : arg.ToAccountID,
			Balance: fmt.Sprintf("%.2f",util.StringToFloat64(toAccount.Balance) + util.StringToFloat64(arg.Amount)),
		})
		if err != nil {
			return err
		}


		return nil
	})

	return result , err

}