package channelStock

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type ChannelStockRepostory interface {
	Refresh(cs ChannelStock) error
}

func NewChannelStockRepository(db *sqlx.DB) ChannelStockRepostory {
	return &channelStockRepostory{
		db: db,
	}
}

type channelStockRepostory struct {
	db *sqlx.DB
}

func (r *channelStockRepostory) Refresh(cs ChannelStock) error {
	_, err := r.db.NamedExec(`REPLACE INTO channel_stock (sku, channel, quantity)
	VALUES (
		:sku,
		:channel,
		(
			SELECT SUM(quantity) FROM source_stock
			WHERE sku = :sku and source IN (
				SELECT source FROM source_channel
				WHERE channel = :channel
			) GROUP BY sku
		)
	)`, cs)

	if err == nil {
		log.Printf("SKU %s in channel %s succesfully refreshed", cs.Sku, cs.Channel)
	}

	return err
}
