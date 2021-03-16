package model

import (
	"gorm.io/gorm"
	"regexp"
	"unicode"
)

const (
	// 密码验证选项 只能含有
	PswOptNumber  uint16 = 1 << iota // 数字 0001
	PswOptLower                      // 小写 0010
	PswOptUpper                      // 大写 0100
	PswOptSpecial                    // 特殊符号 1000
	// alert
	PASSWORDISEMPTY                = "密码为空"
	PASSWORDTOLONG                 = "密码长度过长"
	PASSWORDINCORRECTSPECIFICATION = "密码可为16为，含数字大小写特殊符号"
	USEREMPTY                      = "用户名为空"
	USERTOLONG                     = "用户名过长"
	EMAILEMPTY                     = "邮箱为空"
	EMAILNONSTANDARD               = "邮箱不符合规范"
	USERINCORRECTSPECIFICATION     = "密码可为12为，含数字大小写特殊符号"
)

type VerificationSpecification interface {
	VerifyUserOrEmail(isEmail bool, option, mustOption uint16) (bool, string) // 验证用户名or邮箱是否符合规范
	VerifyPassword(option, mustOption uint16) (bool, string)                  // 验证密码是否符合规范
}
type Register struct {
	Username string `gorm:""` // 用户名
	Password string `gorm:""` // 密码
	Email    string `gorm:""` // 邮箱
	Birthday string `gorm:""` // 生日
	Gender   int64  `gorm:""` // 性别
	gorm.Model
}

// VerifyEmailFormat 验证电子邮箱
func VerifyEmailFormat(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// GeneralValidation 通用验证
// verifyField 验证字段
// option 验证选项
// mustOption 必选项
func GeneralValidation(verifyField string, option, mustOption uint16) bool {
	var verifyResult uint16
	for _, r := range verifyField {
		switch {
		case unicode.IsNumber(r):
			verifyResult = verifyResult | PswOptNumber
		case unicode.IsLower(r):
			verifyResult = verifyResult | PswOptLower
		case unicode.IsUpper(r):
			verifyResult = verifyResult | PswOptUpper
		case unicode.IsPunct(r) || unicode.IsSymbol(r): //标点符号 和 字符
			verifyResult = verifyResult | PswOptSpecial
		default:
			return false
		}
	}
	// 比较option和验证结果
	if option&verifyResult != verifyResult {
		// option和验证不符
		return false
	}
	// 与必选项是否相符合
	return (verifyResult & mustOption) == mustOption
}

// VerifyPassword 验证密码
// option 验证选项 必填
// mustOption 验证必填项 选填
func (user *Register) VerifyPassword(option, mustOption uint16) (bool, string) {
	if user.Password == "" {
		return false, PASSWORDISEMPTY
	}
	if len(user.Password) > 16 {
		return false, PASSWORDTOLONG
	}
	if !GeneralValidation(user.Password, option, mustOption) {
		return false, PASSWORDINCORRECTSPECIFICATION
	}
	return true, ""
}
// VerifyUserOrEmail  验证用户名or邮箱
// isEmail 验证字段是否为Email
// option 验证选项
// mustOption 验证必选项
func (user *Register) VerifyUserOrEmail(isEmail bool, option, mustOption uint16) (bool, string) {
	if user.Username == "" && isEmail == false {
		return false, USEREMPTY
	}
	if user.Email == "" && isEmail == true {
		return false, EMAILEMPTY
	}
	if len(user.Username) > 12 && isEmail == false {
		return false, USERTOLONG
	}
	if isEmail && !VerifyEmailFormat(user.Email) {
		return false, EMAILNONSTANDARD
	}
	if !isEmail && !GeneralValidation(user.Username, option, mustOption) {
		return false, USERINCORRECTSPECIFICATION
	}
	if isEmail {
		return true, ""
	}
	return true, ""
}
