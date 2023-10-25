package dto

import (
	"encoding/json"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestRange(t *testing.T) {
	requestStr := "{\"beg\":2,\"end\":10}"
	req := &Range{}
	_ = json.Unmarshal([]byte(requestStr), req)
	convey.Convey("TestRange", t, func() {
		convey.Convey("GetBeg", func() {
			convey.So(req.GetBeg(), convey.ShouldEqual, 2)
		})
		convey.Convey("GetEnd", func() {
			convey.So(req.GetEnd(), convey.ShouldEqual, 10)
		})
	})
}
