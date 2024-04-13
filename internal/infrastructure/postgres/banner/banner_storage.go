//go:generate mockgen -source=banner_storage.go -destination=mocks/mock.go
package banner_storage

import (
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
)

type logger interface {
	WithFields(fields logrus.Fields) *logrus.Entry
}

type BannerRepository struct {
	db     *sql.DB
	logger logger
}

func NewBannerRepository(db *sql.DB, logger logger) *BannerRepository {
	return &BannerRepository{db: db, logger: logger}
}

func (bn *BannerRepository) Get(bannerParams GetUserBannerCriteria) ([]BannerEntitySql, error) {
	var banner BannerEntitySql

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.Select("content").
		From("banners as b").
		Join("tags as t on b.id = t.banner_id").Join("features as f on b.id = f.banner_id").
		Where("feature_id = ? and tag_id = ? and is_active = true", bannerParams.FeatureId, bannerParams.TagId).ToSql()
	if err != nil {
		bn.logger.WithFields(logrus.Fields{
			"package":  "banner_storage",
			"function": "Get",
			"error":    err,
		}).Error("Error creating select statement")
		return []BannerEntitySql{}, err
	}

	row := bn.db.QueryRow(query, args...)
	err = row.Scan(&banner.Content)
	if err != nil {
		bn.logger.WithFields(logrus.Fields{
			"package":  "banner_storage",
			"function": "Get",
			"error":    err,
		}).Error("Error scanning row")
		return []BannerEntitySql{}, err
	}

	return []BannerEntitySql{banner}, nil
}

func (bn *BannerRepository) GetAllBanners(bannerParams GetBannersListCriteria) ([]BannerEntitySql, error) {

	var bannerResultSlice []BannerEntitySql
	var bannersQuery string
	var args []interface{}
	var err error

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	searchBannersID := psql.Select("ID").
		From("banners as b").
		Join("tags as t on b.id = t.banner_id").Join("features as f on b.id = f.banner_id")

	if bannerParams.FeatureId != 0 {
		searchBannersID = searchBannersID.Where("feature_id = ?", bannerParams.FeatureId)
	}
	if bannerParams.TagId != 0 {
		searchBannersID = searchBannersID.Where("tag_id = ?", bannerParams.TagId)
	}

	bannersQuery, args, err = searchBannersID.ToSql()

	if err != nil {
		bn.logger.WithFields(logrus.Fields{
			"package":  "banner_storage",
			"function": "GetAllBanners",
			"error":    err,
		}).Error("Error while creating select statement bannersQuery")
		return []BannerEntitySql{}, err
	}

	var bannerID int
	var bannerIDSlice []int
	foundID, err := bn.db.Query(bannersQuery, args...)
	if err != nil {
		bn.logger.WithFields(logrus.Fields{
			"package":  "banner_storage",
			"function": "GetAllBanners",
			"error":    err,
		}).Error("Error while running query for foundID")
		return []BannerEntitySql{}, err
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
	} else {
		return []BannerEntitySql{}, nil
	}
	if bannerParams.Limit != 0 {
		selectAllBanners = selectAllBanners.Limit(uint64(bannerParams.Limit))
	}
	if bannerParams.Offset != 0 {
		selectAllBanners = selectAllBanners.Offset(uint64(bannerParams.Offset))
	}

	fetchAllBannersQuery, args, err := selectAllBanners.ToSql()
	if err != nil {
		bn.logger.WithFields(logrus.Fields{
			"package":  "banner_storage",
			"function": "GetAllBanners",
			"error":    err,
		}).Error("Error while creating select statement for fetchAllBannersQuery")
		return []BannerEntitySql{}, err
	}
	bannerRows, err := bn.db.Query(fetchAllBannersQuery, args...)
	if err != nil {
		bn.logger.WithFields(logrus.Fields{
			"package":  "banner_storage",
			"function": "GetAllBanners",
			"error":    err,
		}).Error("Error while running query for bannerRows")
		return []BannerEntitySql{}, err
	}
	defer bannerRows.Close()

	for bannerRows.Next() {
		var banner BannerEntitySql
		var updateTime sql.NullTime
		err = bannerRows.Scan(&banner.ID, &banner.Content, &banner.IsActive,
			&banner.CreatedAt, &updateTime, &banner.FeatureId)

		if updateTime.Valid {
			banner.UpdatedAt = updateTime.Time
		}
		if err != nil {
			bn.logger.WithFields(logrus.Fields{
				"package":  "banner_storage",
				"function": "GetAllBanners",
				"error":    err,
			}).Error("Error while scanning bannerRows")
			return []BannerEntitySql{}, err
		}
		selTags, args, err := psql.Select("tag_id").From("tags").Where("banner_id = ?", banner.ID).ToSql()
		if err != nil {
			bn.logger.WithFields(logrus.Fields{
				"package":  "banner_storage",
				"function": "GetAllBanners",
				"error":    err,
			}).Error("Error while creating select statement for selTags")
			return []BannerEntitySql{}, err
		}
		tagsRows, err := bn.db.Query(selTags, args...)
		if err != nil {
			bn.logger.WithFields(logrus.Fields{
				"package":  "banner_storage",
				"function": "GetAllBanners",
				"error":    err,
			}).Error("Error while running query for tagsRow")
			return []BannerEntitySql{}, err
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

func (bn *BannerRepository) PostBanner(postBannerParams BannerPutPostCriteria) (int, error) {
	var createdID int
	tx, err := bn.db.Begin()

	if err != nil {
		bn.logger.WithFields(logrus.Fields{
			"package":  "banner_storage",
			"function": "PostBanner",
			"error":    err,
		}).Error("Error while  beginning transaction")
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			bn.logger.WithFields(logrus.Fields{
				"package":  "banner_storage",
				"function": "PostBanner",
				"error":    err,
			}).Error("Error while inserting")
			return
		}
		err = tx.Commit()
	}()

	insertBannerQuery := sq.Insert("banners").
		Columns("content", "is_active", "create_time").
		Values(postBannerParams.Content,
			postBannerParams.IsActive,
			postBannerParams.CreatedAt).Suffix("returning id").RunWith(tx).
		PlaceholderFormat(sq.Dollar)

	err = insertBannerQuery.QueryRow().Scan(&createdID)

	_, err = sq.Insert("features").
		Columns("feature_id", "banner_id").
		Values(postBannerParams.FeatureId, createdID).RunWith(tx).PlaceholderFormat(sq.Dollar).Exec()
	for _, tagID := range postBannerParams.TagIds {
		_, err = sq.Insert("tags").
			Columns("tag_id", "banner_id").
			Values(tagID, createdID).RunWith(tx).PlaceholderFormat(sq.Dollar).Exec()
	}

	return createdID, err
}

func (bn *BannerRepository) PutBanner(putBannerParams BannerPutPostCriteria) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	tx, err := bn.db.Begin()

	if err != nil {
		bn.logger.WithFields(logrus.Fields{
			"package":  "banner_storage",
			"function": "PutBanner",
			"error":    err,
		}).Error("Error while beginning transaction")
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			bn.logger.WithFields(logrus.Fields{
				"package":  "banner_storage",
				"function": "PutBanner",
				"error":    err,
			}).Error("Error while updating")
			return
		}
		err = tx.Commit()
	}()

	updateBannerContent := psql.Update("banners").Where("ID = ?", putBannerParams.ID).
		Set("update_time", time.Now())
	if putBannerParams.IfFlagActiveIsSet {
		updateBannerContent = updateBannerContent.Set("is_active", putBannerParams.IsActive)
	}
	updateBannerContent = updateBannerContent.RunWith(tx)
	_, err = updateBannerContent.Exec()

	if len(putBannerParams.TagIds) > 0 {
		deleteTags := psql.Delete("tags").Where("banner_id = ?", putBannerParams.ID).RunWith(tx)
		_, err = deleteTags.Exec()

		for _, tagID := range putBannerParams.TagIds {
			_, err = sq.Insert("tags").
				Columns("tag_id", "banner_id").
				Values(tagID, putBannerParams.ID).RunWith(tx).PlaceholderFormat(sq.Dollar).Exec()
		}
	}

	if putBannerParams.FeatureId != 0 {
		updateFeatureID := psql.Update("features").Where("banner_id = ?", putBannerParams.ID).
			Set("feature_id", putBannerParams.FeatureId).RunWith(tx)
		_, err = updateFeatureID.Exec()
	}

	return err
}

