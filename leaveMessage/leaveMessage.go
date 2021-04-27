package main

import (
	"MicroserviceGoBlog/monomerService/globalVariable"
	"fmt"
	"github.com/hashicorp/consul/api"
	"time"
)

const (
	PORT = 8090
)

// 留言

func main() {
	sql, _ := globalVariable.DbInit().DB()
	sql.SetMaxOpenConns(20)
	sql.SetMaxIdleConns(3)
	sql.SetConnMaxLifetime(5 * time.Minute)
	defer sql.Close()

	// 往consul中注册服务
	leaveMessageServiceConsulConf := api.DefaultConfig()
	leaveMessageClient, _ := api.NewClient(leaveMessageServiceConsulConf)
	leaveMessageClientReg := api.AgentServiceRegistration{
		ID:      "leaveMessageClientRegMain",
		Name:    "leaveMessageClientReg",
		Tags:    []string{"leaveMessageClientReg"},
		Address: globalVariable.TCP,
		Port:    PORT,
		Check: &api.AgentServiceCheck{
			CheckID: "checkLeaveMessage",
			Name:    "authLeaveMessage",
			//HTTP:    fmt.Sprint(globalVariable.HTTP, ":", PORT),
			TCP:      fmt.Sprint(globalVariable.TCP, ":", PORT),
			Interval: "10s",
			Timeout:  "5s",
		},
	}
	_ = leaveMessageClient.Agent().ServiceRegister(&leaveMessageClientReg)
}
