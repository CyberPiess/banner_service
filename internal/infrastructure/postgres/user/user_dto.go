package user

type UserRequest struct {
	TagId           int
	FeatureId       int
	UseLastRevision bool
}

type UserResponse struct {
	Content string
}
