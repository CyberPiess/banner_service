package banner

import "time"

type BannerEntity struct {
	ID        int
	Content   map[string]interface{}
	TagId     []int
	FeatureId int
	IsActive  *bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
