package dto

import (
	"errors"
	"testing"

	"github.com/illidaris/aphrodite/pkg/exception"
	"github.com/smartystreets/goconvey/convey"
)

func TestResponse(t *testing.T) {
	convey.Convey("TestResponse", t, func() {
		convey.Convey("ErrorResponse", func() {
			convey.So(ErrorResponse(errors.New("fail")).Message, convey.ShouldEqual, "fail")
		})
		convey.Convey("OkResponse", func() {
			convey.So(OkResponse(errors.New("fail")).Code, convey.ShouldEqual, 0)
		})
		convey.Convey("NewResponse", func() {
			resp := NewResponse(nil, exception.ERR_BUSI.New("fail"))
			convey.So(resp.Code, convey.ShouldEqual, 30000)
			convey.So(resp.Message, convey.ShouldEqual, "fail")
		})
		convey.Convey("ToException", func() {
			resp := NewResponse(nil, exception.ERR_BUSI_NOFOUND.New("fail"))
			convey.So(resp.Code, convey.ShouldEqual, 30001)
			convey.So(resp.Message, convey.ShouldEqual, "fail")
			ex := resp.ToException()
			convey.So(ex.Code(), convey.ShouldEqual, exception.ERR_BUSI_NOFOUND)
		})
		convey.Convey("ToExceptionFaild", func() {
			resp := NewResponse(nil, exception.ERR_BUSI_NOFOUND.New("fail"))
			resp.Code = 123456
			convey.So(resp.Code, convey.ShouldEqual, 123456)
			convey.So(resp.Message, convey.ShouldEqual, "fail")
			ex := resp.ToException()
			convey.So(ex.Code(), convey.ShouldEqual, 123456)
		})
	})
}
