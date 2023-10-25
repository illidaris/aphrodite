package convert

import "encoding/json"

// Json json marshal
func Json(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(b)
}
