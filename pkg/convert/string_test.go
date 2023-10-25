package convert

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

// TestJson
func TestJson(t *testing.T) {
	convey.Convey("TestJson", t, func() {
		m := map[string]interface{}{"a": "x", "b": 1}
		ft := "{\"a\":\"x\",\"b\":1}"
		convey.Convey("empty", func() {
			raw := Json(m)
			convey.So(raw, convey.ShouldEqual, ft)
		})
	})
}
