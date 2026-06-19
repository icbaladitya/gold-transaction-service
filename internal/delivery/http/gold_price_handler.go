package http

import (
	"gold-transaction-service/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GoldPriceHandler struct {
	useCase domain.GoldPriceUseCaseInterface
}

func GoldPriceHandlerFunc(useCase domain.GoldPriceUseCaseInterface) *GoldPriceHandler {
	return &GoldPriceHandler{useCase: useCase}
}

func (h *GoldPriceHandler) GetGoldPrices(ctx *gin.Context) {
	result := h.useCase.GetGoldPrices()

	switch result.ResultCode {
	case 1:
		ctx.JSON(http.StatusOK, result)
	case 0:
		ctx.JSON(http.StatusBadRequest, result)
	default:
		ctx.JSON(http.StatusInternalServerError, result)
	}
}

func (h *GoldPriceHandler) GenerateGoldPrice(ctx *gin.Context) {
	var req domain.GenerateGoldPriceInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result := h.useCase.GenerateGoldPrice(&req)

	switch result.ResultCode {
	case 1:
		ctx.JSON(http.StatusOK, result)
	case 0:
		ctx.JSON(http.StatusBadRequest, result)
	default:
		ctx.JSON(http.StatusInternalServerError, result)
	}
}
