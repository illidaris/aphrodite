package convert

import "encoding/json"

// Json json marshal
// Json将给定的数据结构转换为JSON格式的字符串
func Json(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(b)
}
