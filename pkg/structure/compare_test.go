package structure

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var _ = IIDSection[int](&mockIIDSection{})

// Define a mock type for IIDSection to use in testing
type mockIIDSection struct {
	id int
}

// Implement the ID() method required by IIDSection
func (m mockIIDSection) ID() int {
	return m.id
}

// TestCompareIDArray tests the CompareIDArray function
func TestCompareIDArray(t *testing.T) {
	convey.Convey("TestCompareIDArray", t, func() {
		srcs := []IIDSection[int]{
			mockIIDSection{id: 1},
			mockIIDSection{id: 1},
			mockIIDSection{id: 2},
			mockIIDSection{id: 2},
			mockIIDSection{id: 3},
			mockIIDSection{id: 3},
			mockIIDSection{id: 4},
			mockIIDSection{id: 4},
			mockIIDSection{id: 9},
			mockIIDSection{id: 10},
		}
		targets := []IIDSection[int]{
			mockIIDSection{10},
			mockIIDSection{4},
			mockIIDSection{10},
			mockIIDSection{2},
			mockIIDSection{3},
			mockIIDSection{10},
			mockIIDSection{11},
			mockIIDSection{12},
		}
		convey.Convey("default", func() {
			count := CompareIDArray(srcs, targets)
			convey.So(count, convey.ShouldEqual, 7)
		})
		convey.Convey("max is 1", func() {
			count := CompareIDArray(srcs, targets, WithMax[IIDSection[int]](1))
			convey.So(count, convey.ShouldEqual, 1)
		})
		convey.Convey("duplicate removal", func() {
			m := NewUnqueFilter[int]()
			right := []int{2, 3, 4, 10}
			result := []int{}
			count := CompareIDArray(srcs, targets,
				WithIterator(func(s, t IIDSection[int]) {
					result = append(result, s.ID())
				}),
				WithFilter[IIDSection[int]](func(s IIDSection[int]) bool {
					return m(s.ID())
				}))
			convey.So(count, convey.ShouldEqual, 4)
			for i := 0; i < int(count); i++ {
				convey.So(result[i], convey.ShouldEqual, right[i])
			}
		})
	})
}

// TestCompareArray tests the CompareArray function
func TestCompareArray(t *testing.T) {
	// Create mock source and target arrays
	src := []int{1, 2, 3}
	target := []int{2, 3, 4}

	// Call CompareArray with no options
	count := CompareArray(src, target)
	if count != 2 {
		t.Errorf("Expected 2 matching elements, got %d", count)
	}

	// Call CompareArray with a custom option to limit the maximum number of comparisons
	count = CompareArray(src, target, WithMax[int](1))
	if count != 1 {
		t.Errorf("Expected 1 matching element with max option, got %d", count)
	}
}
