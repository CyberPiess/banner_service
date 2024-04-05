package banner

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type BannerRepository struct {
	db *sql.DB
}

func NewBannerRepository(db *sql.DB) *BannerRepository {
	return &BannerRepository{db: db}
}

func (ur *BannerRepository) Get(ctx context.Context, bannerParams BannerRequest) (BannerResponse, error) {
	var banner BannerResponse

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.Select("content").
		From("banners").
		Where("feature_id = ? and ? = ANY(tag_id)", bannerParams.FeatureId, bannerParams.TagId).ToSql()
	if err != nil {
		return BannerResponse{}, err
	}

	row := ur.db.QueryRow(query, args...)
	err = row.Scan(&banner.Content)
	if err != nil {
		return BannerResponse{}, err
	}

	return banner, nil
}
