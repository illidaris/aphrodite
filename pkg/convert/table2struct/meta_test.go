package table2struct

import (
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cast"
)

func TestFields(t *testing.T) {
	type Display struct {
		Code string `json:"code" gorm:"comment:编码"`
		Name string `json:"name" gorm:"comment:名称"`
	}

	type DemoSection struct {
		CreateAt int64 `json:"createAt" gorm:"comment:创建时间"`
		CreateBy int64 `json:"createBy" gorm:"comment:创建人"`
	}

	type RangeSection struct {
		Beg    int64  `json:"beg" `
		End    int64  `json:"end"`
		Prefix string `json:"prefix"`
	}

	type Demo struct {
		Id          int64  `json:"id" gorm:"comment:id"`
		Cover       string `json:"cover" gorm:"comment:背景图"`
		Display     `gorm:"extend"`
		DemoSection `gorm:"extend"`
		Ranges      RangeSection `json:"range" gorm:"comment:范围"`
	}

	convey.Convey("TestFields", t, func() {
		demos := []interface{}{}
		for index, v := range []int64{1, 99} {
			demo := Demo{}
			demo.Id = v
			demo.Cover = fmt.Sprintf("http://baiu.com/img/%v", v)
			demo.CreateBy = 999
			demo.CreateAt = time.Now().Unix()
			demo.Code = cast.ToString(index)
			demo.Name = fmt.Sprintf("测试数据%v", index)
			demos = append(demos, demo)
		}
		convey.Convey("success", func() {
			for _, v := range demos {
				_, err := Fields(v, true, "")
				if err != nil {
					convey.So(err, convey.ShouldBeNil)
				}
			}
		})
	})
}
