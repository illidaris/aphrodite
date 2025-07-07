package dto

import (
	"context"
	"testing"

	"github.com/illidaris/aphrodite/pkg/contextex"
	"github.com/illidaris/aphrodite/pkg/dependency"
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

func TestBindFrmCtx(t *testing.T) {
	type Demo struct {
		BizRequest
		IPRequest
	}

	type Demo2 struct {
		BizId int64
		IP    string
	}
	// 使用convey库进行测试断言
	convey.Convey("TestBindFrmCtx", t, func() {
		biz := int64(9989)
		ip := "192.168.1.1"
		ctx := context.Background()
		ctx = contextex.WithBizId(ctx, biz)
		ctx = contextex.WithIP(ctx, ip)

		demo := &Demo{}
		demo.BizId = 1
		demo.IP = "x"

		demo2 := &Demo2{}
		demo2.BizId = 2
		demo2.IP = "x"

		convey.Convey("biz test", func() {
			dependency.BindBizFrmCtx(ctx, demo)
			convey.So(demo.BizId, convey.ShouldEqual, biz)

			dependency.BindBizFrmCtx(ctx, demo2)
			convey.So(demo2.BizId, convey.ShouldEqual, 2)
		})

		convey.Convey("ip test", func() {
			dependency.BindIPFrmCtx(ctx, demo)
			convey.So(demo.IP, convey.ShouldEqual, ip)

			dependency.BindIPFrmCtx(ctx, demo2)
			convey.So(demo2.IP, convey.ShouldEqual, "x")
		})
	})
}
