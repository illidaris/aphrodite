package convert

import (
	"time"

	"github.com/spf13/cast"
)

// IsAdult 判断身份证号是否是成年人
func IsAdult(idCard string) bool {
	age := GetAgeWithIDCard18(idCard, ToCST(time.Now()))
	return age >= 18
}

// GetAgeWithIDCard18 判断身份证相对于参考时间多少年
func GetAgeWithIDCard18(idCard string, refDate time.Time) int {
	if len(idCard) != 18 {
		return 0
	}
	var (
		refYear     = refDate.Year()
		refMonth    = int(refDate.Month())
		refDay      = refDate.Day()
		idCardYear  = cast.ToInt(idCard[6:10])
		idcardMonth = cast.ToInt(idCard[10:12]) // 月
		idcardDay   = cast.ToInt(idCard[12:14]) // 日
	)
	baseAge := refYear - idCardYear
	if refMonth < idcardMonth || (refMonth == idcardMonth && refDay < idcardDay) {
		baseAge--
	}
	if baseAge <= 0 {
		return 0
	}
	return baseAge
}

// GetSexWithIDCard18 判断身份证获取性别;-1未知，1性别为男，0则为女
func GetSexWithIDCard18(idCard string) int {
	if len(idCard) != 18 {
		return -1
	}
	var (
		idCardSex = cast.ToInt(idCard[16:17])
	)
	return idCardSex % 2
}
