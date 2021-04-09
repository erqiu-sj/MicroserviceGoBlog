package main

import (
	loginStruct "MicroserviceGoBlog/login/model"
	"MicroserviceGoBlog/login/protocol"
	"MicroserviceGoBlog/monomerService/globalVariable"
	"MicroserviceGoBlog/register/model"
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"net/http"
)

const (
	PORT = 8083
)

var (
	loginServiceHandle *grpc.Server
)

type LoginImpl struct {
}

func (that *LoginImpl) ReadyLogin(_ context.Context, req *protocol.LoginRequest) (res *protocol.LoginResponse, err error) {
	response := &protocol.LoginResponse{}
	var registerResult *model.Register
	registerResult = new(model.Register)
	globalVariable.Db.Table("registers").Find(&registerResult, "username = ?,email = ?", req.Username, req.Email)
	if registerResult.Username == "" && registerResult.Email == "" {
		// 没有注册的情况
		response.Status = false
		response.HttpCode = http.StatusNoContent
		response.Message = loginStruct.REGISTERFIRST
		return response, nil
	}
	globalVariable.Db.Table("logins").Create(&req)
	response.Status = true
	response.Message = loginStruct.LOGINSUCCESS
	response.HttpCode = http.StatusOK
	return response, nil

}
func main() {
	sql, _ := globalVariable.DbInit().DB()
	defer sql.Close()
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

	globalVariable.MicroserviceInit(func(handle *grpc.Server) {
		loginServiceHandle = handle
		protocol.RegisterLoginServiceServer(loginServiceHandle, new(LoginImpl))
	}, func() {
		fmt.Println("login service start!")
	}, func(err string) {
		panic(err)
	}, 8083)

}
