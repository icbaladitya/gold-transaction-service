package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type GoldTransactionInput struct {
	UserID string                `json:"user_id" binding:"required"`
	Type   string                `json:"type" binding:"required,oneof:BUY SELL"`
	Items  []GoldTransactionItem `json:"items" binding:"required,dive"`
}

type GoldTransactionItem struct {
	ID       string          `json:"id" binding:"required"`
	GoldGram decimal.Decimal `json:"gold_gram" binding:"required"`
	Version  int             `json:"version" binding:"required,gte=1"`
	Qty      int             `json:"qty" binding:"required,gte=1"`
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

type GoldTransactionRepoInterface interface {
	InsertTransactionHeader(input *TransactionHeaderInput) error
	InsertTransactionDetail(input []TransactionDetailInput) error
	// GoldTransactionSell()
	// GoldTransactionHistory()
}

type GoldTransactionUseCaseInterface interface {
	GoldTransactions(input *GoldTransactionInput) BasicResponse[any]
	GoldTransactionHistory(userId *string) BasicResponse[any]
}
