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
	})
}
