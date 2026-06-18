package http

import "github.com/gin-gonic/gin"

func SetupGoldPriceRouter(r *gin.Engine, h *GoldPriceHandler) {
	api := r.Group("/gold")
	{
		api.GET("/price", h.GetGoldPrices)
		api.POST("/price", h.GenerateGoldPrice)
	}
}
