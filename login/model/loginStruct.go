package model

import (
	"gorm.io/gorm"
	"regexp"
	"time"
)

type Login struct {
	Username  string    `gorm:"TYPE:VARCHAR(16);NOT NULL;UNIQUE"` // 用户名  or email
	Password  string    `gorm:"TYPE:VARCHAR(100);NOT NULL;"`      // 密码
	LastLogin time.Time // 最后一次登陆时间
	gorm.Model
}
type LoginVerification interface {
	VerifyEmailFormat() bool
}

// VerifyEmailFormat 验证是否为邮箱
func (that *Login) VerifyEmailFormat() bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(that.Username)
}
