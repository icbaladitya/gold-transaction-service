package http

import "github.com/gin-gonic/gin"

func SetupGoldPriceRouter(r *gin.Engine, h *GoldPriceHandler) {
	api := r.Group("/gold")
	{
		api.GET("/price", h.GetGoldPrices)
		api.POST("/price", h.GenerateGoldPrice)
	}
}

func SetupTransactionRouter(r *gin.Engine, h *GoldTransactionHandler) {
	api := r.Group("/transaction")
	{
		api.POST("/buy", h.CreateTransaction)
		api.POST("/sell", h.CreateTransaction)
		api.GET("/history", h.GetTransactionHistory)
	}
}
