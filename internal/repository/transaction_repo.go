package repository

import (
	"context"
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

func (t *TransactionRepo) InsertTransactionHeader(ctx context.Context, tx *sql.Tx, input *domain.TransactionHeaderInput) error {
	query := `
    INSERT INTO gold_trx_hdr (
      id, user_id, type, total_gold_gram, total_gold_idr, status, description, created_by, total_qty
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
  `
	_, err := tx.ExecContext(ctx, query,
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
	return err
}

func (t *TransactionRepo) InsertTransactionDetail(ctx context.Context, tx *sql.Tx, input []domain.TransactionDetailInput) error {
	if len(input) == 0 {
		return nil
	}

	baseQuery := `
    INSERT INTO gold_trx_dtl (
      id, gold_trx_hdr_id, gold_prices_id, gold_gram, buy_price, sell_price, created_by, qty, total_price, total_gram
    ) VALUES 
  `

	valueStrings := make([]string, 0, len(input))
	valueArgs := make([]interface{}, 0, len(input)*10)
	paramCounter := 1

	for _, item := range input {
		rowQuery := fmt.Sprintf(
			"($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			paramCounter, paramCounter+1, paramCounter+2, paramCounter+3, paramCounter+4, paramCounter+5, paramCounter+6,
			paramCounter+7, paramCounter+8, paramCounter+9,
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
		valueArgs = append(valueArgs, item.TotalPrice)
		valueArgs = append(valueArgs, item.TotalGram)

		paramCounter += 10
	}

	completeQuery := baseQuery + strings.Join(valueStrings, ",\n")

	_, err := tx.ExecContext(ctx, completeQuery, valueArgs...)
	return err
}

func (t *TransactionRepo) UpdateStockGold(ctx context.Context, tx *sql.Tx, goldId *string, stock *int, tipe *string) error {
	var query string

	if *tipe == "BUY" {
		query = `UPDATE mst_gold SET stock = stock - $1 WHERE id = $2`
	} else {
		query = `UPDATE mst_gold SET stock = stock + $1 WHERE id = $2`
	}

	_, err := tx.ExecContext(ctx, query, stock, goldId)
	return err
}

func (t *TransactionRepo) InsertBalanceUser(ctx context.Context, tx *sql.Tx, input *domain.UserBalanceInput) error {
	query := `
    WITH balance AS (
      SELECT COALESCE(MAX(version), 0) + 1 AS new_version
      FROM user_balance
      WHERE user_id = $2 
    )
    INSERT INTO user_balance (id, user_id, idr_balance, gold_balance, version, created_by)
    SELECT $1, $2, $3, $4, b.new_version, $5
    FROM balance b
  `

	_, err := tx.ExecContext(ctx, query,
		input.ID,
		input.UserID,
		input.IDRBalance,
		input.GoldBalance,
		input.CreatedBy,
	)
	return err
}

func (t *TransactionRepo) ValidationStock(ctx context.Context, tx *sql.Tx, goldId *string) (*int, error) {
	var stock int
	query := `SELECT stock FROM mst_gold WHERE id = $1`
	err := tx.QueryRowContext(ctx, query, goldId).Scan(&stock)

	return &stock, err
}

func (t *TransactionRepo) ValidationBalance(ctx context.Context, tx *sql.Tx, userId *string) (*domain.UserBalanceData, error) {
	query := `
    SELECT id, idr_balance, gold_balance, version 
    FROM user_balance
    WHERE version = (SELECT MAX(version) FROM user_balance)
    AND user_id = $1
    `
	var userBalanceData domain.UserBalanceData

	err := tx.QueryRowContext(ctx, query, userId).Scan(
		&userBalanceData.ID,
		&userBalanceData.IDRBalance,
		&userBalanceData.GoldBalance,
		&userBalanceData.Version,
	)

	if err != nil {
		return nil, err
	}
	return &userBalanceData, nil
}

func (t *TransactionRepo) GetGoldPrice(ctx context.Context, tx *sql.Tx, goldPriceId *string) (*domain.GoldPrice, error) {
	query := `
		SELECT gp.id, gp.mst_gold_id, gp.buy_price, gp.sell_price, gp.price_per_gram, gp.version, mg.gold_gram 
		FROM gold_prices gp
		JOIN mst_gold mg ON mg.id = gp.mst_gold_id
		WHERE gp.id = $1`

	var goldPrice domain.GoldPrice
	err := tx.QueryRowContext(ctx, query, goldPriceId).Scan(
		&goldPrice.ID, &goldPrice.GoldID, &goldPrice.BuyPrice, &goldPrice.SellPrice, &goldPrice.PricePerGram, &goldPrice.Version, &goldPrice.GoldGram,
	)

	if err != nil {
		return nil, err
	}
	return &goldPrice, nil
}

func (t *TransactionRepo) GetTransactionHeader(ctx context.Context, tx *sql.Tx, userId *string) ([]domain.TransactionHistoryHeader, error) {
	query := `
		SELECT id, type, total_gold_gram, total_gold_idr, total_qty, status
		FROM gold_trx_hdr
		WHERE user_id = $1
		ORDER BY created DESC`

	rows, err := tx.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listHeaderTransaction []domain.TransactionHistoryHeader

	for rows.Next() {
		var header domain.TransactionHistoryHeader

		err := rows.Scan(
			&header.GoldTrxHdrID,
			&header.Type,
			&header.TotalGoldGram,
			&header.TotalGoldIDR,
			&header.TotalQty,
			&header.Status,
		)
		if err != nil {
			return nil, err
		}

		listHeaderTransaction = append(listHeaderTransaction, header)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return listHeaderTransaction, nil
}

func (t *TransactionRepo) GetTransactionDetail(ctx context.Context, tx *sql.Tx, goldHeaderId *string) ([]domain.TransactionHistoryDetail, error) {
	query := `
		SELECT id, gold_gram, qty, buy_price, sell_price, total_price, total_gram
		FROM gold_trx_dtl
		WHERE gold_trx_hdr_id = $1
		ORDER BY created DESC`

	rows, err := tx.QueryContext(ctx, query, goldHeaderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listDetailTransaction []domain.TransactionHistoryDetail

	for rows.Next() {
		var detail domain.TransactionHistoryDetail

		err := rows.Scan(
			&detail.GoldTrxDetailID,
			&detail.GoldGram,
			&detail.Qty,
			&detail.BuyPrice,
			&detail.SellPrice,
			&detail.TotalPrice,
			&detail.TotalGram,
		)
		if err != nil {
			return nil, err
		}

		listDetailTransaction = append(listDetailTransaction, detail)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return listDetailTransaction, nil
}
