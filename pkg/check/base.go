package check

import (
	"encoding/json"
	"net"
	"net/url"
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

func IsIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil && parsedIP.To16() != nil && parsedIP.To4().String() == ip
}

func IsIPv6(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() == nil && parsedIP.To16() != nil
}

func IsIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func IsPort(portStr string) bool {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return false
	}
	return port >= 0 && port <= 65535
}

func IsURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

func IsURI(uriStr string) bool {
	_, err := url.ParseRequestURI(uriStr)
	return err == nil
}

func IsInternationalPhone(phone string) bool {
	return internationalPhoneRegexp.MatchString(phone)
}

func IsGender(gender string) bool {
	switch gender {
	case "male", "female", "unknown":
		return true
	default:
		return false
	}
}

func IsStringLengthInRange(s string, min, max int) bool {
	length := utf8.RuneCountInString(s)
	return length >= min && length <= max
}

func IsDateTime(dtStr, layout string) bool {
	_, err := time.Parse(layout, dtStr)
	return err == nil
}

func IsUUID(uuid string) bool {
	return uuidRegexp.MatchString(uuid)
}

func IsJSON(jsonStr string) bool {
	return json.Valid([]byte(jsonStr))
}

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

func IsMAC(mac string) bool {
	_, err := net.ParseMAC(mac)
	return err == nil
}

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

func IsMD5(s string) bool {
	return md5Regexp.MatchString(s)
}

func IsSHA1(s string) bool {
	return sha1Regexp.MatchString(s)
}

func IsSHA256(s string) bool {
	return sha256Regexp.MatchString(s)
}

func IsSHA512(s string) bool {
	return sha512Regexp.MatchString(s)
}
