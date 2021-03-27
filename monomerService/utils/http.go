package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ReturnData 返回数据
func ReturnData(C *gin.Context, status int, data interface{}) {
	C.JSON(status, gin.H{
		"data": data,
	})
}

// Cors 跨域设置
func Cors() gin.HandlerFunc {
	return func(C *gin.Context) {
		C.Header("Access-Control-Allow-Origin", "*")
		C.Header("Access-Control-Allow-Headers", "content-type,Authorization")
		C.Header("Access-Control-Allow-Methods", "DELETE,PUT,POST,GET,OPTIONS")
		if C.Request.Method == "OPTIONS" {
			C.AbortWithStatus(http.StatusNoContent)
		} else {
			C.Next()
		}
	}
}
