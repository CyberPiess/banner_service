package banner

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type BannerRepository struct {
	db *sql.DB
}

func NewBannerRepository(db *sql.DB) *BannerRepository {
	return &BannerRepository{db: db}
}

func (bn *BannerRepository) Get(bannerParams BannerRequest) ([]BannerResponse, error) {
	var banner BannerResponse

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.Select("content").
		From("banners as b").
		Join("tags as t on b.id = t.banner_id").Join("features as f on b.id = f.banner_id").
		Where("feature_id = ? and tag_id = ? and is_active = true", bannerParams.FeatureId, bannerParams.TagId).ToSql()
	if err != nil {
		return []BannerResponse{}, err
	}

	row := bn.db.QueryRow(query, args...)
	err = row.Scan(&banner.Content)
	if err != nil {
		return []BannerResponse{}, err
	}

	return []BannerResponse{banner}, nil
}

func (bn *BannerRepository) GetAllBanners(bannerParams BannerRequest) ([]BannerResponse, error) {

	var bannerResultSlice []BannerResponse
	var bannersQuery string
	var args []interface{}
	var err error

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	if bannerParams.FeatureId != 0 || bannerParams.TagId != 0 {
		var searchBannersID sq.SelectBuilder
		psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

		if bannerParams.FeatureId != 0 && bannerParams.TagId == 0 {
			searchBannersID = psql.Select("banner_id").
				From("features").Where("feature_id = ?", bannerParams.FeatureId)
		} else if bannerParams.FeatureId == 0 && bannerParams.TagId != 0 {
			searchBannersID = psql.Select("banner_id").
				From("tags").Where("tag_id = ?", bannerParams.FeatureId)
		} else {
			searchBannersID = psql.Select("ID").
				From("banners as b").
				Join("tags as t on b.id = t.banner_id").Join("features as f on b.id = f.banner_id").
				Where("feature_id = ? and tag_id = ? and is_active = true", bannerParams.FeatureId, bannerParams.TagId)
		}
		bannersQuery, args, err = searchBannersID.ToSql()
	}
	if err != nil {
		return bannerResultSlice, err
	}

	var bannerID int
	var bannerIDSlice []int
	foundID, err := bn.db.Query(bannersQuery, args...)
	if err != nil {
		return bannerResultSlice, err
	}
	defer foundID.Close()
	for foundID.Next() {
		foundID.Scan(&bannerID)
		bannerIDSlice = append(bannerIDSlice, bannerID)
	}

	selectAllBanners := psql.Select("id, content, is_active, create_time, update_time, feature_id").
		From("banners as b").
		Join("features as f on b.ID = f.banner_id")
	if len(bannerIDSlice) > 0 {
		selectAllBanners = selectAllBanners.Where(sq.Eq{"id": bannerIDSlice})
	}
	if bannerParams.Limit != 0 {
		selectAllBanners = selectAllBanners.Limit(uint64(bannerParams.Limit))
	}
	if bannerParams.Offset != 0 {
		selectAllBanners = selectAllBanners.Offset(uint64(bannerParams.Offset))
	}

	fetchAllBannersQuery, args, err := selectAllBanners.ToSql()
	if err != nil {
		return bannerResultSlice, err
	}
	bannerRows, err := bn.db.Query(fetchAllBannersQuery, args...)
	if err != nil {
		return bannerResultSlice, err
	}
	defer bannerRows.Close()
	for bannerRows.Next() {
		var banner BannerResponse
		err = bannerRows.Scan(&banner.ID, &banner.Content, &banner.IsActive,
			&banner.CreatedAt, &sql.NullTime{Time: banner.UpdatedAt}, &banner.FeatureId)

		if err != nil {
			return bannerResultSlice, err
		}
		selTags, args, err := psql.Select("tag_id").From("tags").Where("banner_id = ?", banner.ID).ToSql()
		if err != nil {
			return []BannerResponse{}, err
		}
		tagsRows, err := bn.db.Query(selTags, args...)
		if err != nil {
			return []BannerResponse{}, err
		}
		defer tagsRows.Close()
		var tagID int
		var tags []int
		for tagsRows.Next() {
			tagsRows.Scan(&tagID)
			tags = append(tags, tagID)
		}
		banner.TagId = tags

		bannerResultSlice = append(bannerResultSlice, banner)
	}
	return bannerResultSlice, nil
}

func (bn *BannerRepository) IfTokenValid(token string) (bool, error) {
	var exists bool

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("token").
		Prefix("SELECT EXISTS (").From("valid_tokens").
		Where("token = ?", token).Suffix(")").ToSql()
	if err != nil {
		return false, err
	}

	row := bn.db.QueryRow(query, args...)
	err = row.Scan(&exists)

	return exists, err

}

func (bn *BannerRepository) IfAdminTokenValid(token string) (bool, error) {
	var exists bool

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("token").
		Prefix("SELECT EXISTS (").From("valid_tokens").
		Where("token = ? and permission_level = 'admin'", token).Suffix(")").ToSql()
	if err != nil {
		return false, err
	}

	row := bn.db.QueryRow(query, args...)
	err = row.Scan(&exists)

	return exists, err

}

func (bn *BannerRepository) IfBannerExists(tagId int, featureId int) (bool, error) {
	var exists bool

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("content").
		Prefix("SELECT EXISTS (").
		From("banners as b").
		Join("tags as t on b.id = t.banner_id").Join("features as f on b.id = f.banner_id").
		Where("feature_id = ? and tag_id = ? and is_active = true", featureId, tagId).Suffix(")").ToSql()
	if err != nil {
		return false, err
	}
	row := bn.db.QueryRow(query, args...)
	err = row.Scan(&exists)
	return exists, err
}
