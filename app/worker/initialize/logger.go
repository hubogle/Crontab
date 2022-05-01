package initialize

import "go.uber.org/zap"

func InitLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger) // 设置一个全局的 logger
}
