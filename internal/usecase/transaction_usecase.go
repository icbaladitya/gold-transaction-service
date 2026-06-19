package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"gold-transaction-service/internal/domain"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type GoldTransactionUsecase struct {
	db       *sql.DB
	repoFunc domain.GoldTransactionRepoInterface
}

func GoldTransactionUsecaseFunc(db *sql.DB, r domain.GoldTransactionRepoInterface) domain.GoldTransactionUseCaseInterface {
	return &GoldTransactionUsecase{
		db:       db,
		repoFunc: r,
	}
}

func (u *GoldTransactionUsecase) GoldTransactions(ctx context.Context, input *domain.GoldTransactionInput) domain.BasicResponse[any] {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.ErrorResponse[any]("Gagal memulai transaksi database: " + err.Error())
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	userBalance, err := u.repoFunc.ValidationBalance(ctx, tx, &input.UserID)
	if err != nil {
		tx.Rollback()
		return domain.FailResponse[any]("Data saldo tidak di temukan : " + err.Error())
	}

	var totalGoldGram decimal.Decimal
	var totalGoldIDR decimal.Decimal
	var totalQty int = 0

	headerTrxID := strings.ReplaceAll(uuid.New().String(), "-", "")
	var detailInputs []domain.TransactionDetailInput

	for _, item := range input.Items {

		goldPrice, err := u.repoFunc.GetGoldPrice(ctx, tx, &item.ID)
		if err != nil {
			tx.Rollback()
			return domain.FailResponse[any]("Harga emas untuk ID " + goldPrice.ID + " tidak tersedia")
		}

		valStock := u.repoFunc.ValidationStock(ctx, tx, &goldPrice.GoldID)
		if valStock != nil {
			tx.Rollback()
			return domain.FailResponse[any]("Stok emas untuk ID " + item.ID + " tidak tersedia")
		}

		qtyDecimal := decimal.NewFromInt(int64(item.Qty))
		pricePerGram := goldPrice.PricePerGram

		itemTotalIDR := pricePerGram.Mul(item.GoldGram).Mul(qtyDecimal)
		itemTotalGram := item.GoldGram.Mul(qtyDecimal)

		totalGoldGram = totalGoldGram.Add(itemTotalGram)
		totalGoldIDR = totalGoldIDR.Add(itemTotalIDR)
		totalQty += item.Qty

		err = u.repoFunc.UpdateStockGold(ctx, tx, &goldPrice.GoldID, &item.Qty, &input.Type)
		if err != nil {
			tx.Rollback()
			return domain.ErrorResponse[any]("Gagal memperbarui stok emas: " + err.Error())
		}

		detailInputs = append(detailInputs, domain.TransactionDetailInput{
			ID:           strings.ReplaceAll(uuid.New().String(), "-", ""),
			GoldTrxHdrID: headerTrxID,
			GoldPricesID: goldPrice.ID,
			GoldGram:     item.GoldGram,
			BuyPrice:     goldPrice.BuyPrice,
			SellPrice:    goldPrice.SellPrice,
			Created:      time.Now(),
			CreatedBy:    "system",
			Qty:          item.Qty,
		})
	}

	var newIDRBalance decimal.Decimal
	var newGoldBalance decimal.Decimal

	switch input.Type {
	case "BUY":
		if userBalance.IDRBalance.LessThan(totalGoldIDR) {
			tx.Rollback()
			return domain.FailResponse[any]("Saldo uang Anda tidak mencukupi")
		}
		newIDRBalance = userBalance.IDRBalance.Sub(totalGoldIDR)
		newGoldBalance = userBalance.GoldBalance.Add(totalGoldGram)
	case "SELL":
		if userBalance.GoldBalance.LessThan(totalGoldGram) {
			tx.Rollback()
			return domain.FailResponse[any]("Saldo emas Anda tidak mencukupi untuk dijual")
		}
		newIDRBalance = userBalance.IDRBalance.Add(totalGoldIDR)
		newGoldBalance = userBalance.GoldBalance.Sub(totalGoldGram)
	default:
		tx.Rollback()
		return domain.FailResponse[any]("Tipe transaksi tidak valid")
	}

	balanceInput := &domain.UserBalanceInput{
		ID:          strings.ReplaceAll(uuid.New().String(), "-", ""),
		UserID:      input.UserID,
		IDRBalance:  newIDRBalance,
		GoldBalance: newGoldBalance,
		CreatedBy:   input.UserID,
	}

	err = u.repoFunc.InsertBalanceUser(ctx, tx, balanceInput)
	if err != nil {
		tx.Rollback()
		return domain.ErrorResponse[any]("Gagal memperbarui saldo pengguna: " + err.Error())
	}

	desc := fmt.Sprintf("Transaksi %s emas sukses", input.Type)
	headerInput := &domain.TransactionHeaderInput{
		ID:            headerTrxID,
		UserID:        input.UserID,
		Type:          input.Type,
		TotalGoldGram: totalGoldGram,
		TotalGoldIDR:  totalGoldIDR,
		Status:        "SUCCESS",
		Description:   &desc,
		Created:       time.Now(),
		CreatedBy:     "system",
		TotalQty:      totalQty,
	}

	err = u.repoFunc.InsertTransactionHeader(ctx, tx, headerInput)
	if err != nil {
		tx.Rollback()
		return domain.ErrorResponse[any]("Gagal membuat nota transaksi: " + err.Error())
	}

	err = u.repoFunc.InsertTransactionDetail(ctx, tx, detailInputs)
	if err != nil {
		tx.Rollback()
		return domain.ErrorResponse[any]("Gagal menyimpan rincian transaksi: " + err.Error())
	}

	if err := tx.Commit(); err != nil {
		return domain.ErrorResponse[any]("Gagal melakukan commit database: " + err.Error())
	}

	summaryData := map[string]interface{}{
		"trx_id":          headerTrxID,
		"status":          "SUCCESS",
		"current_idr":     newIDRBalance.StringFixed(2),
		"current_gold_gr": newGoldBalance.StringFixed(4),
	}

	return domain.SuccessDataResponse[any](summaryData, "Transaksi dan pembaruan saldo berhasil!")
}

func (u *GoldTransactionUsecase) GoldTransactionHistory(ctx context.Context, userId *string) domain.BasicResponse[any] {
	return domain.SuccessDataResponse[any](nil, "Belum diimplementasikan")
}
