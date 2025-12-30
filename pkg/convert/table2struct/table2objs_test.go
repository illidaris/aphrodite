package table2struct

import (
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cast"
)

func TestObjs2Table(t *testing.T) {
	type Display struct {
		Code string `json:"code" gorm:"comment:编码"`
		Name string `json:"name" gorm:"comment:名称"`
	}
	type DemoSection struct {
		CreateAt int64 `json:"createAt" gorm:"comment:创建时间"`
		CreateBy int64 `json:"createBy" gorm:"comment:创建人"`
	}
	type RangeSection struct {
		Beg    int64  `json:"beg" gorm:"comment:开始"`
		End    int64  `json:"end" gorm:"comment:范围结束结束"`
		Prefix string `json:"prefix" gorm:"comment:范围前缀"`
	}
	type Demo struct {
		Id          int64  `json:"id" gorm:"comment:id"`
		Cover       string `json:"cover" gorm:"comment:背景图"`
		Display     `gorm:"extend"`
		DemoSection `gorm:"extend"`
		Range       RangeSection  `json:"range" gorm:"comment:范围"`
		ShowRange   *RangeSection `json:"showRange" gorm:"comment:范围"`
	}

	convey.Convey("TestObjs2Table", t, func() {
		demos := []interface{}{}
		raws := []Demo{}
		for index, v := range []int64{1} {
			demo := Demo{
				ShowRange: &RangeSection{},
			}
			demo.Id = v
			demo.Cover = fmt.Sprintf("http://baiu.com/img/%v", v)
			demo.CreateBy = 999
			demo.CreateAt = time.Now().Unix()
			demo.Code = cast.ToString(index)
			demo.Name = fmt.Sprintf("测试数据%v", index)
			demo.Range.Beg = 1
			demo.Range.End = 3
			demo.Range.Prefix = "px"
			demo.ShowRange.Beg = 2
			demo.ShowRange.End = 6
			demo.ShowRange.Prefix = "asdas"
			demos = append(demos, demo)
			raws = append(raws, demo)
		}

		convey.Convey("success", func() {
			headers, rows, err := Objs2Table(demos, WithDeep())
			convey.So(err, convey.ShouldBeNil)
			all := [][]string{}
			all = append(all, headers...)
			all = append(all, rows...)
			target := []Demo{}
			Table2Objs(&target, all, WithDeep(), WithStartRowIndex(2), WithHeadIndex(1))
			convey.So(raws, convey.ShouldEqual, target)
		})

		convey.Convey("x", func() {
			headers, rows, err := Objs2Table(demos,
				WithCustom("x1", "字段1", func(i interface{}) string {
					return "1"
				}),
				WithCustom("x2", "字段2", func(i interface{}) string {
					return "2"
				}),
			)
			convey.So(err, convey.ShouldBeNil)
			all := [][]string{}
			all = append(all, headers...)
			all = append(all, rows...)
			println(all)
		})
	})
}
