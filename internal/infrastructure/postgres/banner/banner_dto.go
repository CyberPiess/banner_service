package banner

import "time"

type BannerRequest struct {
	TagId           int
	FeatureId       int
	UseLastRevision bool
	Limit           int
	Offset          int
}

type BannerResponse struct {
	ID        int    `db:"id"`
	Content   string `db:"content"`
	TagId     []int
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"create_time"`
	UpdatedAt time.Time `db:"update_time"`
	FeatureId int       `db:"feature_id"`
}
