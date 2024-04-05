package user

import "github.com/CyberPiess/banner_sevice/internal/infrastructure/postgres/user"

type userStorage interface {
	Get(bannerParams user.UserRequest) (user.UserResponse, error)
}

type UserService struct {
	store userStorage
}

func NewUserService(storage userStorage) *UserService {
	return &UserService{store: storage}
}

func (u *UserService) GetBanner(userParams UserParams) (UserResponse, error) {
	//TODO: добавить проверки

	userRequest := user.UserRequest{
		TagId:           userParams.TagId,
		FeatureId:       userParams.FeatureId,
		UseLastRevision: userParams.UseLastRevision,
	}

	response, err := u.store.Get(userRequest)
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{Content: response.Content}, nil
}
