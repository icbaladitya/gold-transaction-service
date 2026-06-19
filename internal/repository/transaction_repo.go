package repository

import (
	"database/sql"
	"fmt"
	"gold-transaction-service/internal/domain"
	"strings"
)

type TransactionRepo struct {
	db *sql.DB
}

func GoldTransactionRepoFunc(db *sql.DB) domain.GoldTransactionRepoInterface {
	return &TransactionRepo{
		db: db,
	}
}

// InsertTransactionDetail implements [domain.GoldTransactionRepoInterface].
func (t *TransactionRepo) InsertTransactionDetail(input []domain.TransactionDetailInput) error {
	// panic("unimplemented")
	if len(input) == 0 {
		return nil
	}

	baseQuery := `
		INSERT INTO gold_transaction_details (
			id, gold_trx_hdr_id, gold_prices_id, gold_gram, buy_price, sell_price, created_by, qty
		) VALUES 
	`

	valueStrings := make([]string, 0, len(input))
	valueArgs := make([]interface{}, 0, len(input)*8)

	paramCounter := 1

	for _, item := range input {
		rowQuery := fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			paramCounter, paramCounter+1, paramCounter+2, paramCounter+3, paramCounter+4, paramCounter+5, paramCounter+6, paramCounter+7,
		)
		valueStrings = append(valueStrings, rowQuery)

		valueArgs = append(valueArgs, item.ID)
		valueArgs = append(valueArgs, item.GoldTrxHdrID)
		valueArgs = append(valueArgs, item.GoldPricesID)
		valueArgs = append(valueArgs, item.GoldGram)
		valueArgs = append(valueArgs, item.BuyPrice)
		valueArgs = append(valueArgs, item.SellPrice)
		valueArgs = append(valueArgs, item.CreatedBy)
		valueArgs = append(valueArgs, item.Qty)

		paramCounter += 8
	}

	completeQuery := baseQuery + strings.Join(valueStrings, ",\n")

	_, err := t.db.Exec(completeQuery, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

// InsertTransactionHeader implements [domain.GoldTransactionRepoInterface].
func (t *TransactionRepo) InsertTransactionHeader(input *domain.TransactionHeaderInput) error {
	// panic("unimplemented")
	query := `
		INSERT INTO gold_trx_hdr (
			id, user_id, type, total_gold_gram, total_gold_idr, status, description, created_by, total_qty
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := t.db.Exec(query,
		input.ID,
		input.UserID,
		input.Type,
		input.TotalGoldGram,
		input.TotalGoldIDR,
		input.Status,
		input.Description,
		input.CreatedBy,
		input.TotalQty,
	)

	if err != nil {
		return err
	}

	return nil
}
