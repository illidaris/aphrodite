package idsnow

import (
	"math"
	"math/rand"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestOffset(t *testing.T) {
	convey.Convey("TestOffset", t, func() {
		convey.Convey("success", func() {
			lens := []int{3, 4, 6, 11}
			convey.So(Offset(lens, 0), convey.ShouldEqual, 24)
			convey.So(Offset(lens, 2), convey.ShouldEqual, 17)
			convey.So(Offset(lens, 10), convey.ShouldEqual, 0)
		})
	})
}

func TestIdPartsFrmVals(t *testing.T) {
	convey.Convey("TestIdPartsFrmVals", t, func() {
		convey.Convey("success", func() {
			lens := []int{
				defaultBitsTime,
				defaultBitsClock,
				defaultBitsSequence,
				defaultBitsMachine,
				defaultBitGene,
			}
			vals := []int64{
				123456789,
				0,
				rand.Int63n(int64(math.Pow(2, defaultBitsSequence))),
				rand.Int63n(int64(math.Pow(2, defaultBitsMachine))),
				rand.Int63n(int64(math.Pow(2, defaultBitGene))),
			}
			// 组装
			parts := IdPartsFrmVals(lens, vals...)
			id := int64(0)
			for _, v := range parts {
				id |= v
			}
			// 分割
			resVals := GetValsFrmId(lens, id)
			convey.So(resVals, convey.ShouldEqual, vals)
		})
	})
}
