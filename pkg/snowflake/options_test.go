package snowflake

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestOptionsOffset(t *testing.T) {
	convey.Convey("TestOptionsOffset", t, func() {
		opt := newOptions()
		convey.Convey("Success", func() {
			results := []int64{63, 22, 21, 11, 4}
			for index, res := range results {
				v := Offset(opt.LenSlice(), index)
				convey.So(v, convey.ShouldEqual, res)
			}
		})
		convey.Convey("Out Of Range", func() {
			convey.So(Offset(opt.LenSlice(), 1000), convey.ShouldEqual, 0)
		})
	})
}

func TestOptions(t *testing.T) {
	convey.Convey("TestOptions", t, func() {
		opt := newOptions()
		vals := []int64{1, 1, 3, 4, 5}
		vs := IdPartsFrmVals(opt.LenSlice(), vals...)

		convey.Convey("IdPartsFrmVals", func() {
			results := []int64{4194304, 2097152, 6144, 64, 5}
			convey.So(len(vs), convey.ShouldEqual, len(results))
			convey.So(vs, convey.ShouldEqual, results)
		})

		convey.Convey("GetValsFrmId", func() {
			id := int64(0)
			for _, v := range vs {
				id |= v
			}
			convey.So(GetValsFrmId(opt.LenSlice(), id), convey.ShouldEqual, vals)
		})
	})
}
