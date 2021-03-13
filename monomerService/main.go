package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	Server := gin.Default()
	Server.POST("/register",)
	Server.Run(":3000")

}
