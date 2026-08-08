// Copyright (c) 2019. icemsoft.net
// Author: Bruce Created:2019/1/15

package helper

import (
	"log"
	"net/url"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const RecordNotFound = "record not found"

func GetDsnFromConfig() string {

	loc, err := time.LoadLocation(viper.GetString("mysql_location"))
	if err != nil {
		log.Panicf("Time.LoadLocation by %s failed error: %s", viper.GetString("mysql_location"), err.Error())
		return ""
	}

	o := mysql.Config{
		User:                 viper.GetString("mysql_user"),
		Passwd:               viper.GetString("mysql_password"),
		Net:                  viper.GetString("mysql_net"),
		Addr:                 viper.GetString("mysql_host"),
		DBName:               viper.GetString("mysql_name"),
		Collation:            viper.GetString("mysql_collation"),
		Loc:                  loc,
		ParseTime:            viper.GetBool("mysql_parse_time"),
		AllowNativePasswords: true,
	}

	return o.FormatDSN() + "&" + url.PathEscape(viper.GetString("mysql_params"))
}

func NewDBConn(dsn string) (*gorm.DB, error) {
	DBConn, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("gorm.Open by %s failed error: %s", dsn, err.Error())
		return DBConn, err
	}
	sqlDB, err := DBConn.DB()
	if err != nil {
		log.Panicf("gorm.DB() failed error: %s", err.Error())
		return DBConn, err
	}
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Second)
	if sqlDB.Ping() != nil {
		log.Panicf("gorm.Open by %s failed error: %s", dsn, err.Error())
		return DBConn, err
	}
	return DBConn, nil
}

// BatchInsert 批量插入数据。使用 gorm v2 的 CreateInBatches 实现。
func BatchInsert(db *gorm.DB, objArr []interface{}) error {
	// If there is no data, nothing to do.
	if len(objArr) == 0 {
		return nil
	}
	return db.CreateInBatches(objArr, len(objArr)).Error
}
