// Copyright (c) 2019. icemsoft.net
// Author: Bruce Created:2019/1/2

package login

import (
	. "login-service/models"
)

type Repository interface {
	FindByUserId(userId int64) (*UserInfo, error)
	FindByPhone(phone string) (*UserInfo, error)
	UpdatePassword(newPass string, userId int64) error
	UpdateUserLastLogin(userid int) error
}

type CacheRepository interface {
	SetSession(sid string, fields map[string]interface{}, ttl int64) error
	GetSession(sid string) (map[string]string, error)
	DoubleCheckSession(phone string) error
	Delete(key ...interface{}) error
	Get(key string) (string, error)
}
