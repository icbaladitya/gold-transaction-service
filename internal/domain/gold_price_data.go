package domain

import "github.com/shopspring/decimal"

type GoldPriceData struct {
	ID           string          `json:"id"`
	GoldGram     float64         `json:"gold_gram"`
	Stock        int             `json:"stock"`
	BuyPrice     decimal.Decimal `json:"buy_price"`
	SellPrice    decimal.Decimal `json:"sell_price"`
	PricePerGram decimal.Decimal `json:"price_per_gram"`
	Version      int             `json:"version"`
}

type GenerateGoldPriceInput struct {
	PricePerGram     decimal.Decimal `json:"price_per_gram"`
	BuyPricePerGram  decimal.Decimal `json:"buy_price_per_gram"`
	SellPricePerGram decimal.Decimal `json:"sell_price_per_gram"`
}

type GenerateGoldPriceData struct {
	MstGoldID        string          `json:"mst_gold_id"`
	BuyPrice         decimal.Decimal `json:"buy_price"`
	SellPrice        decimal.Decimal `json:"sell_price"`
	PricePerGram     decimal.Decimal `json:"price_per_gram"`
	BuyPricePerGram  decimal.Decimal `json:"buy_price_per_gram"`
	SellPricePerGram decimal.Decimal `json:"sell_price_per_gram"`
	CreatedBy        string          `json:"created_by"`
}

type GoldPriceRepoInterface interface {
	GetGoldPrices() ([]GoldPriceData, error)
	GetGoldGrams() ([]BasicData, error)
	BulkInsertPrices(input []GenerateGoldPriceData) (int, error)
}

type GoldPriceUseCaseInterface interface {
	GetGoldPrices() BasicResponse[GoldPriceData]
	GenerateGoldPrice(input *GenerateGoldPriceInput) BasicResponse[any]
}
