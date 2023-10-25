package dto

import (
	"encoding/json"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestPage(t *testing.T) {
	convey.Convey("TestPage", t, func() {
		requestStr := "{\"page\":2,\"pageSize\":10,\"sorts\":[\"createAt|desc\",\"id\"]}"
		req := &Page{}
		_ = json.Unmarshal([]byte(requestStr), req)
		convey.Convey("GetPageIndex", func() {
			convey.So(req.GetBegin(), convey.ShouldEqual, 10)
		})
		convey.Convey("GetSorts", func() {
			convey.So(len(req.GetSorts()), convey.ShouldEqual, 2)
			convey.So(req.GetSorts()[0].GetIsDesc(), convey.ShouldBeTrue)
			convey.So(req.GetSorts()[1].GetField(), convey.ShouldEqual, "id")
		})
		convey.Convey("GetAfterID", func() {
			convey.So(req.GetAfterID(), convey.ShouldBeNil)
		})
	})
}
