package router

import (
	"Completed_Api/user-web/api"
	"Completed_Api/user-web/middlewares"
	"github.com/gin-gonic/gin"
)

//TODO UserServer
func InitUserServer(Router *gin.RouterGroup) {
	userRouter := Router.Group("user")
	{ //加入登录验证和身份验证机制
		userRouter.GET("GetUserList", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		userRouter.POST("PasswordLogin", api.PasswordLogin)
		userRouter.POST("Register", api.Register)
	}
}
