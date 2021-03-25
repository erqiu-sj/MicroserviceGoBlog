package main

import (
	"MicroserviceGoBlog/monomerService/globalVariable"
	"MicroserviceGoBlog/register/model"
	"MicroserviceGoBlog/register/protocol"
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
)

const (
	PORT = 8081
)

var (
	result                *model.Register
	registerServiceHandle *grpc.Server
)

type RegisterServiceImpl struct {
}

// IsRegistered 判断是否注册user
func (registerInfo *RegisterServiceImpl) IsRegistered(_ context.Context, req *protocol.IsRegisterRequest) (*protocol.RegisterResponse, error) {
	response := &protocol.RegisterResponse{}
	globalVariable.Db.Find(&result, "username = ?", req.Username)
	if result.Username == "" {
		response.Status = false
		response.Message = "用户不存在"
	}
	return response, nil
}

// ReadyRegister 准备注册user
func (registerInfo *RegisterServiceImpl) ReadyRegister(ctx context.Context, req *protocol.RegisterRequest) (*protocol.RegisterResponse, error) {
	msg, _ := metadata.FromIncomingContext(ctx)
	fmt.Println(msg)
	response := &protocol.RegisterResponse{}
	register := model.Register{Username: req.Username, Password: req.Password, Email: req.Email, Birthday: req.Birthday, Gender: req.Gender}
	if register.Email != "" {
		emailState, emailMessage := register.VerifyUserOrEmail(true, 0, 0)
		if !emailState {
			response.Status = emailState
			response.Message = emailMessage
			return response, nil
		}
	}
	usrNameState, usrMessage := register.VerifyUserOrEmail(false, model.PswOptLower&model.PswOptNumber&model.PswOptSpecial&model.PswOptUpper, model.PswOptSpecial)
	if !usrNameState {
		response.Status = usrNameState
		response.Message = usrMessage
		return response, nil
	}
	globalVariable.Db.Table("register").Create(&req)
	response.Status = true
	return response, nil
}

// StopRegister 关闭注册服务
func (registerInfo *RegisterServiceImpl) StopRegister(_ context.Context, req *protocol.StopRequest) (*protocol.StopResponse, error) {
	if req.StopMethod == "stop" {
		// 关闭
		registerServiceHandle.Stop()
	}
	if req.StopMethod == "gracefulStop" {
		// 优雅的关闭
		registerServiceHandle.GracefulStop()
	}
	return &protocol.StopResponse{Message: req.StopMethod, Status: true}, nil
}

func main() {
	// 往consul中注册服务
	registerServiceConsulConf := api.DefaultConfig()
	registerClient, _ := api.NewClient(registerServiceConsulConf)
	registerReg := api.AgentServiceRegistration{
		ID:      "registerMain",
		Name:    "register",
		Tags:    []string{"register", "verifyRegister"},
		Address: globalVariable.TCP,
		Port:    PORT,
		Check: &api.AgentServiceCheck{
			CheckID: "checkRegister",
			Name:    "registerUser",
			//HTTP:    fmt.Sprint(globalVariable.HTTP, ":", PORT),
			TCP:      fmt.Sprint(globalVariable.TCP, ":", PORT),
			Interval: "10s",
			Timeout:  "5s",
		},
	}
	_ = registerClient.Agent().ServiceRegister(&registerReg)

	registerServiceHandle = grpc.NewServer()
	protocol.RegisterRegisterServiceServer(registerServiceHandle, new(RegisterServiceImpl))
	monitor, _ := net.Listen("tcp", ":8082")
	fmt.Println("register service start!")
	registerServiceHandle.Serve(monitor)
}
