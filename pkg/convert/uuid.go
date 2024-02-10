package convert

import "github.com/google/uuid"

// RandomID 生成一个随机的UUID字符串
func RandomID() string {
	// 使用uuid包生成一个随机的UUID字符串
	return uuid.NewString()
}

// Sha1ID 生成一个基于给定数据的SHA1 ID。
// 参数：
//   data：要生成ID的数据
// 返回值：
//   生成的SHA1 ID字符串
func Sha1ID(data []byte) string {
	source, err := uuid.NewDCEPerson()
	if err != nil {
		return ""
	}
	id := uuid.NewSHA1(source, data)
	return id.String()
}

// Sha1IDSimple 生成一个基于SHA1哈希值的UUID字符串
func Sha1IDSimple(data []byte) string {
	// 使用给定的数据生成一个UUID
	id := uuid.NewSHA1(uuid.Nil, data)
	// 返回UUID的字符串表示
	return id.String()
}
