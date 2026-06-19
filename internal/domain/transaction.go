package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type GoldTransactionInput struct {
	UserID string                `json:"user_id" binding:"required"`
	Type   string                `json:"type,omitempty"`
	Items  []GoldTransactionItem `json:"items" binding:"required,dive"`
}

type GoldTransactionItem struct {
	ID       string          `json:"id" binding:"required"`
	GoldGram decimal.Decimal `json:"gold_gram" binding:"required"`
	Version  int             `json:"version" binding:"required,gte=1"`
	Qty      int             `json:"qty" binding:"required,gte=1"`
	// PricePerGram decimal.Decimal `json:"price_per_gram" binding:"required"`
}

type TransactionHeaderInput struct {
	ID            string          `json:"id" db:"id"`
	UserID        string          `json:"user_id" db:"user_id"`
	Type          string          `json:"type" db:"type"`
	TotalGoldGram decimal.Decimal `json:"total_gold_gram" db:"total_gold_gram"`
	TotalGoldIDR  decimal.Decimal `json:"total_gold_idr" db:"total_gold_idr"`
	Status        string          `json:"status" db:"status"`
	Description   *string         `json:"description" db:"description"`
	Created       time.Time       `json:"created" db:"created"`
	CreatedBy     string          `json:"created_by" db:"created_by"`
	TotalQty      int             `json:"total_qty" db:"total_qty"`
}

type TransactionDetailInput struct {
	ID           string          `json:"id" db:"id"`
	GoldTrxHdrID string          `json:"gold_trx_hdr_id" db:"gold_trx_hdr_id"`
	GoldPricesID string          `json:"gold_prices_id" db:"gold_prices_id"`
	GoldGram     decimal.Decimal `json:"gold_gram" db:"gold_gram"`
	BuyPrice     decimal.Decimal `json:"buy_price" db:"buy_price"`
	SellPrice    decimal.Decimal `json:"sell_price" db:"sell_price"`
	Created      time.Time       `json:"created" db:"created"`
	CreatedBy    string          `json:"created_by" db:"created_by"`
	Qty          int             `json:"qty" db:"qty"`
}

type UserBalanceData struct {
	ID          string          `json:"id" db:"id"`
	IDRBalance  decimal.Decimal `json:"idr_balance" db:"idr_balance"`
	GoldBalance decimal.Decimal `json:"gold_balance" db:"gold_balance"`
	Version     int             `json:"version" db:"version"`
}

type UserBalanceInput struct {
	ID          string          `json:"id" db:"id"`
	UserID      string          `json:"user_id" db:"user_id"`
	IDRBalance  decimal.Decimal `json:"idr_balance" db:"idr_balance"`
	GoldBalance decimal.Decimal `json:"gold_balance" db:"gold_balance"`
	Version     int             `json:"version" db:"version"`
	CreatedBy   string          `json:"created_by" db:"created_by"`
}

type GoldPrice struct {
	ID           string          `json:"id" db:"id"`
	GoldID       string          `json:"mst_gold_id" db:"mst_gold_id"`
	BuyPrice     decimal.Decimal `json:"buy_price" db:"buy_price"`
	SellPrice    decimal.Decimal `json:"sell_price" db:"sell_price"`
	PricePerGram decimal.Decimal `json:"price_per_gram" db:"price_per_gram"`
	Version      int             `json:"version" db:"version"`
}

type GoldTransactionRepoInterface interface {
	InsertTransactionHeader(ctx context.Context, tx *sql.Tx, input *TransactionHeaderInput) error
	InsertTransactionDetail(ctx context.Context, tx *sql.Tx, input []TransactionDetailInput) error
	UpdateStockGold(ctx context.Context, tx *sql.Tx, goldId *string, stock *int, tipe *string) error
	ValidationBalance(ctx context.Context, tx *sql.Tx, userId *string) (*UserBalanceData, error)
	ValidationStock(ctx context.Context, tx *sql.Tx, goldId *string) error
	InsertBalanceUser(ctx context.Context, tx *sql.Tx, input *UserBalanceInput) error
	GetGoldPrice(ctx context.Context, tx *sql.Tx, goldPriceId *string) (*GoldPrice, error)
}

type GoldTransactionUseCaseInterface interface {
	GoldTransactions(ctx context.Context, input *GoldTransactionInput) BasicResponse[any]
	GoldTransactionHistory(ctx context.Context, userId *string) BasicResponse[any]
}
