package main

import (
	"MicroserviceGoBlog/monomerService/globalVariable"
	"fmt"
	"github.com/hashicorp/consul/api"
)

const (
	PORT = 8083
)

func main() {
	loginServiceConsulConf := api.DefaultConfig()
	loginClient, _ := api.NewClient(loginServiceConsulConf)
	loginReg := api.AgentServiceRegistration{
		ID:      "loginMain",
		Name:    "login",
		Tags:    []string{"login", "verifyLogin"},
		Address: globalVariable.TCP,
		Port:    PORT,
		Check: &api.AgentServiceCheck{
			CheckID:  "checkLogin",
			Name:     "loginUser",
			TCP:      fmt.Sprint(globalVariable.TCP, ":", PORT),
			Interval: "5s",
			Timeout:  "5s",
		},
	}
	_ = loginClient.Agent().ServiceRegister(&loginReg)

	fmt.Println("login service start!")

}
