package dto

import (
	"testing"

	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/smartystreets/goconvey/convey"
)

func TestNewDataPager(t *testing.T) {
	convey.Convey("TestNewDataPager", t, func() {
		convey.Convey("NewDataPager", func() {
			pager := NewDataPager(nil, 2, 10, 20)
			convey.So(pager.TotalPage, convey.ShouldEqual, 2)
			convey.So(pager.GetTotal(), convey.ShouldEqual, 2)
		})
	})
}

func TestNewRecordPager(t *testing.T) {
	convey.Convey("TestNewRecordPager", t, func() {
		convey.Convey("NewRecordPager", func() {
			pager := NewRecordPager(2, 10, 20, dependency.EmptyPo{})
			convey.So(pager.TotalPage, convey.ShouldEqual, 2)
			convey.So(pager.GetTotal(), convey.ShouldEqual, 2)
		})
	})
}
