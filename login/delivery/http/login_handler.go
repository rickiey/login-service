package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"login-service/helper"
	"login-service/login/dto"
	"login-service/models"
	"net/http"
	"time"
)

func (m *LoginHttpHandler) Middleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Token")
		w.Header().Set("Access-Control-Expose-Headers", "Token")
		next.ServeHTTP(w, r)
	})
}

// 登录 v1
func (m *LoginHttpHandler) login(e echo.Context) error {

	var login dto.LoginDTO
	var resp models.Response

	// 接受数据json解析并绑定到结构体上，这里传的指针
	err := e.Bind(&login)
	if err != nil {
		logrus.Errorf("login.无法解析数据:%v", err)
		resp.Failed("无法解析数据")
		return e.JSON(http.StatusOK, resp)
	}

	// 验证参数
	err = helper.IsValid(login)
	if err != nil {
		logrus.Errorf("login.验证数据:%+v,失败:%v\n", login, err.Error())
		resp.Failed(helper.ValidateError(err))
		return e.JSON(http.StatusOK, resp)
	}
	fmt.Println(login)

	user, err := m.usecase.Login(nil, login)
	if err != nil {
		merr := new(models.Err)
		// 这个 merr 是个指针
		if errors.As(err, &merr) {
			resp.Failed(merr.Msg)
			return e.JSON(http.StatusOK, resp)
		}
		logrus.Errorf("login.Login:%+v,失败:%v\n", login, err.Error())
		resp.Failed("internal error")
		return e.JSON(http.StatusOK, resp)
	}

	jwt, err := m.usecase.GetJWT(context.Background(), user)
	if err != nil {
		merr := new(models.Err)
		// 判断这个错误是不是自定义的错误（models.Err），这个 merr 是个指针
		if errors.As(err, &merr) {
			resp.Failed(merr.Msg)
			return e.JSON(http.StatusOK, resp)
		}
		logrus.Errorf("login.GetJWT:%+v,失败:%v\n", login, err.Error())
		resp.Failed(err.Error())
		return e.JSON(http.StatusOK, resp)
	}

	userFields := user.RemoveSensitive()

	token := map[string]interface{}{
		"jwt":  jwt,
		"user": userFields,
	}
	// 设置cookie
	e.SetCookie(&http.Cookie{Name: "access_token", MaxAge: 86400 * 7, Expires: time.Now().AddDate(0, 0, 7), Value: jwt, Path: "/"})
	// 返回的 201 ， json格式
	return e.JSON(http.StatusCreated, token)
}

// logout
// 注销接口， 删除 token
func (m *LoginHttpHandler) logout(e echo.Context) error {
	var resp models.APIResponse

	token := helper.GetToken(e.Request())
	err := m.usecase.Logout(token)
	if err != nil {
		resp.InternalError(err.Error())
		return e.JSON(http.StatusInternalServerError, resp)
	}
	resp.Success(nil)
	return e.JSON(http.StatusOK, resp)
}
