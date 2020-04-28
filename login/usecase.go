package login

import (
	"context"

	"login-service/login/dto"
	"login-service/models"
)

//Usecase describe login usecase
type Usecase interface {
	Login(ctx context.Context, user dto.LoginDTO) (models.UserInfo, error)
	Logout(token string) error
	GetJWT(ctx context.Context, user models.UserInfo) (string, error)
	ChangePassword(ctx context.Context, user dto.ChangePassword) error
}
