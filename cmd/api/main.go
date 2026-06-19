package main

import (
	"gold-transaction-service/config"
	delivery "gold-transaction-service/internal/delivery/http"
	"gold-transaction-service/internal/repository"
	"gold-transaction-service/internal/usecase"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db := config.InitDB()
	defer db.Close()

	goldPriceRepo := repository.GoldPriceRepoFunc(db)
	goldPrice := usecase.GoldPriceUsecaseFunc(goldPriceRepo)
	goldPriceHandler := delivery.GoldPriceHandlerFunc(goldPrice)

	goldTransactionRepo := repository.GoldTransactionRepoFunc(db)
	goldTransactionUsecase := usecase.GoldTransactionUsecaseFunc(db, goldTransactionRepo)
	goldTransactionHandler := delivery.GoldTransactionHandlerFunc(goldTransactionUsecase)

	r := gin.Default()
	delivery.SetupGoldPriceRouter(r, goldPriceHandler)
	delivery.SetupTransactionRouter(r, goldTransactionHandler)

	r.Run(os.Getenv("APP_PORT"))
}
