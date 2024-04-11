package banner

type GetFilter struct {
	TagId           int
	FeatureId       int
	UseLastRevision bool
}

type GetAllFilter struct {
	TagId     int
	FeatureId int
	Limit     int
	Offset    int
}
