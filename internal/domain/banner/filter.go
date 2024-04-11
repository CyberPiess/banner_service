package banner

type GetFilter struct {
	TagId           int
	FeatureId       int
	UseLastRevision bool
}

type GetAllFilter struct {
	TagId     int `schema:"tag_id"`
	FeatureId int `schema:"feature_id"`
	Limit     int `schema:"limit"`
	Offset    int `schema:"offset"`
}
