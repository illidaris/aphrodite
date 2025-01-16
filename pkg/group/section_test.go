package group

import (
	"encoding/json"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

type testDemo struct {
	X int64
	Y int64
}

func (t *testDemo) GetBeg() int64 {
	return t.X
}
func (t *testDemo) GetEnd() int64 {
	return t.Y
}
func TestSection(t *testing.T) {
	convey.Convey("TestSection", t, func() {
		res := Section(1, 8, 2)
		right := [][]int64{
			{1, 3},
			{3, 5},
			{5, 7},
			{7, 8},
		}
		convey.So(right, convey.ShouldEqual, res)
	})
	convey.Convey("TestSection failed", t, func() {
		res := Section(7, 6, 2)
		convey.So(len(res), convey.ShouldEqual, 0)
	})
}

func TestSectionAny(t *testing.T) {
	convey.Convey("TestSectionAny", t, func() {
		t := &testDemo{
			X: 1,
			Y: 10,
		}
		res := SectionAny(t, 2, func(b, e int64) *testDemo {
			return &testDemo{
				X: b,
				Y: e,
			}
		})
		rigntStr := `[{"X":1,"Y":3},{"X":3,"Y":5},{"X":5,"Y":7},{"X":7,"Y":9},{"X":9,"Y":10}]`
		right := []*testDemo{}
		json.Unmarshal([]byte(rigntStr), &right)
		convey.So(right, convey.ShouldEqual, res)
	})
	convey.Convey("TestSectionAny failed", t, func() {
		t := &testDemo{
			X: 11,
			Y: 10,
		}
		res := SectionAny(t, 2, func(b, e int64) *testDemo {
			return &testDemo{
				X: b,
				Y: e,
			}
		})
		convey.So(len(res), convey.ShouldEqual, 0)
	})
}
