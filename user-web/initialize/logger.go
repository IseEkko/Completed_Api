package initialize

import "go.uber.org/zap"

//TODO 初始化Logger代码
func InitLogger() {
	logger, _ := zap.NewProduction()
	zap.ReplaceGlobals(logger)
	zap.S().Info("logger init success")
}
