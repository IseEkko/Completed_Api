package global

import (
	config2 "Completed_Api/user-web/config"
	"Completed_Api/user-web/proto"
	ut "github.com/go-playground/universal-translator"
)

var (
	//翻译器
	Trans ut.Translator
	//配置文件
	ServConfig *config2.ServerConfig = &config2.ServerConfig{}
	//客户端链接对象
	UserSrvClient proto.UserClient
)
