package banner

type GetFilter struct {
	TagId           int  `schema:"tag_id,required"`
	FeatureId       int  `schema:"feature_id,required"`
	UseLastRevision bool `schema:"use_last_revision,default:false"`
}

type GetAllFilter struct {
	TagId     int `schema:"tag_id"`
	FeatureId int `schema:"feature_id"`
	Limit     int `schema:"limit"`
	Offset    int `schema:"offset"`
}
