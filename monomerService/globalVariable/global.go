package globalVariable

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 全局变量

const (
	TCP  = "127.0.0.1"
	HTTP = "http://127.0.0.1"
	// 注册
	//Error reading registration information
	ERRORREADINGREGISTERTIONINFOMATION = "读取注册信息出错"
)

var (
	GinService *gin.Engine
	Db         *gorm.DB
)
