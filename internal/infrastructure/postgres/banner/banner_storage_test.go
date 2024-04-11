package banner

import (
	"database/sql"
	"log"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
)

//	type bannerStorage interface {
//		Get(bannerParams banner.BannerCriteria) ([]banner.BannerEntitySql, error)
//		IfTokenValid(token string) (bool, error)
//		IfBannerExists(featureId int, tagId int) (bool, error)
//		IfAdminTokenValid(token string) (bool, error)
//		SearchBannerByID(bannerID int) (bool, error)
//		GetAllBanners(bannerParams banner.BannerCriteria) ([]banner.BannerEntitySql, error)
//		PostBanner(postBannerParams banner.BannerPutPostCriteria) (int, error)
//		PutBanner(putBannerParams banner.BannerPutPostCriteria) error
//		DeleteBanner(deleteBannerParams banner.BannerPutPostCriteria) error
//	}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

type GetUserBannerArgs struct {
	getUserBannerParams GetUserBannerCriteria
}

func TestGet(t *testing.T) {
	db, mock := NewMock()

	bannerStorage := NewBannerRepository(db)

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
					TagId:           1,
					FeatureId:       1,
					UseLastRevision: false,
				},
			},
			wantAnswer: []BannerEntitySql{
				{
					Content: `{"some content":"content"}`,
				},
			},
			wantErr: nil,
		},
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	for _, tt := range tests {

		query, args, _ := psql.Select("content").
			From("banners as b").
			Join("tags as t on b.id = t.banner_id").Join("features as f on b.id = f.banner_id").
			Where("feature_id = ? and tag_id = ? and is_active = true", tt.args.getUserBannerParams.FeatureId, tt.args.getUserBannerParams.TagId).
			ToSql()

		//query := "SELECT content FROM banners as b JOIN tags as t on b.id = t.banner_id JOIN features as f on b.id = f.banner_id WHERE feature_id = \\? and tag_id = \\? and is_active = true;"
		rows := sqlmock.NewRows([]string{"content"}).AddRow(`{"some content":"content"}`)

		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args).WillReturnRows(rows)
		foundContent, err := bannerStorage.Get(tt.args.getUserBannerParams)
		assert.Equal(t, tt.wantAnswer, foundContent)
		assert.Equal(t, tt.wantErr, err)
	}

}
