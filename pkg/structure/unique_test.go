package structure

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

type demo struct {
	id   string
	sort int
}

func (d *demo) ID() string {
	return d.id
}
func (d *demo) Sort() int {
	return d.Sort()
}
func TestUniqueArray(t *testing.T) {
	convey.Convey("TestUniqueArray", t, func() {
		convey.Convey("NewUniqueArray", func() {
			arr := NewUniqueArray[int]()
			arr.Append(3, 4, 5, 1, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4)
			convey.So(arr.Len(), convey.ShouldEqual, 4)
		})
		convey.Convey("NewUniqueAnyArray", func() {
			arr := NewUniqueAnyArray[string]()
			arr.Append(
				&demo{
					id:   "b2",
					sort: 1,
				},
				&demo{
					id:   "a1",
					sort: 2,
				},
				&demo{
					id:   "b2",
					sort: 3,
				},
				&demo{
					id:   "a1",
					sort: 3,
				})
			convey.So(arr.Len(), convey.ShouldEqual, 2)
		})
	})
}

func F(a, b int) int {
	return a / b
}
func FuzzUniqueArray(f *testing.F) {
	f.Fuzz(func(t *testing.T, a, b int) {
		v := F(a, b)
		println(v)
	})
}
