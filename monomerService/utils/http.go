package utils

import "github.com/gin-gonic/gin"

// ReturnData 返回数据
func ReturnData(C *gin.Context, status int, data interface{}) {
	C.JSON(status, gin.H{
		"data": data,
	})
}
