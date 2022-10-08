package initialize

import (
	"Completed_Api/user-web/middlewares"
	router2 "Completed_Api/user-web/router"
	"github.com/gin-gonic/gin"
)

//TODO 初始化Router代码
func Init_Router() *gin.Engine {
	router := gin.Default()
	/**
	解决跨域问题
	*/
	router.Use(middlewares.Cors())
	GlobalRouter := router.Group("/v1")
	router2.InitUserServer(GlobalRouter)
	//验证码路由
	router2.InitBaseRouter(GlobalRouter)
	return router
}
