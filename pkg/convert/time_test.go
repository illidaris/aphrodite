package convert

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

// TestToCST
func TestToCST(t *testing.T) {
	convey.Convey("TestToCST", t, func() {
		convey.Convey("empty", func() {
			now := time.Now()
			println(now.Format(DefaultTimeFormat))
			utc := now.UTC()
			println(utc.Format(DefaultTimeFormat))
			newNow := ToCST(now)
			println(newNow.Format(DefaultTimeFormat))
			println("---")
		})
	})
}

// TestFormUnixToString
func TestFormUnixToString(t *testing.T) {
	convey.Convey("TestFormUnixToString", t, func() {
		ts := int64(1696905231)
		def := "2023-10-10 10:33:51"
		num := "20231010103351"
		convey.Convey("DefaultTimeFormat", func() {
			ft := FormUnixToString(ts, DefaultTimeFormat)
			convey.So(ft, convey.ShouldEqual, def)
		})
		convey.Convey("NumberTimeFormat", func() {
			ft := FormUnixToString(ts, NumberTimeFormat)
			convey.So(ft, convey.ShouldEqual, num)
		})
	})
}

// TestTimeFormat
func TestTimeFormat(t *testing.T) {
	convey.Convey("TestFormUnixToString", t, func() {
		ts := int64(1696905231)
		t := time.Unix(ts, 0)
		def := "2023-10-10 10:33:51"
		num := "20231010103351"
		convey.Convey("DefaultTimeFormat", func() {
			ft := TimeFormat(t)
			convey.So(ft, convey.ShouldEqual, def)
		})
		convey.Convey("NumberTimeFormat", func() {
			ft := TimeNumber(t)
			convey.So(ft, convey.ShouldEqual, num)
		})
	})
}
