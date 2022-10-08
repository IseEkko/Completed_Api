package main

import (
	"Completed_Api/user-web/initialize"
	"fmt"
	"go.uber.org/zap"
)

var (
	httpPort = 8000
)

func main() {
	//初始化logger
	initialize.InitLogger()
	//初始化配置
	initialize.Config()
	//初始化grpc客户端
	//initialize.GrpcClient()
	//初始化routers
	Router := initialize.Init_Router()
	//初始化翻译器
	if err := initialize.InitTrans("zh"); err != nil {
		zap.S().Infof("翻译器初始化失败")
		panic(err)
	}
	//初始化用户服务连接
	initialize.InitSrvConn()
	if err := Router.Run(fmt.Sprintf(":%d", httpPort)); err != nil {
		zap.S().Panicf("启动服务失败，%s", err.Error())
	}
}
