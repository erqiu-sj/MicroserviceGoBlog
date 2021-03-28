package globalVariable

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局变量

const (
	TCP  = "127.0.0.1"
	HTTP = "http://127.0.0.1"
	// 注册
	//Error reading registration information
	ERRORREADINGREGISTERTIONINFOMATION = "读取注册信息出错"
	DbConf                             = "root:Qsj.0228@tcp(127.0.0.1:3306)/Blog?charset=utf8&parseTime=True&loc=Local"
)

var (
	GinService *gin.Engine
	Db         *gorm.DB
)

// DbInit 初始化数据库
func DbInit() *gorm.DB {
	Db, _ = gorm.Open(mysql.Open(DbConf), &gorm.Config{})
	return Db
}
