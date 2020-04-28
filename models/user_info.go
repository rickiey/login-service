// Copyright (c) 2019. icemsoft.net
// Author: Bruce Created:2019/1/31

package models

import (
	"errors"
	"github.com/fatih/structs"
)

const NeedChangePasswordFirstLogin = 1

var (
	ErrNeedChangePasswordFirstLogin = errors.New("第一次登录需要修改密码")
)

type UserInfo struct {
	ID             int    `json:"id" structs:"id" `
	Email          string `json:"email" structs:"email"`
	Name           string `json:"name" structs:"name"`
	Phone          string `json:"phone" structs:"phone"`
	PasswordDigest string `json:"password_digest" structs:"password_digest"`
	Enable         bool   `json:"enable" structs:"enable"`
	Deleted        bool   `json:"deleted" structs:"deleted"`
}

func (UserInfo) TableName() string {
	return "users"
}

//for ruby bool string
func (u UserInfo) ConvertSessionFields() map[string]interface{} {
	userFields := structs.Map(u)
	if u.Enable {
		userFields["enable"] = "true"
	}
	if u.Deleted {
		userFields["deleted"] = "true"
	}
	return userFields
}

//
func (u UserInfo) RemoveSensitive() map[string]interface{} {
	userFields := structs.Map(u)
	delete(userFields, "password_digest")
	delete(userFields, "deleted")
	return userFields
}
