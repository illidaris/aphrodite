package convert

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

// TestIsAdult
func TestIsAdult(t *testing.T) {
	convey.Convey("TestIsAdult", t, func() {
		convey.Convey("empty", func() {
			id := ""
			isAdult := IsAdult(id)
			convey.So(isAdult, convey.ShouldBeFalse)
		})
		convey.Convey("not adult", func() {
			id := "310229299502233214"
			isAdult := IsAdult(id)
			convey.So(isAdult, convey.ShouldBeFalse)
		})
		convey.Convey("adult", func() {
			id := "310229200502233214"
			isAdult := IsAdult(id)
			convey.So(isAdult, convey.ShouldBeTrue)
		})
	})
}

// TestGetAgeWithIDCard18
func TestGetAgeWithIDCard18(t *testing.T) {
	convey.Convey("TestGetAgeWithIDCard18", t, func() {
		refDate := time.Date(2023, 2, 23, 0, 0, 0, 0, time.Local)
		convey.Convey("empty", func() {
			id := ""
			age := GetAgeWithIDCard18(id, refDate)
			convey.So(age, convey.ShouldAlmostEqual, 0)
		})
		convey.Convey("success 19910222 20230223", func() {
			id := "310229199102223214"
			age := GetAgeWithIDCard18(id, refDate)
			convey.So(age, convey.ShouldAlmostEqual, 32)
		})
		convey.Convey("success 19910224 20230223", func() {
			id := "310229199102243214"
			age := GetAgeWithIDCard18(id, refDate)
			convey.So(age, convey.ShouldAlmostEqual, 31)
		})
		convey.Convey("success 19911124 20230223", func() {
			id := "310229199111243214"
			age := GetAgeWithIDCard18(id, refDate)
			convey.So(age, convey.ShouldAlmostEqual, 31)
		})
	})
}

// TestGetSexWithIDCard18
func TestGetSexWithIDCard18(t *testing.T) {
	convey.Convey("TestGetSexWithIDCard18", t, func() {
		convey.Convey("empty", func() {
			id := ""
			sex := GetSexWithIDCard18(id)
			convey.So(sex, convey.ShouldAlmostEqual, -1)
		})
		convey.Convey("boy", func() {
			id := "310229199111243294"
			sex := GetSexWithIDCard18(id)
			convey.So(sex, convey.ShouldAlmostEqual, 1)
		})
		convey.Convey("girl", func() {
			id := "310229199111243264"
			sex := GetSexWithIDCard18(id)
			convey.So(sex, convey.ShouldAlmostEqual, 0)
		})
	})
}
