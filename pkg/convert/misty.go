package convert

import "strings"

// Misty 模糊化数据，脱敏
func Misty(raw string, b, e int, k string) string {
	var (
		prefix, suffix string
	)
	source := []rune(raw)
	l := len(source)
	switch {
	case l == 0:
		return raw
	case b > e:
		return raw
	case l < b:
		return raw
	}
	prefix = string(source[:b])
	if l > e {
		suffix = string(source[e:])
	} else {
		e = l
	}
	return prefix + strings.Repeat(k, e-b) + suffix
}

func MistyDefault(raw string) string {
	l := len([]rune(raw))
	b := l / 4
	if l%4 > 0 {
		b++
	}
	e := b + l/2
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
