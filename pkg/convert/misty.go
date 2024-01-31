package convert

import "strings"

// Misty 模糊化数据，脱敏
// Misty函数接收四个参数：raw，b，e和k，返回一个字符串。
func Misty(raw string, b, e int, k string) string {
	var (
		prefix, suffix string
	)
	// 将raw转换为rune类型的切片
	source := []rune(raw)
	// 获取source的长度
	l := len(source)
	switch {
	case l == 0:
		// 如果source为空，则直接返回raw
		return raw
	case b > e:
		// 如果b大于e，则直接返回raw
		return raw
	case l < b:
		// 如果source的长度小于b，则直接返回raw
		return raw
	}
	// 将source的前b个字符赋值给prefix
	prefix = string(source[:b])
	if l > e {
		// 如果source的长度大于e，则将source的后e-l个字符赋值给suffix
		suffix = string(source[e:])
	} else {
		// 如果source的长度小于等于e，则将e设为source的长度
		e = l
	}
	// 返回prefix、重复k字符(e-b)次、suffix拼接而成的字符串
	return prefix + strings.Repeat(k, e-b) + suffix
}

// MistyDefault将给定的字符串进行处理，并返回处理后的字符串。
// 参数raw为原始字符串。
// 返回值为处理后的字符串。
func MistyDefault(raw string) string {
	// 获取原始字符串的长度
	l := len([]rune(raw))
	// 计算需要替换的起始位置
	b := l / 4
	// 如果原始字符串长度不能被4整除，则需要额外增加一个替换位置
	if l%4 > 0 {
		b++
	}
	// 计算需要替换的结束位置
	e := b + l/2
	// 调用Misty函数进行字符串处理，传入原始字符串、起始位置、结束位置和替换字符
	return Misty(raw, b, e, "*")
}

func MistyMobile(raw string) string {
	if len(raw) != 11 {
		return MistyDefault(raw)
	}
	return Misty(raw, 3, 7, "*")
}

func MistyMail(raw string) string {
	keys := strings.Split(raw, "@")
	if len(keys) < 2 {
		return MistyDefault(raw)
	}
	keys[0] = MistyDefault(keys[0])
	return strings.Join(keys, "@")
}
