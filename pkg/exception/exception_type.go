package exception

type ExceptionType int32

// 通用错误大类
const (
	ERR_COMMON     ExceptionType = (iota + 1) * 10000 // 通用错误码
	ERR_VERIFYCODE                                    // 验证码错误
	ERR_BUSI                                          // 通用业务操作错误
)

// 通用错误码
const (
	ERR_COMMON_BUSY         = ERR_COMMON + (iota + 1) // 服务器繁忙
	ERR_COMMON_REQ_TOOMANEY                           // 请求过多
	ERR_COMMON_BADPARAM                               // 参数错误
	ERR_COMMON_UNAUTH                                 // 登录失效
	ERR_COMMON_NOPERMISSION                           // 没有权限
	ERR_COMMON_USER                                   // 操作人错误
	ERR_COMMON_TIMEOUT                                // 超时
)

// 验证码错误
const (
	ERR_VERIFYCODE_SENDFAIL   = ERR_VERIFYCODE + (iota + 1) // 发送验证码失败
	ERR_VERIFYCODE_HASSEND                                  // 验证码已经发送过
	ERR_VERIFYCODE_HASEXPIRED                               // 验证码已经过期
)

// 业务错误
const (
	ERR_BUSI_NOFOUND  = ERR_BUSI + (iota + 1) // 没有找到
	ERR_BUSI_HASEXIST                         // 已经存在
	ERR_BUSI_CREATE                           // 创建失败
	ERR_BUSI_UPDATE                           // 更新失败
	ERR_BUSI_DELETE                           // 删除失败
	ERR_BUSI_STATUS                           // 状态错误
)

func (ex ExceptionType) New(s string) Exception {
	return New(ex, s)
}

func (ex ExceptionType) Wrap(err error, msgs ...string) Exception {
	return Wrap(ex, err, msgs...)
}
