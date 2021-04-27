package main

import (
	"MicroserviceGoBlog/monomerService/globalVariable"
	"fmt"
	"github.com/hashicorp/consul/api"
	"time"
)

const (
	PORT = 8089
)

//  对文章进行权限控制
func main() {
	sql, _ := globalVariable.DbInit().DB()
	sql.SetMaxOpenConns(20)
	sql.SetMaxIdleConns(3)
	sql.SetConnMaxLifetime(5 * time.Minute)
	defer sql.Close()

	// 往consul中注册服务
	authArticleServiceConsulConf := api.DefaultConfig()
	authArticleClient, _ := api.NewClient(authArticleServiceConsulConf)
	authArticleReg := api.AgentServiceRegistration{
		ID:      "authArticleMain",
		Name:    "authArticle",
		Tags:    []string{"authArticle", "verifyArticle"},
		Address: globalVariable.TCP,
		Port:    PORT,
		Check: &api.AgentServiceCheck{
			CheckID: "checkArticle",
			Name:    "authArticleAdmin",
			//HTTP:    fmt.Sprint(globalVariable.HTTP, ":", PORT),
			TCP:      fmt.Sprint(globalVariable.TCP, ":", PORT),
			Interval: "10s",
			Timeout:  "5s",
		},
	}
	_ = authArticleClient.Agent().ServiceRegister(&authArticleReg)
}
