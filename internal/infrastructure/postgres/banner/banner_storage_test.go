package banner_storage

import (
	"database/sql"
	"log"
	"testing"

	"github.com/CyberPiess/banner_service/internal/infrastructure/logging"
	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

type GetUserBannerArgs struct {
	getUserBannerParams GetUserBannerCriteria
}

func TestGet(t *testing.T) {
	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "banner_storage_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		logger.Fatal("error init mock", err)
	}
	defer db.Close()

	bannerStorage := NewBannerRepository(db, logger)

	getUserBannerParams := GetUserBannerCriteria{
		TagId:     1,
		FeatureId: 1,
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, _, _ := psql.Select("content").
		From("banners as b").
		Join("tags as t on b.id = t.banner_id").Join("features as f on b.id = f.banner_id").
		Where("feature_id = ? and tag_id = ? and is_active = true").
		ToSql()
	rows := sqlmock.NewRows([]string{"content"}).AddRow(`{"some content":"content"}`)

	mock.ExpectQuery(query).WithArgs(getUserBannerParams.FeatureId, getUserBannerParams.TagId).WillReturnRows(rows)
	mock.ExpectQuery(query).WithArgs(getUserBannerParams.FeatureId, getUserBannerParams.TagId).WillReturnError(sql.ErrConnDone)

	tests := []struct {
		name       string
		args       GetUserBannerArgs
		wantAnswer []BannerEntitySql
		wantErr    error
	}{
		{
			name: "Correct data",
			args: GetUserBannerArgs{
				getUserBannerParams: GetUserBannerCriteria{
					TagId:     1,
					FeatureId: 1,
				},
			},
			wantAnswer: []BannerEntitySql{
				{
					Content: `{"some content":"content"}`,
				},
			},
			wantErr: nil,
		},
		{
			name: "Error on query",
			args: GetUserBannerArgs{
				getUserBannerParams: GetUserBannerCriteria{
					TagId:     1,
					FeatureId: 1,
				},
			},
			wantAnswer: []BannerEntitySql{},
			wantErr:    sql.ErrConnDone,
		},
	}

	for _, tt := range tests {

		foundContent, err := bannerStorage.Get(tt.args.getUserBannerParams)
		assert.Equal(t, tt.wantAnswer, foundContent)
		assert.Equal(t, tt.wantErr, err)
	}
}
func TestIfTokenValid(t *testing.T) {
	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "banner_storage_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		logger.Fatal("error init mock", err)
	}
	defer db.Close()

	bannerStorage := NewBannerRepository(db, logger)

	tokenString := "tokenString"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, _, _ := psql.Select("token").
		Prefix("SELECT EXISTS (").From("valid_tokens").
		Where("token = ?").Suffix(")").ToSql()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow("true")

	mock.ExpectQuery(query).WithArgs(tokenString).WillReturnRows(rows)
	mock.ExpectQuery(query).WithArgs(tokenString).WillReturnError(sql.ErrConnDone)

	tests := []struct {
		name      string
		token     string
		wantFound bool
		wantErr   error
	}{
		{
			name:      "Correct data",
			token:     "tokenString",
			wantFound: true,
			wantErr:   nil,
		},
		{
			name:      "Error on query",
			token:     "tokenString",
			wantFound: false,
			wantErr:   sql.ErrConnDone,
		},
	}

	for _, tt := range tests {

		foundContent, err := bannerStorage.IfTokenValid(tt.token)
		assert.Equal(t, tt.wantFound, foundContent)
		assert.Equal(t, tt.wantErr, err)
	}
}

func TestIfAdminTokenValid(t *testing.T) {
	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "banner_storage_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		logger.Fatal("error init mock", err)
	}
	defer db.Close()

	bannerStorage := NewBannerRepository(db, logger)

	tokenString := "tokenString"

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, _, _ := psql.Select("token").
		Prefix("SELECT EXISTS (").From("valid_tokens").
		Where("token = ? and permission_level = 'admin'").Suffix(")").ToSql()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow("true")

	mock.ExpectQuery(query).WithArgs(tokenString).WillReturnRows(rows)
	mock.ExpectQuery(query).WithArgs(tokenString).WillReturnError(sql.ErrConnDone)

	tests := []struct {
		name      string
		token     string
		wantFound bool
		wantErr   error
	}{
		{
			name:      "Correct data",
			token:     "tokenString",
			wantFound: true,
			wantErr:   nil,
		},
		{
			name:      "Error on query",
			token:     "tokenString",
			wantFound: false,
			wantErr:   sql.ErrConnDone,
		},
	}

	for _, tt := range tests {

		foundContent, err := bannerStorage.IfAdminTokenValid(tt.token)
		assert.Equal(t, tt.wantFound, foundContent)
		assert.Equal(t, tt.wantErr, err)
	}
}

