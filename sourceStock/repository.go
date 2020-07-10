package sourceStock

import (
	"log"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

type SourceStockRepository interface {
	AddHistory(ss SourceStock) error
	Refresh(ss SourceStock) error
}

func NewSourceStockRepository(db *sqlx.DB) SourceStockRepository {
	return &sourceStockRepository{
		db: db,
	}
}

type sourceStockRepository struct {
	db *sqlx.DB
}

func (ssr *sourceStockRepository) AddHistory(ss SourceStock) error {
	_, err := ssr.db.NamedExec("INSERT INTO source_stock_history (sku, source, quantity, type, date) VALUES (:sku, :source, :quantity, :type, :date)", ss)

	return err
}

func (ssr *sourceStockRepository) Refresh(ss SourceStock) error {
	_, err := ssr.db.NamedExec(`REPLACE INTO source_stock (sku, source, quantity)
	VALUES (
		:sku,
		:source,
		IFNULL((
			SELECT IFNULL(MIN(quantity),0) + IFNULL(
				(
					SELECT SUM(quantity) FROM source_stock_history
					WHERE sku = ssh.sku and source = ssh.source and type <> 'Set' and date > ssh.date
				),0
			)
			FROM source_stock_history ssh 
			WHERE sku = :sku and source = :source and type = 'Set'
			AND date = (
				SELECT MAX(date) FROM source_stock_history
				WHERE sku = ssh.sku and source = ssh.source and type = 'Set'
			)
			GROUP BY source, sku, date
		),0)
	)`, ss)

	if err == nil {
		log.Printf("SKU %s in source %s succesfully refreshed", ss.Sku, ss.Source)
	}

	return err
}
