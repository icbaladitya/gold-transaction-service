package usecase

import (
	"gold-transaction-service/internal/domain"

	"github.com/shopspring/decimal"
)

type GoldPriceUsecase struct {
	repo domain.GoldPriceRepoFunc
}

func GoldPriceUsecaseFunc(r domain.GoldPriceRepoFunc) domain.GoldPriceUsecaseFunc {
	return &GoldPriceUsecase{repo: r}
}

func (g *GoldPriceUsecase) GetGoldPrices() domain.BasicResponse[domain.GoldPriceData] {
	items, err := g.repo.GetGoldPrices()

	if err != nil {
		return domain.ErrorResponse[domain.GoldPriceData]("Error : " + err.Error())
	}

	if items == nil {
		// items = []domain.GoldPriceData{}
		return domain.FailResponse[domain.GoldPriceData]("Data is empty")
	}
	return domain.SuccessListResponse(items, "successfully")
}

func (g *GoldPriceUsecase) GenerateGoldPrice(input *domain.GenerateGoldPriceInput) domain.BasicResponse[any] {
	// panic("unimplemented")
	items, err := g.repo.GetGoldGrams()
	if err != nil {
		return domain.ErrorResponse[any]("Error : " + err.Error())
	}

	if items == nil {
		items = []domain.BasicData{}
	}

	var generateGoldPrices []domain.GenerateGoldPriceData

	for _, item := range items {
		gramVal, _ := decimal.NewFromString(item.Name)

		// totalPrice := gramVal.Mul(input.PricePerGram)
		totalBuyPrice := gramVal.Mul(input.BuyPricePerGram)
		totalSellPrice := gramVal.Mul(input.SellPricePerGram)

		// fmt.Printf("Generated -> %s gr | Price: Rp %s | Buy: Rp %s | Sell: Rp %s\n",
		// 	item.Name, totalPrice.StringFixed(2), totalBuyPrice.StringFixed(2), totalSellPrice.StringFixed(2))

		generateGoldPrices = append(generateGoldPrices, domain.GenerateGoldPriceData{
			MstGoldID:        item.ID,
			BuyPrice:         totalBuyPrice,
			SellPrice:        totalSellPrice,
			BuyPricePerGram:  input.BuyPricePerGram,
			SellPricePerGram: input.SellPricePerGram,
			CreatedBy:        "system",
		})
	}

	_, err = g.repo.BulkInsertPrices(generateGoldPrices)

	if err != nil {
		return domain.ErrorResponse[any]("Error insert DB: " + err.Error())
	}

	return domain.SuccessResponse[any]("Generate Success")
}