func TestIfBannerExists(t *testing.T) {
	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "banner_storage_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		logger.Fatal("error init mock", err)
	}
	defer db.Close()

	bannerStorage := NewBannerRepository(db, logger)

	ifBannerExistParams := GetUserBannerCriteria{
		TagId:     1,
		FeatureId: 1,
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, _, _ := psql.Select("content").
		Prefix("SELECT EXISTS (").
		From("banners as b").
		Join("tags as t on b.id = t.banner_id").Join("features as f on b.id = f.banner_id").
		Where("feature_id = ? and tag_id = ? and is_active = true").Suffix(")").ToSql()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow("true")

	mock.ExpectQuery(query).WithArgs(ifBannerExistParams.FeatureId, ifBannerExistParams.TagId).WillReturnRows(rows)
	mock.ExpectQuery(query).WithArgs(ifBannerExistParams.FeatureId, ifBannerExistParams.TagId).WillReturnError(sql.ErrConnDone)

	tests := []struct {
		name      string
		args      GetUserBannerArgs
		wantFound bool
		wantErr   error
	}{
		{
			name: "Correct data",
			args: GetUserBannerArgs{
				getUserBannerParams: GetUserBannerCriteria{
					TagId:     1,
					FeatureId: 1,
				},
			},
			wantFound: true,
			wantErr:   nil,
		},
		{
			name: "Error on query",
			args: GetUserBannerArgs{
				getUserBannerParams: GetUserBannerCriteria{
					TagId:     1,
					FeatureId: 1,
				},
			},
			wantFound: false,
			wantErr:   sql.ErrConnDone,
		},
	}

	for _, tt := range tests {

		foundContent, err := bannerStorage.IfBannerExists(tt.args.getUserBannerParams.FeatureId, tt.args.getUserBannerParams.TagId)
		assert.Equal(t, tt.wantFound, foundContent)
		assert.Equal(t, tt.wantErr, err)
	}
}

