package convert

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestConvertToStruct(t *testing.T) {
	convey.Convey("TestConvertToStruct", t, func() {
		type Demo struct {
			Name  string    `json:"name"`
			Age   int       `json:"age"`
			Birth time.Time `json:"birth"`
		}
		demo := &Demo{
			Name:  "xxx",
			Age:   123,
			Birth: time.Now(),
		}
		jsBs, _ := json.Marshal(demo)
		var obj interface{}
		_ = json.Unmarshal(jsBs, &obj)

		convey.Convey("ConvertToStructByRef", func() {
			// res, err := ConvertToStructByRef(obj, reflect.TypeOf(Demo{}))
			// convey.So(err, convey.ShouldBeNil)

			// demoV, ok := res.(Demo)
			// convey.So(ok, convey.ShouldBeTrue)
			// convey.So(demoV.Name, convey.ShouldEqual, demo.Name)
			// convey.So(demoV.Age, convey.ShouldEqual, demo.Age)
			// convey.So(demoV.Birth.Unix(), convey.ShouldEqual, demo.Birth.Unix())
		})

		convey.Convey("ConvertToStructByJson", func() {
			res, err := ConvertToStructByJson(obj, reflect.TypeOf(Demo{}))
			convey.So(err, convey.ShouldBeNil)

			demoV, ok := res.(*Demo)
			convey.So(ok, convey.ShouldBeTrue)
			convey.So(demoV.Name, convey.ShouldEqual, demo.Name)
			convey.So(demoV.Age, convey.ShouldEqual, demo.Age)
			convey.So(demoV.Birth.Unix(), convey.ShouldEqual, demo.Birth.Unix())
		})
	})
}
