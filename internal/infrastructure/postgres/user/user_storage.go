package user

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Get(bannerParams UserRequest) (UserResponse, error) {
	var content string

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, args, err := psql.Select("content").
		From("banners").
		Where("feature_id = ? and (tag_id && '{?}')", bannerParams.FeatureId, bannerParams.TagId).ToSql()
	if err != nil {
		return UserResponse{}, err
	}

	row := ur.db.QueryRow(query, args...)
	err = row.Scan(&content)
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{Content: content}, nil
}
