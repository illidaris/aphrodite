package embedded

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestIComponent(t *testing.T) {
	convey.Convey("TestIComponent", t, func() {
		convey.Convey("TestIComponentNil", func() {
			testComponent := NewComponent[int]()
			writer := testComponent.GetWriter("demo")
			convey.So(writer, convey.ShouldBeGreaterThanOrEqualTo, 0)
			reader := testComponent.GetReader("demo")
			convey.So(reader, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
		convey.Convey("TestIComponentSingle", func() {
			testComponent := NewComponent[int]()
			testComponent.NewWriter("demo", 4)
			writer := testComponent.GetWriter("demo")
			convey.So(writer, convey.ShouldBeGreaterThanOrEqualTo, 4)
		})
		convey.Convey("TestIComponentMulti", func() {
			testComponent := NewComponent[int]()
			testComponent.NewReader("demo", []int{1, 2, 3}...)
			testComponent.NewWriter("demo", []int{4, 5, 6}...)
			reader := testComponent.GetReader("demo")
			convey.So(reader, convey.ShouldBeLessThan, 4)
			writer := testComponent.GetWriter("demo")
			convey.So(writer, convey.ShouldBeGreaterThanOrEqualTo, 3)
		})
		convey.Convey("TestIComponentBalance", func() {
			testComponent := NewComponent[int]()
			testComponent.SetReaderBalance(func(ts ...IInstance[int]) IInstance[int] {
				return NewInstance[int]("demo", 666)
			})
			reader := testComponent.GetReader("demo")
			convey.So(reader, convey.ShouldBeGreaterThanOrEqualTo, 666)
			testComponent.SetWriterBalance(func(ts ...IInstance[int]) IInstance[int] {
				return NewInstance[int]("demo", 777)
			})
			writer := testComponent.GetWriter("demo")
			convey.So(writer, convey.ShouldBeGreaterThanOrEqualTo, 777)
		})
	})
}
