package main

import (
	"MicroserviceGoBlog/monomerService/globalVariable"
	registerModel "MicroserviceGoBlog/register/model"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func init() {
	dbConf := "root:Qsj.0228@tcp(127.0.0.1:3306)/Blog?charset=utf8&parseTime=True&loc=Local"
	globalVariable.GinService = gin.Default()
	globalVariable.Db, _ = gorm.Open(mysql.Open(dbConf), &gorm.Config{})
	globalVariable.Db.AutoMigrate(&registerModel.Register{})
	globalVariable.GinService.Run(":3000")
}

func main() {
	Server := gin.Default()
	Server.POST("/register")
	if err := Server.Run(":3000"); err != nil {
		log.Fatal("服务器启动失败")
	}

}
