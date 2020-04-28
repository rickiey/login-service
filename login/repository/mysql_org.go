// Copyright (c) 2019. icemsoft.net
// Author: Bruce Created:2019/1/2

package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"login-service/login"
	. "login-service/models"
	"time"
)

type mysqlUserRepository struct {
	Conn *gorm.DB
}

//  mysqlUserRepository 结构体的指针实现了接口 Repository
// NewMysqlOrgRepository will create an object that represent the organization.Repository interface
func NewMysqlUserRepository(Conn *gorm.DB) login.Repository {
	return &mysqlUserRepository{Conn}
}

func (r *mysqlUserRepository) UpdatePassword(newPass string, userId int64) error {

	err := r.Conn.Debug().Exec(`UPDATE users  
		SET updated_at = now(), password_digest=?,login_action=(CASE login_action WHEN 1 THEN 9 ELSE login_action END)
	    WHERE id=?`, newPass, userId).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *mysqlUserRepository) FindByUserId(userId int64) (*UserInfo, error) {
	u := UserInfo{}
	err := r.Conn.Debug().Model(u.TableName()).Find(&u, "deleted=0 AND id=?", userId).Error
	if err != nil {
		return &u, err
	}
	return &u, nil
}

func (r *mysqlUserRepository) FindByPhone(phone string) (*UserInfo, error) {
	u := UserInfo{}
	err := r.Conn.Debug().Model(u.TableName()).Find(&u, "deleted=0 AND (email=? or phone=? )", phone, phone).Error
	if err != nil {
		return &u, err
	}
	return &u, nil
}

func (r *mysqlUserRepository) UpdateUserLastLogin(userid int) error {
	s := `update users set last_login = ? where id = ? `
	err := r.Conn.Debug().Exec(s, time.Now(), userid).Error
	if err != nil {
		logrus.Error("update users set last_login  failed:", err)
	}
	return err
}
