// Copyright (c) 2019. icemsoft.net
// Author: Bruce Created:2019/1/15

package helper

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"strings"
	"time"
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
	DBConn, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Panicf("gorm.Open by %s failed error: %s", dsn, err.Error())
		return DBConn, err
	}
	DBConn.DB().SetMaxOpenConns(200)
	DBConn.DB().SetMaxIdleConns(10)
	DBConn.DB().SetConnMaxLifetime(30 * time.Second)
	if DBConn.DB().Ping() != nil {
		log.Panicf("gorm.Open by %s failed error: %s", dsn, err.Error())
		return DBConn, err
	}
	return DBConn, nil
}

func BatchInsert(db *gorm.DB, objArr []interface{}) error {
	// If there is no data, nothing to do.
	if len(objArr) == 0 {
		return nil
	}

	mainObj := objArr[0]
	mainScope := db.NewScope(mainObj)
	mainFields := mainScope.Fields()
	quoted := make([]string, 0, len(mainFields))
	for i := range mainFields {
		// If primary key has blank value (0 for int, "" for string, nil for interface ...), skip it.
		// If field is ignore field, skip it.
		if (mainFields[i].IsPrimaryKey && mainFields[i].IsBlank) || (mainFields[i].IsIgnored) {
			continue
		}
		quoted = append(quoted, mainScope.Quote(mainFields[i].DBName))
	}

	placeholdersArr := make([]string, 0, len(objArr))

	for _, obj := range objArr {
		scope := db.NewScope(obj)
		fields := scope.Fields()
		placeholders := make([]string, 0, len(fields))
		for i := range fields {
			if (fields[i].IsPrimaryKey && fields[i].IsBlank) || (fields[i].IsIgnored) {
				continue
			}
			placeholders = append(placeholders, scope.AddToVars(fields[i].Field.Interface()))
		}
		placeholdersStr := "(" + strings.Join(placeholders, ", ") + ")"
		placeholdersArr = append(placeholdersArr, placeholdersStr)
		// add real variables for the replacement of placeholders' '?' letter later.
		mainScope.SQLVars = append(mainScope.SQLVars, scope.SQLVars...)
	}

	mainScope.Raw(fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		mainScope.QuotedTableName(),
		strings.Join(quoted, ", "),
		strings.Join(placeholdersArr, ", "),
	))

	if _, err := mainScope.SQLDB().Exec(mainScope.SQL, mainScope.SQLVars...); err != nil {
		return err
	}
	return nil
}
