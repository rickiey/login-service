package usecase

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"login-service/login/dto"
	"login-service/utils"
	"login-service/utils/password"
)

// 更改密码
func (m *LoginUsecase) ChangePassword(ctx context.Context, changePass dto.ChangePassword) error {
	var userId int64
	if changePass.Token == "" {
		return errors.New("无效的token")
	}

	sessionFields, err := m.cacheRepo.GetSession(changePass.Token)
	if err != nil {
		return err
	}

	uId, ok := sessionFields["id"]
	if !ok {
		return errors.New("token值无效")
	}

	userId, _ = utils.ToInt64(uId)
	if userId <= 0 {
		return errors.New("token id值无效")
	}
	//user info
	info, err := m.loginRepo.FindByUserId(userId)
	if err != nil {
		return err
	}
	if info.ID <= 0 {
		return errors.New("用户不存在")
	}
	logrus.Info("login user:", info)

	if !password.Compare(info.PasswordDigest, changePass.Password) {
		return errors.New("用户/密码不正确")
	}
	if changePass.NewPassword != changePass.ConfirmPassword {
		return errors.New("确认密码不正确")
	}


	newPass, err := password.Get(changePass.NewPassword)
	if err != nil {
		return errors.New("修改密码,发生错误")
	}
	err = m.loginRepo.UpdatePassword(newPass, userId)
	if err != nil {
		return errors.New("修改密码,发生错误")
	}
	return nil
}
