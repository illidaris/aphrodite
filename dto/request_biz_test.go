package dto

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIBizRequest(t *testing.T) {
	bizId := int64(1)
	req := &BizRequest{
		BizId: bizId,
	}
	convey.Convey("TestIBizRequest", t, func() {
		convey.Convey("GetBizId", func() {
			convey.So(req.GetBizId(), convey.ShouldEqual, bizId)
		})
		convey.Convey("ToUrlQuery", func() {
			convey.So(req.ToUrlQuery().Encode(), convey.ShouldEqual, "bizId=1")
		})
	})
}
