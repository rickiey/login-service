package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"login-service/login"
)

const (
	microServicePath = "/account"
)

// LoginHttpHandler
type LoginHttpHandler struct {
	usecase login.Usecase
}

func NewHttpLoginHandler(uc login.Usecase) *LoginHttpHandler {
	return &LoginHttpHandler{uc}
}

func (m *LoginHttpHandler) InitApi(e *echo.Echo) {

	// 组控制，统一增加前缀
	g := e.Group(microServicePath)
	g.POST("/v1/login", m.login)
	g.GET("/v1/logout", m.logout)
	g.POST("/v1/password/update", m.updatePassword)

	// 可跨域配置
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

}
