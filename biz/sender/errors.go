package sender

import "github.com/illidaris/aphrodite/pkg/exception"

var (
	ExStoreNil  = exception.ERR_VERIFYCODE.New("系统配置库未初始化")
	ExSenderErr = exception.ERR_VERIFYCODE.New("系统配置错误")
)
