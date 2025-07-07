package check

import "regexp"

func IsValidEmail(email string) bool {
	// 邮箱正则表达式
	pattern := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	// 编译正则表达式
	reg := regexp.MustCompile(pattern)
	// 使用正则表达式匹配字符串
	return reg.MatchString(email)
}