func TestSearchBannerById(t *testing.T) {
	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "banner_storage_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		logger.Fatal("error init mock", err)
	}
	defer db.Close()

	bannerStorage := NewBannerRepository(db, logger)

	bannerID := 1

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query, _, _ := psql.Select("id").
		Prefix("SELECT EXISTS (").
		From("banners as b").
		Where("id = ?").Suffix(")").ToSql()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow("true")

	mock.ExpectQuery(query).WithArgs().WillReturnRows(rows)
	mock.ExpectQuery(query).WithArgs().WillReturnError(sql.ErrConnDone)

	tests := []struct {
		name      string
		bannerID  int
		wantFound bool
		wantErr   error
	}{
		{
			name:      "Correct data",
			bannerID:  bannerID,
			wantFound: true,
			wantErr:   nil,
		},
		{
			name:      "Error on query",
			bannerID:  bannerID,
			wantFound: false,
			wantErr:   sql.ErrConnDone,
		},
	}

	for _, tt := range tests {

		foundContent, err := bannerStorage.SearchBannerByID(tt.bannerID)
		assert.Equal(t, tt.wantFound, foundContent)
		assert.Equal(t, tt.wantErr, err)
	}
}
func TestGetAllBanners(t *testing.T) {
	logger, err := logging.LoggerCreate(logging.Config{LogLevel: "info",
		LogFile: "banner_storage_test.log"})
	if err != nil {
		log.Fatal("error init logger")
	}

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		logger.Fatal("error init mock", err)
	}
	defer db.Close()

	bannerStorage := NewBannerRepository(db, logger)

	selectBannerIDs := GetBannersListCriteria{
		TagId:     1,
		FeatureId: 1,
		Limit:     1,
		Offset:    1,
	}
	bannerID := 1
	timeFromDb := sql.NullTime{}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	selectBannerIDQuery, _, _ := psql.Select("ID").
		From("banners as b").
		Join("tags as t on b.id = t.banner_id").
		Join("features as f on b.id = f.banner_id").Where("feature_id = ?").Where("tag_id = ?").ToSql()

	returnedID := sqlmock.NewRows([]string{"id"}).AddRow(bannerID)
	returnedID2 := sqlmock.NewRows([]string{"id"}).AddRow(bannerID)
	returnedID3 := sqlmock.NewRows([]string{"id"}).AddRow(bannerID)

	selectBannerList, _, _ := psql.Select("id, content, is_active, create_time, update_time, feature_id").
		From("banners as b").
		Join("features as f on b.ID = f.banner_id").
		Where("id IN (?)").
		Limit(uint64(selectBannerIDs.Limit)).
		Offset(uint64(selectBannerIDs.Offset)).
		ToSql()
	returnedRow := sqlmock.NewRows([]string{"id", "content", "is_active", "create_time", "update_time", "feature_id"}).
		AddRow("1", `{"some_content":"string"}`, "true", timeFromDb.Time, timeFromDb.Time, "1")
	returnedRow2 := sqlmock.NewRows([]string{"id", "content", "is_active", "create_time", "update_time", "feature_id"}).
		AddRow("1", `{"some_content":"string"}`, "true", timeFromDb.Time, timeFromDb.Time, "1")

	selTags, _, _ := psql.Select("tag_id").From("tags").Where("banner_id = ?").ToSql()
	returnedTags := sqlmock.NewRows([]string{"tag_id"}).AddRow("1")

	mock.ExpectQuery(selectBannerIDQuery).WithArgs(selectBannerIDs.FeatureId, selectBannerIDs.TagId).WillReturnRows(returnedID)
	mock.ExpectQuery(selectBannerList).WithArgs(bannerID).WillReturnRows(returnedRow)
	mock.ExpectQuery(selTags).WithArgs(bannerID).WillReturnRows(returnedTags)

	mock.ExpectQuery(selectBannerIDQuery).WithArgs(selectBannerIDs.FeatureId, selectBannerIDs.TagId).WillReturnRows(returnedID2)
	mock.ExpectQuery(selectBannerList).WithArgs(bannerID).WillReturnRows(returnedRow2)
	mock.ExpectQuery(selTags).WithArgs(bannerID).WillReturnError(sql.ErrConnDone)

	mock.ExpectQuery(selectBannerIDQuery).WithArgs(selectBannerIDs.FeatureId, selectBannerIDs.TagId).WillReturnRows(returnedID3)
	mock.ExpectQuery(selectBannerList).WithArgs(bannerID).WillReturnError(sql.ErrConnDone)

	mock.ExpectQuery(selectBannerIDQuery).WithArgs(selectBannerIDs.FeatureId, selectBannerIDs.TagId).WillReturnError(sql.ErrConnDone)

	tests := []struct {
		name                string
		getBannerListParams GetBannersListCriteria
		wantFound           []BannerEntitySql
		wantErr             error
	}{
		{
			name:                "Correct data",
			getBannerListParams: selectBannerIDs,
			wantFound: []BannerEntitySql{
				{
					ID:        1,
					Content:   `{"some_content":"string"}`,
					IsActive:  true,
					CreatedAt: timeFromDb.Time,
					UpdatedAt: timeFromDb.Time,
					TagId:     []int{1},
					FeatureId: 1,
				},
			},
			wantErr: nil,
		},
		{
			name:                "Error while getting tags",
			getBannerListParams: selectBannerIDs,
			wantFound:           []BannerEntitySql{},
			wantErr:             sql.ErrConnDone,
		},
		{
			name:                "Error while getting all banner info",
			getBannerListParams: selectBannerIDs,
			wantFound:           []BannerEntitySql{},
			wantErr:             sql.ErrConnDone,
		},
		{
			name:                "Error while getting all banners ID",
			getBannerListParams: selectBannerIDs,
			wantFound:           []BannerEntitySql{},
			wantErr:             sql.ErrConnDone,
		},
	}

	for _, tt := range tests {

		foundContent, err := bannerStorage.GetAllBanners(tt.getBannerListParams)
		assert.Equal(t, tt.wantFound, foundContent)
		assert.Equal(t, tt.wantErr, err)
	}
}
