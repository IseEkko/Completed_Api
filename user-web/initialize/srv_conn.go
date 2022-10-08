package initialize

import (
	"Completed_Api/user-web/global"
	"Completed_Api/user-web/proto"
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func InitSrvConn() {
	/**
	从注册中心获取相关信息
	*/
	cfg := api.DefaultConfig()
	consulInfo := global.ServConfig.CosulConfiginfo
	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)

	userSrvHost := ""
	userSrvPort := 0
	clients, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	data, err := clients.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServConfig.UserSrv.Name))
	if err != nil {
		panic(err)
	}
	for _, vaule := range data {
		userSrvHost = vaule.Address
		userSrvPort = vaule.Port
		break
	}
	/**
	增强业务逻辑
	*/
	if userSrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		return
	}
	zap.S().Infof("[InitSrvConn] ip:%s port:%d", userSrvHost, userSrvPort)
	/**
	获取用户链接
	*/
	dial, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	/**
	这里我们需要考虑的问题：
	1. 服务下线了
	2. 改变端口
	3. 改变ip
	*/
	/**
	这么做的好处就是不用重复的tcp三次握手
	但是使用的时候我们需要注意的是，多个协程公用一个这样会影响一定的性能
	解决的办法：负载均衡、连接池
	https://github.com/processout/grpc-go-pool/blob/master/pool.go
	👆链接池开源代码
	*/
	client := proto.NewUserClient(dial)
	global.UserSrvClient = client
}
