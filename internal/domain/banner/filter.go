package banner

type Filter struct {
	TagId           string `json:"tag_id"`
	FeatureId       string `json:"feature_id"`
	UseLastRevision string `json:"use_last_revision"`
}
