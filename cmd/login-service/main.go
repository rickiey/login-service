package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	_helper "login-service/helper"
	_deliveryHttp "login-service/login/delivery/http"
	_repo "login-service/login/repository"
	_uc "login-service/login/usecase"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})

	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.InfoLevel)
	viper.AddConfigPath("/etc/conf/")
	viper.AddConfigPath("conf")
	viper.AddConfigPath("$HOME/")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

}

func main() {

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Warnf("config file not found/read : %s \n", err)
	}
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)

	//OrgDsn := "root:windows2010..A@tcp(192.168.8.120:3306)/organization_develop?charset=utf8&parseTime=True&loc=Local"
	OrgDsn := _helper.GetDsnFromConfig()

	fmt.Println(OrgDsn)
	fmt.Println(viper.GetString("cache_address"))
	fmt.Println(viper.GetString("jwt_signed"))
	dbConn, _ := _helper.NewDBConn(OrgDsn)

	opt := _helper.RedisOption{
		Address:  viper.GetString("cache_address"),
		Password: viper.GetString("cache_password"),
	}
	cacheRepo := _repo.NewSessionRepository(opt)

	authorRepo := _repo.NewMysqlUserRepository(dbConn)
	orgUC := _uc.NewLoginUsecase(authorRepo, cacheRepo, 5*time.Second)
	handler := _deliveryHttp.NewHttpLoginHandler(orgUC)
	e := echo.New()
	s := &http.Server{
		Addr: viper.GetString("http_address"),
	}
	handler.InitApi(e)
	e.HideBanner = true

	idleConnClosed := make(chan struct{})

	// 优雅关闭
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint,
			os.Kill,
			os.Interrupt,
			syscall.SIGHUP, // interrupt signal sent from terminal
			syscall.SIGINT,
			syscall.SIGTERM, // sigterm signal sent from kubernetes
			syscall.SIGQUIT)

		// We received an interrupt signal, shut down.
		select {
		case <-sigint:
			logrus.Info("server going to shut down")
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
			defer cancel()

			err := s.Shutdown(ctx)
			if err != nil {
				fmt.Println("shutdown ....", err)
			}
			logrus.Info("presenter.server close done!")
		}

		close(idleConnClosed)
	}()

	e.Logger.Fatal(e.StartServer(s))
	<-idleConnClosed
	logrus.Info("server done!")
}
