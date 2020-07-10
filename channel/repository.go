package channel

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

type ChannelRepository interface {
	GetChannelsRelatedToSource(s string) ([]string, error)
}

func NewChannelRepository(db *sqlx.DB) ChannelRepository {
	return &channelRepository{
		db: db,
	}
}

type channelRepository struct {
	db *sqlx.DB
}

func (cr *channelRepository) GetChannelsRelatedToSource(s string) ([]string, error) {
	chs := []string{}
	err := cr.db.Select(&chs, "SELECT channel FROM source_channel where source='"+s+"'")

	if err != nil {
		return nil, err
	}

	return chs, nil
}
