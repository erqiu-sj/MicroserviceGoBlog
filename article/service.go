package main

import (
	"MicroserviceGoBlog/article/model"
	"MicroserviceGoBlog/article/protocol"
	"MicroserviceGoBlog/monomerService/globalVariable"
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

var (
	articleServiceHandle *grpc.Server
)

const (
	PORT = 8088
)

type ArticleImpl struct {
}

// CreateArticle 新增文章
func (that *ArticleImpl) CreateArticle(steam protocol.ArticleService_CreateArticleServer) error {
	return nil
}

// ModifyArticles 编辑文章
func (that *ArticleImpl) ModifyArticles(_ context.Context, req *protocol.EditArticleRequest) (*protocol.Response, error) {
	res := new(protocol.Response)
	return res, nil
}

// StrikeOutArticle 删除文章
func (that *ArticleImpl) StrikeOutArticle(_ context.Context, req *protocol.DeleteArticleRequest) (*protocol.Response, error) {
	res := new(protocol.Response)
	result := model.Article{}
	if globalVariable.Db.Table("articles").Delete(&result, req.GetArticleId()).Error != nil {
		// 删除失败
		res.HttpStatus = http.StatusOK
		res.Status = true
	} else {
		// 删除成功
		res.Status = false
		res.HttpStatus = http.StatusOK
	}
	return res, nil
}

// DemandArticle 查询文章
func (that *ArticleImpl) DemandArticle(_ context.Context, req *protocol.QueryArticleRequest) (*protocol.Response, error) {
	res := new(protocol.Response)
	result := model.Article{}
	globalVariable.Db.Table("articles").Find(&result, "id = ?", req.GetArticleId())
	if result.ID != 0 {
		//有数据
		res.HttpStatus = http.StatusOK
		res.Status = true
	} else {
		// 没有数据
		res.Status = false
		res.HttpStatus = http.StatusOK
	}
	return res, nil
}
func main() {
	sql, _ := globalVariable.DbInit().DB()
	sql.SetMaxOpenConns(20)
	sql.SetMaxIdleConns(3)
	sql.SetConnMaxLifetime(5 * time.Minute)
	defer sql.Close()

	articleServiceConsulConf := api.DefaultConfig()
	articleClient, _ := api.NewClient(articleServiceConsulConf)
	articleReg := api.AgentServiceRegistration{
		ID:      "articleMain",
		Name:    "article",
		Tags:    []string{"createArticle", "deleteArticle", "editArticle", "lookingArticle"},
		Address: globalVariable.TCP,
		Port:    PORT,
		Check: &api.AgentServiceCheck{
			CheckID:  "checkArticle",
			Name:     "articleAdmin",
			TCP:      fmt.Sprint(globalVariable.TCP, ":", PORT),
			Interval: "10s",
			Timeout:  "5s",
		},
	}
	_ = articleClient.Agent().ServiceRegister(&articleReg)

	globalVariable.MicroserviceInit(func(handle *grpc.Server) {
		articleServiceHandle = handle
		protocol.RegisterArticleServiceServer(articleServiceHandle, new(ArticleImpl))
	}, func() {
		fmt.Println("article start!")
	}, func(err string) {
		panic(err)
	}, PORT)

}
