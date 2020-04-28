// Copyright (c) 2019. icemsoft.net
// Author: Bruce Created:2019/1/31

package models

import "login-service/utils"

type Users struct {
	Id             int    `json:"id";gorm:"column:id"`
	Email          string `json:"email";gorm:"column:email"`
	Phone          string `json:"phone";gorm:"column:phone"`
	Name           string `json:"name";gorm:"column:name"`
	PasswordDigest string `json:"password_digest";gorm:"column:password_digest"`

	Enable  bool `json:"enable";gorm:"column:enable"`
	Deleted bool `json:"deleted";gorm:"column:deleted"`

	Description string     `json:"description";gorm:"column:description"`
	CreatedAt   utils.Time `json:"created_at";gorm:"column:created_at"`
	UpdatedAt   utils.Time `json:"update_at";gorm:"column:update_at"`

	Remark      string `json:"remark";gorm:"column:remark"`
	LoginAction int    `json:"login_action";gorm:"column:login_action"`
}

func (Users) TableName() string {
	return "users"
}
