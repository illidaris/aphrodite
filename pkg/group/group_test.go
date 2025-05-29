package group

import (
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestCount(t *testing.T) {
	convey.Convey("TestGroup", t, func() {
		convey.Convey("CountBase batch 0 total 1000", func() {
			datas := []int64{}
			for i := 0; i < 1000; i++ {
				datas = append(datas, int64(i))
			}
			res := CountBase(0, datas...)
			convey.So(res, convey.ShouldEqual, 100)
		})
		convey.Convey("CountBase batch 0 total 99", func() {
			datas := []int64{}
			for i := 0; i < 99; i++ {
				datas = append(datas, int64(i))
			}
			res := CountBase(0, datas...)
			convey.So(res, convey.ShouldEqual, 99)
		})
		convey.Convey("CountBase batch 1 total 10000", func() {
			datas := []int64{}
			for i := 0; i < 10000; i++ {
				datas = append(datas, int64(i))
			}
			res := CountBase(1, datas...)
			convey.So(res, convey.ShouldEqual, 1000)
		})
		convey.Convey("CountBase batch 1 total 99", func() {
			datas := []int64{}
			for i := 0; i < 99; i++ {
				datas = append(datas, int64(i))
			}
			res := CountBase(1, datas...)
			convey.So(res, convey.ShouldEqual, 99)
		})
		convey.Convey("CountBase batch 3 total 1000", func() {
			datas := []int64{}
			for i := 0; i < 1000; i++ {
				datas = append(datas, int64(i))
			}
			res := CountBase(3, datas...)
			convey.So(res, convey.ShouldEqual, 334)
		})
		convey.Convey("CountBase batch 3 total 99", func() {
			datas := []int64{}
			for i := 0; i < 99; i++ {
				datas = append(datas, int64(i))
			}
			res := CountBase(3, datas...)
			convey.So(res, convey.ShouldEqual, 99)
		})
	})
}

// TestGroup
func TestGroup(t *testing.T) {
	convey.Convey("TestGroup", t, func() {
		type demo struct {
			Name string
			Age  int
		}
		demos := []*demo{
			{Name: "x4", Age: 4},
			{Name: "x3", Age: 3},
			{Name: "x2", Age: 2},
			{Name: "x5", Age: 5},
			{Name: "x1", Age: 1},
			{Name: "x6", Age: 6},
			{Name: "x7", Age: 7},
		}
		batch := 3
		total := len(demos)
		var count int
		if int(total)%batch == 0 {
			count = int(total) / batch
		} else {
			count = int(total)/batch + 1
		}
		convey.Convey("GroupCount", func() {
			res := Count(batch, demos...)
			convey.So(res, convey.ShouldEqual, count)
		})
		convey.Convey("Group", func() {
			res := Group(batch, demos...)
			for gId, g := range res {
				for iId, i := range g {
					d := demos[batch*gId+iId]
					convey.So(i.Name, convey.ShouldEqual, d.Name)
					convey.So(i.Age, convey.ShouldEqual, d.Age)
				}
			}
		})
	})
}

// TestGroupFunc
func TestGroupFunc(t *testing.T) {
	convey.Convey("TestGroup", t, func() {
		type demo struct {
			Name string
			Age  int
		}
		demos := []*demo{
			{Name: "x4", Age: 4},
			{Name: "x3", Age: 3},
			{Name: "x2", Age: 2},
			{Name: "x5", Age: 5},
			{Name: "x1", Age: 1},
			{Name: "x6", Age: 6},
			{Name: "x7", Age: 7},
		}
		batch := 3
		total := len(demos)
		convey.Convey("GroupFunc", func() {
			affect, errM := GroupFunc(func(v ...*demo) (int64, error) {
				for _, item := range v {
					println(item.Name)
				}
				return int64(len(v)), nil
			}, batch, demos...)
			convey.So(affect, convey.ShouldEqual, total)
			convey.So(len(errM), convey.ShouldEqual, 0)
		})
	})
}

// TestGroupFunc
func TestGroupBaseFunc(t *testing.T) {
	convey.Convey("TestGroupBase", t, func() {
		demos := []int64{
			5, 7, 8, 9, 1, 2, 3, 11, 55, 88,
		}
		batch := 2
		total := len(demos)
		convey.Convey("GroupBaseFunc", func() {
			affect, errM := GroupBaseFunc(func(v ...int64) (int64, error) {
				time.Sleep(time.Millisecond * 10)
				for _, item := range v {
					println(item)
				}
				return int64(len(v)), nil
			}, batch, demos...)
			convey.So(affect, convey.ShouldEqual, total)
			convey.So(len(errM), convey.ShouldEqual, 0)
		})
	})
}

// TestGroupFuncWithErr
func TestGroupFuncWithErr(t *testing.T) {
	convey.Convey("TestGroup", t, func() {
		type demo struct {
			Name string
			Age  int
		}
		demos := []*demo{
			{Name: "x4", Age: 4},
			{Name: "x3", Age: 3},
			{Name: "x2", Age: 2},
			{Name: "x5", Age: 5},
			{Name: "x1", Age: 1},
			{Name: "x6", Age: 6},
			{Name: "x7", Age: 7},
		}
		batch := 3
		curTOtal := 0
		for _, v := range demos {
			if v.Age <= 5 {
				curTOtal += 1
			}
		}
		convey.Convey("GroupFunc_Error", func() {
			affect, errM := GroupFunc(func(v ...*demo) (int64, error) {
				var err error
				result := []*demo{}
				for _, item := range v {
					if item.Age > 5 {
						err = fmt.Errorf("find err: %d", item.Age)
						continue
					}
					result = append(result, item)
					println(item.Name)
				}
				return int64(len(result)), err
			}, batch, demos...)
			convey.So(affect, convey.ShouldEqual, curTOtal)
			convey.So(len(errM), convey.ShouldEqual, 2)
		})
	})
}
