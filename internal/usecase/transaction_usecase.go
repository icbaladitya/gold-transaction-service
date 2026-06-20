package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"gold-transaction-service/internal/domain"
	"strconv"
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

	var totalGoldGramHdr decimal.Decimal
	var totalGoldIDRHdr decimal.Decimal
	var totalQtyHdr int = 0

	headerTrxID := strings.ReplaceAll(uuid.New().String(), "-", "")
	var detailInputs []domain.TransactionDetailInput

	for _, item := range input.Items {

		goldPrice, err := u.repoFunc.GetGoldPrice(ctx, tx, &item.ID)
		if err != nil {
			tx.Rollback()
			return domain.FailResponse[any]("Harga emas untuk ID " + goldPrice.ID + " tidak tersedia")
		}

		valStock, err := u.repoFunc.ValidationStock(ctx, tx, &goldPrice.GoldID)
		if err != nil {
			tx.Rollback()
			return domain.FailResponse[any]("Gagal cek stok " + item.ID + " : " + err.Error())
		}

		if input.Type == "BUY" {
			if *valStock < item.Qty {
				tx.Rollback()
				if *valStock == 0 {
					return domain.FailResponse[any]("Stock " + item.GoldGram.String() + " tidak tersedia")
				} else {
					return domain.FailResponse[any]("Stock " + item.GoldGram.String() + " hanya tersedia " + strconv.Itoa(*valStock))
				}
			}
		}

		err = u.repoFunc.UpdateStockGold(ctx, tx, &goldPrice.GoldID, &item.Qty, &input.Type)
		if err != nil {
			tx.Rollback()
			return domain.ErrorResponse[any]("Gagal memperbarui stok emas: " + err.Error())
		}

		qtyDecimal := decimal.NewFromInt(int64(item.Qty))

		totalGram := goldPrice.GoldGram.Mul(qtyDecimal)
		totalPrice := goldPrice.BuyPrice.Mul(qtyDecimal)

		totalGoldGramHdr = totalGoldGramHdr.Add(totalGram)
		totalGoldIDRHdr = totalGoldIDRHdr.Add(totalPrice)
		totalQtyHdr += item.Qty

		detailInputs = append(detailInputs, domain.TransactionDetailInput{
			ID:           strings.ReplaceAll(uuid.New().String(), "-", ""),
			GoldTrxHdrID: headerTrxID,
			GoldPricesID: goldPrice.ID,
			GoldGram:     item.GoldGram,
			BuyPrice:     goldPrice.BuyPrice,
			SellPrice:    goldPrice.SellPrice,
			TotalPrice:   totalPrice,
			TotalGram:    totalGram,
			Created:      time.Now(),
			CreatedBy:    "system",
			Qty:          item.Qty,
		})
	}

	var newIDRBalance decimal.Decimal
	var newGoldBalance decimal.Decimal

	switch input.Type {
	case "BUY":
		if userBalance.IDRBalance.LessThan(totalGoldIDRHdr) {
			tx.Rollback()
			return domain.FailResponse[any]("Saldo uang Anda tidak mencukupi")
		}
		newIDRBalance = userBalance.IDRBalance.Sub(totalGoldIDRHdr)
		newGoldBalance = userBalance.GoldBalance.Add(totalGoldGramHdr)
	case "SELL":
		if userBalance.GoldBalance.LessThan(totalGoldGramHdr) {
			tx.Rollback()
			return domain.FailResponse[any]("Saldo emas Anda tidak mencukupi untuk dijual")
		}
		newIDRBalance = userBalance.IDRBalance.Add(totalGoldIDRHdr)
		newGoldBalance = userBalance.GoldBalance.Sub(totalGoldGramHdr)
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
		TotalGoldGram: totalGoldGramHdr,
		TotalGoldIDR:  totalGoldIDRHdr,
		Status:        "SUCCESS",
		Description:   &desc,
		Created:       time.Now(),
		CreatedBy:     "system",
		TotalQty:      totalQtyHdr,
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
		"gold_trx_hdr_id": headerTrxID,
		"status":          "SUCCESS",
		"idr_balance":     newIDRBalance.StringFixed(2),
		"gold_balance":    newGoldBalance.StringFixed(2),
	}

	return domain.SuccessDataResponse[any](summaryData, "Transaksi dan pembaruan saldo berhasil!")
}

func (u *GoldTransactionUsecase) GoldTransactionHistory(ctx context.Context, userId *string) domain.BasicResponse[domain.TransactionHistoryHeader] {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.ErrorResponse[domain.TransactionHistoryHeader]("Gagal memulai transaksi database: " + err.Error())
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	dataHeader, err := u.repoFunc.GetTransactionHeader(ctx, tx, userId)
	if err != nil {
		tx.Rollback()
		return domain.ErrorResponse[domain.TransactionHistoryHeader]("Transaksi header : " + err.Error())
	}

	if len(dataHeader) == 0 {
		tx.Rollback()
		return domain.FailResponse[domain.TransactionHistoryHeader]("Data transaksi tidak ada")
	}
	fmt.Println("masuk sinii ", dataHeader)

	var dataResponse []domain.TransactionHistoryHeader

	for _, item := range dataHeader {
		fmt.Println("masuk sinii looping ", item.GoldTrxHdrID)
		dataDetail, err := u.repoFunc.GetTransactionDetail(ctx, tx, &item.GoldTrxHdrID)
		if err != nil {
			tx.Rollback()
			return domain.ErrorResponse[domain.TransactionHistoryHeader]("Transaksi detail : " + err.Error())
		}

		fmt.Println("masuk sinii dataDetail ", dataDetail)

		history := domain.TransactionHistoryHeader{
			GoldTrxHdrID:  item.GoldTrxHdrID,
			Type:          item.Type,
			TotalGoldGram: item.TotalGoldGram,
			TotalGoldIDR:  item.TotalGoldIDR,
			TotalQty:      item.TotalQty,
			Status:        item.Status,
			Items:         dataDetail,
		}

		dataResponse = append(dataResponse, history)
	}

	fmt.Println("masuk sinii dataResponse ", dataResponse)

	if err := tx.Commit(); err != nil {
		return domain.ErrorResponse[domain.TransactionHistoryHeader]("Gagal melakukan commit database: " + err.Error())
	}

	return domain.SuccessListResponse[domain.TransactionHistoryHeader](dataResponse, "Successfully")
}
