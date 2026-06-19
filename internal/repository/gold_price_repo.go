package repository

import (
	"database/sql"
	"fmt"
	"gold-transaction-service/internal/domain"
	"strings"
)

type GoldPriceRepo struct {
	db *sql.DB
}

func GoldPriceRepoFunc(db *sql.DB) domain.GoldPriceRepoInterface {
	return &GoldPriceRepo{
		db: db,
	}
}

func (g *GoldPriceRepo) GetGoldPrices() ([]domain.GoldPriceData, error) {
	// panic("unimplemented")
	query := `SELECT gp.id, mg.gold_gram, mg.stock, gp.buy_price, gp.sell_price, gp.price_per_gram, gp.version
		FROM gold_prices gp
		JOIN mst_gold mg ON mg.id = gp.mst_gold_id
		WHERE gp.version = (SELECT MAX(version) FROM gold_prices)
		ORDER BY mg.gold_gram`

	rows, err := g.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []domain.GoldPriceData

	for rows.Next() {
		var p domain.GoldPriceData

		err := rows.Scan(
			&p.ID,
			&p.GoldGram,
			&p.Stock,
			&p.BuyPrice,
			&p.SellPrice,
			&p.PricePerGram,
			&p.Version,
		)
		if err != nil {
			return nil, err
		}

		prices = append(prices, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return prices, nil
}

func (g *GoldPriceRepo) GetGoldGrams() ([]domain.BasicData, error) {
	// panic("unimplemented")
	query := `SELECT id, gold_gram FROM mst_gold`
	rows, err := g.db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []domain.BasicData
	for rows.Next() {
		var b domain.BasicData
		if err := rows.Scan(&b.ID, &b.Name); err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
}

func (g *GoldPriceRepo) BulkInsertPrices(input []domain.GenerateGoldPriceData) (int, error) {
	// panic("unimplemented")
	baseQuery := `
    WITH next_version AS (
      SELECT COALESCE(MAX(version), 0) + 1 AS new_ver FROM gold_prices
    )
    INSERT INTO gold_prices (
      id, mst_gold_id, buy_price, sell_price, price_per_gram, buy_price_per_gram, sell_price_per_gram, version, created_by
    )
    VALUES 
  `

	valueStrings := make([]string, 0, len(input))
	valueArgs := make([]interface{}, 0, len(input)*7)

	paramCounter := 1

	for _, item := range input {
		rowQuery := fmt.Sprintf(
			"(replace(gen_random_uuid()::text, '-', ''), $%d, $%d, $%d, $%d, $%d, $%d, (SELECT new_ver FROM next_version), $%d)",
			paramCounter, paramCounter+1, paramCounter+2, paramCounter+3, paramCounter+4, paramCounter+5, paramCounter+6,
		)
		valueStrings = append(valueStrings, rowQuery)

		valueArgs = append(valueArgs, item.MstGoldID)
		valueArgs = append(valueArgs, item.BuyPrice)
		valueArgs = append(valueArgs, item.SellPrice)
		valueArgs = append(valueArgs, item.PricePerGram)
		valueArgs = append(valueArgs, item.BuyPricePerGram)
		valueArgs = append(valueArgs, item.SellPricePerGram)
		valueArgs = append(valueArgs, item.CreatedBy)

		paramCounter += 7
	}

	completeQuery := baseQuery + strings.Join(valueStrings, ",\n")

	_, err := g.db.Exec(completeQuery, valueArgs...)
	if err != nil {
		return 0, fmt.Errorf("error executing bulk insert gold prices : %w", err)
	}

	return 1, nil
}
