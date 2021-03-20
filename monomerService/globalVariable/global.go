package globalVariable

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 全局变量

const (
	TCP  = "127.0.0.1"
	HTTP = "http://127.0.0.1"
)

var (
	GinService *gin.Engine
	Db         *gorm.DB
)
