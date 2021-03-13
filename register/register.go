package main

import (
	"MicroserviceGoBlog/register/model"
	"MicroserviceGoBlog/register/protocol"
	"context"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net"
)

var (
	DB     *gorm.DB
	result *model.Register
)

type RegisterServiceImpl struct {
}

// IsRegistered 判断是否注册user
func (registerInfo *RegisterServiceImpl) IsRegistered(ctx context.Context, req *protocol.IsRegisterRequest) (*protocol.RegisterResponse, error) {
	var response *protocol.RegisterResponse
	DB.Find(&result, "username = ?", req.Username)
	if result.Username == "" {
		response = new(protocol.RegisterResponse)
		response.Status = false
	}
	return response, nil
}

func (registerInfo *RegisterServiceImpl) ReadyRegister(ctx context.Context, req *protocol.RegisterRequest) (*protocol.RegisterResponse, error) {
	var response *protocol.RegisterResponse
	DB.Table("register").Create(&req)
	response = new(protocol.RegisterResponse)
	response.Status = true
	return response, nil
}

func main() {
	dbConf := "root:Qsj.0228@tcp(127.0.0.1:3306)/Blog?charset=utf8&parseTime=True&loc=Local"

	DB, _ = gorm.Open(mysql.Open(dbConf), &gorm.Config{})
	DB.AutoMigrate(&model.Register{})
	service := grpc.NewServer()
	protocol.RegisterRegisterServiceServer(service, new(RegisterServiceImpl))
	monitor, _ := net.Listen("tcp", ":8081")
	service.Serve(monitor)
}
