package user

type UserParams struct {
	TagId           int  `json:"tag_id"`
	FeatureId       int  `json:"feature_id"`
	UseLastRevision bool `json:"use_last_revision"`
	Token           string
}

type UserResponse struct {
	Content string
}
