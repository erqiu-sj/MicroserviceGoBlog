package main

import (
	"MicroserviceGoBlog/monomerService/globalVariable"
	"MicroserviceGoBlog/monomerService/router"
	"MicroserviceGoBlog/monomerService/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

func init() {
	globalVariable.GinService = gin.Default()
	globalVariable.GinService.Use(utils.Cors())
	globalVariable.DbInit()
	sql, _ := globalVariable.Db.DB()
	sql.SetMaxIdleConns(10)
	sql.SetMaxOpenConns(100)
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	globalVariable.GinService.POST("/register", router.Register)
	fmt.Println("gin service start!")
	fmt.Println(globalVariable.Db, "TYPE!")
	_ = globalVariable.GinService.Run(":3000")
}
