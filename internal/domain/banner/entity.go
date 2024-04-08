package banner

import "time"

type BannerEntity struct {
	ID        int
	Content   map[string]interface{} `json:"content"`
	TagId     []int                  `json:"tag_ids"`
	FeatureId int                    `json:"feature_id"`
	IsActive  bool                   `json:"is_active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
