package usecase

import (
	"context"
	"errors"
	"time"

	"login-service/login"
	"login-service/login/dto"
	"login-service/models"
	"login-service/utils/password"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type LoginUsecase struct {
	loginRepo      login.Repository
	cacheRepo      login.CacheRepository
	contextTimeout time.Duration
}

// LoginUsecase 结构的指针（注意：是他的指针，不是这个结构体）实现了  LoginUsecase 接口
func NewLoginUsecase(o login.Repository, c login.CacheRepository, timeout time.Duration) *LoginUsecase {
	var et time.Duration
	if timeout.Seconds() <= 0 {
		et = time.Duration(5 * time.Second)
	}
	return &LoginUsecase{
		loginRepo:      o,
		cacheRepo:      c,
		contextTimeout: et,
	}
}

func (m *LoginUsecase) Login(ctx context.Context, login dto.LoginDTO) (models.UserInfo, error) {

	var user *models.UserInfo
	if login.Phone == "" && login.Email == "" {
		logrus.Info("用户名不能为空")
		return models.UserInfo{}, models.NewErr("用户/密码不正确")
	}

	if login.Phone == "" {
		login.Phone = login.Email
	}
	// 根据 phone 或 email 查询 user
	user, err := m.loginRepo.FindByPhone(login.Phone)
	if err != nil && err == gorm.ErrRecordNotFound {
		// 查询出错则返回错误
		logrus.Error(err.Error())
		return *user, err
	}

	err = checkPasswd(user.PasswordDigest, login.Password)
	if err != nil {
		return models.UserInfo{}, err
	}

	return *user, nil
}

// 验证密码
func checkPasswd(ercpasswd, passwd string) error {
	if len(passwd) < 6 || len(passwd) > 60 {
		logrus.Info("密码太长/太短")
		return models.NewErr("用户/密码不正确")
	}
	// 验证密码，用的加密包 golang.org/x/crypto/bcrypt
	if !password.Compare(ercpasswd, passwd) {
		return models.NewErr("用户/密码不正确")
	}
	return nil
}

var (
	// 特权列表
	superuser = map[string]bool{
		"18361273739": true,
		"12345678901": true,
		"13814842044": true,
	}
)

// 根据用户信息生成 jwt token, 其实可以把解析token获取用户信息，但我们只是作为redis的key
func (m *LoginUsecase) GetJWT(ctx context.Context, user models.UserInfo) (string, error) {

	exp := time.Now().Add(time.Hour * 24 * 7).Unix()
	claim := jwt.MapClaims{
		"exp": exp,
		"sub": user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	jwtStr, err := token.SignedString([]byte(viper.GetString("jwt_signed")))
	if err != nil {
		return "", errors.New("登录申请失败")
	}
	userFields := user.ConvertSessionFields()
	delete(userFields, "password_digest")

	logrus.Info("jwt created: ", jwtStr)
	logrus.Info("jwt fields: ", userFields)
	// token 写进redis   7 天过期
	err = m.cacheRepo.SetSession(jwtStr, userFields, 7*86400)
	if err != nil {
		logrus.Error("登录失败:", userFields, err)
		return "", errors.New("登录失败")
	}
	// 更新用户最后一次登录时间
	go m.loginRepo.UpdateUserLastLogin(user.ID)
	return jwtStr, err
}

// 注销 token
func (m *LoginUsecase) Logout(token string) error {
	// 从redis 删除 token
	err := m.cacheRepo.Delete(token)
	return err
}