func (bn *BannerRepository) DeleteBanner(putBannerParams BannerPutPostCriteria) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	tx, err := bn.db.Begin()

	if err != nil {
		bn.logger.WithFields(logrus.Fields{
			"package":  "banner_storage",
			"function": "DeleteBanner",
			"error":    err,
		}).Error("Error while beginning transaction")
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			bn.logger.WithFields(logrus.Fields{
				"package":  "banner_storage",
				"function": "DeleteBanner",
				"error":    err,
			}).Error("Error while deleting")
			return
		}
		err = tx.Commit()
	}()

	updateBannerContent := psql.Delete("banners").Where("ID = ?", putBannerParams.ID).RunWith(tx)
	_, err = updateBannerContent.Exec()

	deleteTags := psql.Delete("tags").Where("banner_id = ?", putBannerParams.ID).RunWith(tx)
	_, err = deleteTags.Exec()

	updateFeatureID := psql.Delete("features").Where("banner_id = ?", putBannerParams.ID).RunWith(tx)
	_, err = updateFeatureID.Exec()

	return err
}

func (bn *BannerRepository) IfTokenValid(token string) (bool, error) {
	var exists bool

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("token").
		Prefix("SELECT EXISTS (").From("valid_tokens").
		Where("token = ?", token).Suffix(")").ToSql()
	if err != nil {
		bn.logger.WithFields(logrus.Fields{
			"package":  "banner_storage",
			"function": "IfTokenvalid",
			"error":    err,
		}).Error("Error while creating select statement")
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
		bn.logger.WithFields(logrus.Fields{
			"package":  "banner_storage",
			"function": "IfAdminTokenValid",
			"error":    err,
		}).Error("Error while creating select statement")
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
		bn.logger.WithFields(logrus.Fields{
			"package":  "banner_storage",
			"function": "IfBannerExists",
			"error":    err,
		}).Error("Error while creating select statement")
		return false, err
	}
	row := bn.db.QueryRow(query, args...)
	err = row.Scan(&exists)
	return exists, err
}

func (bn *BannerRepository) SearchBannerByID(bannerID int) (bool, error) {
	var exists bool

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	query, args, err := psql.Select("id").
		Prefix("SELECT EXISTS (").
		From("banners as b").
		Where("id = ?", bannerID).Suffix(")").ToSql()
	if err != nil {
		bn.logger.WithFields(logrus.Fields{
			"package":  "banner_storage",
			"function": "SearchBannerByID",
			"error":    err,
		}).Error("Error while creating select statement")
		return false, err
	}
	row := bn.db.QueryRow(query, args...)
	err = row.Scan(&exists)
	return exists, err
}
