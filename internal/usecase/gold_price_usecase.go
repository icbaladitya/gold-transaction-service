package usecase

import (
	"gold-transaction-service/internal/domain"

	"github.com/shopspring/decimal"
)

type GoldPriceUsecase struct {
	repo domain.GoldPriceRepoInterface
}

func GoldPriceUsecaseFunc(r domain.GoldPriceRepoInterface) domain.GoldPriceUseCaseInterface {
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

		totalBuyPrice := gramVal.Mul(input.BuyPricePerGram)
		totalSellPrice := gramVal.Mul(input.SellPricePerGram)

		generateGoldPrices = append(generateGoldPrices, domain.GenerateGoldPriceData{
			MstGoldID:        item.ID,
			BuyPrice:         totalBuyPrice,
			SellPrice:        totalSellPrice,
			PricePerGram:     input.PricePerGram,
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
