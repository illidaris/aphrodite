package check

import "regexp"

// "github.com/nyaruka/phonenumbers" 国外手机号验证

// isValidChineseMobile 验证中国大陆的手机号格式
func IsValidChineseMobile(mobile string) bool {
	// 正则表达式匹配中国大陆的手机格式
	// 1[3-9]\d{9} 表示以1开头，第二位是3-9之间的任意数字，后面跟上9位数字
	regex := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(regex, mobile)
	return matched
}
