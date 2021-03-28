package main

import (
	"MicroserviceGoBlog/monomerService/globalVariable"
	"MicroserviceGoBlog/monomerService/router"
	"MicroserviceGoBlog/monomerService/utils"
	registerModel "MicroserviceGoBlog/register/model"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
)

var (
	sqlCon *sql.DB
)

func init() {
	globalVariable.GinService = gin.Default()
	globalVariable.GinService.Use(utils.Cors())
	sqlCon, _ = globalVariable.DbInit().DB()
	globalVariable.Db.AutoMigrate(&registerModel.Register{})
	sqlCon.SetMaxIdleConns(10)
	sqlCon.SetMaxOpenConns(100)
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	defer sqlCon.Close()
	globalVariable.GinService.POST("/register", router.Register)
	fmt.Println("gin service start!")
	_ = globalVariable.GinService.Run(":3000")
}
