package dto

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewDataPager(t *testing.T) {
	convey.Convey("TestNewDataPager", t, func() {
		convey.Convey("NewDataPager", func() {
			pager := NewDataPager(nil, 2, 10, 20)
			convey.So(pager.TotalPage, convey.ShouldEqual, 2)
			convey.So(pager.GetTotalRecord(), convey.ShouldEqual, 20)
		})
	})
}
