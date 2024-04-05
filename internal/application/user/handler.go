package user

import (
	"net/http"
	"strconv"

	"github.com/CyberPiess/banner_sevice/internal/domain/user"
)

type userService interface {
	GetBanner(userParams user.UserParams) (user.UserResponse, error)
}

type User struct {
	service userService
}

func NewUserHandler(service userService) *User {
	return &User{service: service}
}

func (u *User) GetUserBanner(w http.ResponseWriter, r *http.Request) {
	var userParams user.UserParams
	userParams.TagId, _ = strconv.Atoi(r.FormValue("tag_id"))
	userParams.FeatureId, _ = strconv.Atoi(r.FormValue("feature_id"))
	userParams.UseLastRevision, _ = strconv.ParseBool(r.FormValue("use_last_revision"))
	userParams.Token = r.Header.Get("token")

	u.service.GetBanner(userParams)
}
