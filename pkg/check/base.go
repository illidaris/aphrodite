package check

import (
	"encoding/json"
	"net"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"
)

var (
	chineseMobileRegexp      = regexp.MustCompile(`^1[3-9]\d{9}$`)
	internationalPhoneRegexp = regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	uuidRegexp               = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
	chineseIDCardRegexp      = regexp.MustCompile(`^\d{17}(\d|X)$`)
	md5Regexp                = regexp.MustCompile(`^[a-fA-F0-9]{32}$`)
	sha1Regexp               = regexp.MustCompile(`^[a-fA-F0-9]{40}$`)
	sha256Regexp             = regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
	sha512Regexp             = regexp.MustCompile(`^[a-fA-F0-9]{128}$`)
)

// IsIPv4 判断字符串是否为规范格式的 IPv4 地址。
func IsIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil && parsedIP.To16() != nil && parsedIP.To4().String() == ip
}

// IsIPv6 判断字符串是否为 IPv6 地址。
func IsIPv6(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() == nil && parsedIP.To16() != nil
}

// IsIP 判断字符串是否为有效的 IP 地址（IPv4 或 IPv6）。
func IsIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// IsPort 判断字符串是否为 0 到 65535 范围内的端口号。
func IsPort(portStr string) bool {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return false
	}
	return port >= 0 && port <= 65535
}

// IsInternationalPhone 判断字符串是否符合国际电话号码格式（+ 后跟 2 到 15 位数字）。
func IsInternationalPhone(phone string) bool {
	return internationalPhoneRegexp.MatchString(phone)
}

// IsGender 判断性别值是否为 male、female 或 unknown。
func IsGender(gender string) bool {
	switch gender {
	case "male", "female", "unknown":
		return true
	default:
		return false
	}
}

// IsStringLengthInRange 判断字符串的 Unicode 字符数是否在闭区间 [min, max] 内。
func IsStringLengthInRange(s string, min, max int) bool {
	length := utf8.RuneCountInString(s)
	return length >= min && length <= max
}

// IsDateTime 判断日期时间字符串是否能按指定布局解析。
func IsDateTime(dtStr, layout string) bool {
	_, err := time.Parse(layout, dtStr)
	return err == nil
}

// IsUUID 判断字符串是否为带连字符的 UUID（版本 1 到 5）。
func IsUUID(uuid string) bool {
	return uuidRegexp.MatchString(uuid)
}

// IsJSON 判断字符串是否为语法有效的 JSON。
func IsJSON(jsonStr string) bool {
	return json.Valid([]byte(jsonStr))
}

// IsPasswordStrong 判断密码是否至少 8 个字符且包含四类字符中的至少三类。
func IsPasswordStrong(password string) bool {
	if utf8.RuneCountInString(password) < 8 {
		return false
	}
	level := 0
	if ok, _ := regexp.MatchString(`[a-z]`, password); ok {
		level++
	}
	if ok, _ := regexp.MatchString(`[A-Z]`, password); ok {
		level++
	}
	if ok, _ := regexp.MatchString(`[0-9]`, password); ok {
		level++
	}
	if ok, _ := regexp.MatchString(`[\W_]`, password); ok { // \W 匹配非单词字符
		level++
	}
	return level >= 3
}

// IsMAC 判断字符串是否为有效的 MAC 地址。
func IsMAC(mac string) bool {
	_, err := net.ParseMAC(mac)
	return err == nil
}

// IsIDCard 判断字符串是否为校验位正确的 18 位中国居民身份证号码。
func IsIDCard(id string) bool {
	if !chineseIDCardRegexp.MatchString(id) {
		return false
	}

	var (
		weight = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
		check  = []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
		sum    = 0
	)

	for i := 0; i < 17; i++ {
		// '0' 的 ASCII 码是 48
		sum += int(id[i]-'0') * weight[i]
	}

	return check[sum%11] == id[17]
}

// IsMD5 判断字符串是否为 32 位十六进制 MD5 摘要。
func IsMD5(s string) bool {
	return md5Regexp.MatchString(s)
}

// IsSHA1 判断字符串是否为 40 位十六进制 SHA-1 摘要。
func IsSHA1(s string) bool {
	return sha1Regexp.MatchString(s)
}

// IsSHA256 判断字符串是否为 64 位十六进制 SHA-256 摘要。
func IsSHA256(s string) bool {
	return sha256Regexp.MatchString(s)
}

// IsSHA512 判断字符串是否为 128 位十六进制 SHA-512 摘要。
func IsSHA512(s string) bool {
	return sha512Regexp.MatchString(s)
}
