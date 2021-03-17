package model

import (
	"gorm.io/gorm"
	"time"
)

type Login struct {
	Username  string    `gorm:""` // 用户名
	Password  string    `gorm:""` // 密码
	Email     string    `gorm:""` // 邮箱
	LastLogin time.Time `gorm:""` // 最后一次登陆时间
	gorm.Model
}
type LoginVerification interface {
}
