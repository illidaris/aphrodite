package sender

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type CodeGenateHandle func(int) string                        // 函数：生成指定的随机码
type ArgsFmtHandle func(any) string                           // 函数：Args格式化
type NotifyMsgHandle func(context.Context, any, string) error // 函数：发送消息

func RandVerifyCode(num int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var captcha string
	// 循环生成指定位数的数字字符
	for i := 0; i < num; i++ {
		randomNumber := r.Intn(10)
		captcha += fmt.Sprintf("%d", randomNumber)
	}
	return captcha
}

func DailyOffset(now time.Time, min time.Duration) time.Duration {
	// 构造次日零点时间对象
	y := now.Year()
	m := now.Month()
	d := now.Day() + 1
	limit := time.Date(y, m, d, 0, 0, 0, 0, now.Location())

	// 计算时间差并进行最小值保护
	diff := limit.Sub(now)
	if diff < min {
		return min
	}
	return diff
}
