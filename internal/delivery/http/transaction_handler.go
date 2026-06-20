package http

import (
	"context"
	"gold-transaction-service/internal/domain"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type GoldTransactionHandler struct {
	usecase domain.GoldTransactionUseCaseInterface
}

func GoldTransactionHandlerFunc(u domain.GoldTransactionUseCaseInterface) *GoldTransactionHandler {
	return &GoldTransactionHandler{
		usecase: u,
	}
}

func (h *GoldTransactionHandler) CreateTransaction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	var input domain.GoldTransactionInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse[any]("Validasi input gagal: "+err.Error()))
		return
	}

	if strings.Contains(c.Request.URL.Path, "/buy") {
		input.Type = "BUY"
	} else if strings.Contains(c.Request.URL.Path, "/sell") {
		input.Type = "SELL"
	}

	result := h.usecase.GoldTransactions(ctx, &input)

	if result.ResultCode == -1 || result.ResultCode == 0 {
		c.JSON(http.StatusUnprocessableEntity, result)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *GoldTransactionHandler) GetTransactionHistory(c *gin.Context) {
	userID := c.Query("user-id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, domain.FailResponse[any]("Parameter user-id wajib diisi"))
		return
	}

	ctx := c.Request.Context()
	response := h.usecase.GoldTransactionHistory(ctx, &userID)

	if response.ResultCode == -1 {
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if response.ResultCode == 0 {
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
