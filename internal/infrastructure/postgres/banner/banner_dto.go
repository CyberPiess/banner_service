package banner

type BannerRequest struct {
	TagId           int
	FeatureId       int
	UseLastRevision bool
}

type BannerResponse struct {
	Content string
}
