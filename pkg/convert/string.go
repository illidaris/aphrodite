package convert

import (
	"encoding/json"
	"unicode/utf8"
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

// TruncateMySQLVarchar 按 MySQL utf8mb4 存储字节长度截断字符串。
func TruncateMySQLVarchar(raw string, length int) string {
	if length <= 0 || raw == "" {
		return ""
	}

	used := 0
	for i, r := range raw {
		size := utf8.RuneLen(r)
		if used+size > length {
			return raw[:i]
		}
		used += size
	}

	return raw
}
