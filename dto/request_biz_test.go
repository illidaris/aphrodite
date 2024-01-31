package dto

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

// TestIBizRequest 是测试函数，用于测试IBizRequest结构体的功能
func TestIBizRequest(t *testing.T) {
	bizId := int64(1)
	req := &BizRequest{
		BizId: bizId,
	}

	// 使用convey库进行测试断言
	convey.Convey("TestIBizRequest", t, func() {
		convey.Convey("GetBizId", func() {
			convey.So(req.GetBizId(), convey.ShouldEqual, bizId)
		})

		convey.Convey("ToUrlQuery", func() {
			convey.So(req.ToUrlQuery().Encode(), convey.ShouldEqual, "bizId=1")
		})
	})
}
