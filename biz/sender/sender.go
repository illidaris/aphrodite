package sender

import (
	"context"

	"github.com/illidaris/aphrodite/pkg/exception"
)

/*
	eg:
		"{APP应用名}:{BIZID}:{BUSI业务名}:{CODE功能类别}:{ID唯一值}"

		"_aphrodite:5078:login_code:code:19922028858"
		"_aphrodite:5078:login_code:locked:19922028858"
		"_aphrodite:5078:login_code:limiters_ip:192.0.0.1"
		"_aphrodite:5078:login_code:limiters_uid:11234"
*/

var (
	globalApp = "_aphrodite"
)

func SetApp(v string) {
	globalApp = v
}

type CodeRequest struct{}

func SendFunc(opts ...SenderOption) func(context.Context, any) (int64, exception.Exception) {
	opt := NewSenderOptions(opts...)
	return func(ctx context.Context, id any) (int64, exception.Exception) {
		// 检查是否可以发送验证码
		ttl, ex := opt.CheckExistLocked(ctx, id)
		if ex != nil {
			return ttl, ex
		}
		if ttl > 0 {
			return ttl, nil
		}
		// 其他限制器校验

		//

		return int64(lockDur.Seconds()), false, nil
	}
}

func VerifyFunc(opts ...SenderOption) func(context.Context, any) (int64, exception.Exception) {
	return func(ctx context.Context, id any) (int64, exception.Exception) {
		return 0, nil
	}
}
