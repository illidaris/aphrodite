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
			convey.So(pager.GetTotal(), convey.ShouldEqual, 20)
		})
	})
}

func TestNewRowPager(t *testing.T) {
	convey.Convey("TestNewRowPager", t, func() {
		convey.Convey("NewRowPager", func() {
			pager := NewRowPager(2, 10, 20)
			convey.So(pager.TotalPage, convey.ShouldEqual, 2)
			convey.So(pager.GetTotal(), convey.ShouldEqual, 20)
		})
	})
}

func TestNewRecordPager(t *testing.T) {
	convey.Convey("TestNewRecordPager", t, func() {
		convey.Convey("NewRecordPager", func() {
			pager := NewRecordPager(2, 10, 20, 1)
			convey.So(pager.TotalPage, convey.ShouldEqual, 2)
			convey.So(pager.GetTotal(), convey.ShouldEqual, 20)
		})
	})
}
