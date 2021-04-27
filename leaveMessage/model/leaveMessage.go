package model

import "gorm.io/gorm"

const (
	MESSAGELENGTHISTOOLONG = "留言长度过长"
	FROMADRESS             = "留言地址不能为空"
)

type LeaveMessage struct {
	gorm.Model
	To        string `gorm:"TYPE:VARCHAR(50);"`                  // 回复谁？
	ArticleId uint   `gorm:"TYPE:INT;"`                          // 如果给文章留言该字段不可能为空
	from      string `gorm:"TYPE:VARCHAR(100);NOT NULL;"`        // 在哪留言
	Message   string `gorm:"TYPE:VARCHAR(255);NOT NULL;UNIQUE;"` // 留言信息
	MeInfo    string `gorm:"TYPE:VARCHAR(255);DEFAULT:0;"`       // 留言用户
	Report    uint   `gorm:"TYPE:INT;DEFAULT:0;"`                // 举报次数 前端不予显示
	Start     uint   `gorm:"TYPE:INT;DEFAULT:0;"`                // 点赞次数
}

type LeaveMessageVerification interface {
	VerificationMessage() (string, bool) // 验证留言信息
	VerificationFrom() (string, bool)    // 验证留言地址
}

func (that *LeaveMessage) VerificationMessage() (string, bool) {
	if len(that.Message) > 200 {
		return MESSAGELENGTHISTOOLONG, false
	}
	return "", true
}

func (that *LeaveMessage) VerificationFrom() (string, bool) {
	if len(that.from) == 0 {
		return FROMADRESS, false
	}
	return "", true
}
