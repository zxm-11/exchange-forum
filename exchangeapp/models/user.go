package models

import "gorm.io/gorm"

type User struct {

	gorm.Model //包含id,created_at,updated_at字段
	Username string `gorm:"unique"` //用户名唯一
	Password string

}