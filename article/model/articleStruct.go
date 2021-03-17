package model

import "gorm.io/gorm"

type Article struct {
	Title        string `gorm:"TYPE:VARCHAR(50);UNIQUE;NOT NULL"`  // 标题
	Introduction string `gorm:"TYPE:VARCHAR(255);UNIQUE;NOT NULL"` // 简介
	Content      string `gorm:"TYPE:TEXT;NOT NULL;"`               // 文章内容
	Start        uint   `gorm:"TYPE:INT;DEFAULT:0"`                // 点赞
	Tag          string `gorm:"TYPE:VARCHAR(225);NOT NULL"`        // 文章类型
	Browse       uint   `gorm:"TYPE:INT;DEFAULT:0;"`               // 浏览数
	Top          uint   `gorm:"TYPE:INT;DEFAULT:0"`                // 是否置顶 奇数=false 偶数=true
	gorm.Model
}
type ArticleVerification interface {
	VerificationTitle() bool        // 标题是否标准
	VerificationIntroduction() bool // 简介是否标准
	VerificationContent() bool      // 内容是否标准
	VerificationTag() bool          // 文章类型是否标准
}

func (article *Article) VerificationTitle() bool {
	if len(article.Title) > 50 {
		return false
	}
	return true
}

func (article *Article) VerificationIntroduction() bool {
	if len(article.Introduction) > 100 {
		return false
	}
	return true
}

func (article *Article) VerificationContent() bool {
	if len(article.Content) <= 50 {
		return false
	}
	return true
}

func (article *Article) VerificationTag() bool {
	if len(article.Tag) < 2 {
		return false
	}
	return true
}
