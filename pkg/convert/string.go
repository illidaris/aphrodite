package convert

import (
	"encoding/json"
)

// Json json marshal
// Json将给定的数据结构转换为JSON格式的字符串
func Json(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(b)
}

// TruncateMySQLVarchar 存储字节长度截断字符串。
func TruncateMySQLVarchar(raw string, length int) string {
	if length <= 0 || raw == "" {
		return ""
	}
	runes := []rune(raw)
	if len(runes) <= length {
		return raw
	}
	return string(runes[:length])
}
