package convert

import (
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

// TestMisty
func TestMisty(t *testing.T) {
	convey.Convey("test Misty", t, func() {
		convey.Convey("source is empty", func() {
			source := ""
			dStr := Misty(source, 1, 3, "*")
			convey.ShouldEqual(dStr, source)
		})
		convey.Convey("source is \"310229198805193214\"", func() {
			source := "310229198805193214"
			resultMap := map[string][]interface{}{
				"****29198805193214": {0, 4, "*"},
				"3102**********3214": {4, 14, "*"},
				"31022919880519****": {14, 26, "*"},
				"310229198805193214": {20, 26, "*"},
			}
			for res, params := range resultMap {
				b := params[0].(int)
				e := params[1].(int)
				k := params[2].(string)
				title := fmt.Sprintf("[%d,%d] is %s, %s", b, e, k, res)
				convey.Convey(title, func() {
					dStr := Misty(source, b, e, k)
					convey.So(dStr, convey.ShouldEqual, res)
				})
			}
		})
	})
}

func TestMistyDefault(t *testing.T) {
	convey.Convey("test MistyDefault", t, func() {
		convey.Convey("source is empty \"\"", func() {
			source := ""
			dStr := MistyDefault(source)
			convey.ShouldEqual(dStr, source)
		})
		convey.Convey("source is no empty", func() {
			resultMap := map[string]string{
				"大苏打撒旦":   "大苏**旦",
				"AAAAVVV": "AA***VV",
				"bbbb":    "b**b",
				"a":       "a",
			}
			for s, t := range resultMap {
				title := fmt.Sprintf("%s ===> %s", s, t)
				convey.Convey(title, func() {
					dStr := MistyDefault(s)
					println(s, "=>", dStr)
					convey.So(dStr, convey.ShouldEqual, t)
				})
			}
		})
	})
}

func TestMistyMobile(t *testing.T) {
	convey.Convey("test MistyMobile", t, func() {
		convey.Convey("source is empty \"\"", func() {
			source := ""
			dStr := MistyMobile(source)
			convey.ShouldEqual(dStr, source)
		})
		convey.Convey("source is no empty", func() {
			resultMap := map[string]string{
				"13912345678": "139****5678",
				"19922334455": "199****4455",
			}
			for s, t := range resultMap {
				title := fmt.Sprintf("%s ===> %s", s, t)
				convey.Convey(title, func() {
					dStr := MistyMobile(s)
					println(s, "=>", dStr)
					convey.So(dStr, convey.ShouldEqual, t)
				})
			}
		})
	})
}

func TestMistyMail(t *testing.T) {
	convey.Convey("test MistyMail", t, func() {
		convey.Convey("source is empty \"\"", func() {
			source := ""
			dStr := MistyMail(source)
			convey.ShouldEqual(dStr, source)
		})
		convey.Convey("source is no empty", func() {
			resultMap := map[string]string{
				"wda&asd@126.com": "wd***sd@126.com",
				"wsd@126.com":     "w*d@126.com",
				"a@126.com":       "a@126.com",
			}
			for s, t := range resultMap {
				title := fmt.Sprintf("%s ===> %s", s, t)
				convey.Convey(title, func() {
					dStr := MistyMail(s)
					println(s, "=>", dStr)
					convey.So(dStr, convey.ShouldEqual, t)
				})
			}
		})
	})
}
