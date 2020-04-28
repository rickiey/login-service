package http

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"login-service/helper"
	"login-service/login/dto"
	"login-service/models"
	"net/http"
)

// updatePassword
// 更改密码， 需要手机号，原密码，新密码
func (m *LoginHttpHandler) updatePassword(e echo.Context) error {

	var changePass dto.ChangePassword
	var resp models.Response
	err := e.Bind(&changePass)
	if err != nil {
		logrus.Errorf("login.无法解析数据:%v", err)
		resp.Failed("无法解析数据")
		return e.JSON(http.StatusBadRequest, resp)
	}

	err = helper.IsValid(changePass)
	if err != nil {
		logrus.Errorf("login.验证数据:%+v,失败:%v\n", changePass, err.Error())
		resp.Failed(helper.ValidateError(err))
		return e.JSON(http.StatusBadRequest, resp)
	}

	changePass.Token = helper.GetToken(e.Request())

	err = m.usecase.ChangePassword(e.Request().Context(), changePass)

	if err != nil {
		// 出错失败返回
		logrus.Errorf("outLogin.Login:%+v,失败:%v\n", changePass, err.Error())
		resp.Failed(err.Error())
		return e.JSON(http.StatusOK, resp)
	}
	resp.Success(nil)
	return e.JSON(http.StatusOK, resp)
}
