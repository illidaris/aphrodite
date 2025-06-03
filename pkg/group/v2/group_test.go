package v2

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

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

		total := len(demos)
		batch := 3

		opts := []Option{WithBatch(batch)}
		var count int
		if int(total)%batch == 0 {
			count = int(total) / batch
		} else {
			count = int(total)/batch + 1
		}
		convey.Convey("GroupCount", func() {
			res := Count(demos, opts...)
			convey.So(res, convey.ShouldEqual, count)
		})
		convey.Convey("Group", func() {
			res := Group(demos, opts...)
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
		opts := []Option{WithBatch(batch)}
		convey.Convey("GroupFunc", func() {
			affect, errM := GroupFunc(func(v ...*demo) (int64, error) {
				for _, item := range v {
					println(item.Name)
				}
				return int64(len(v)), nil
			}, demos, opts...)
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
		opts := []Option{WithBatch(batch)}
		convey.Convey("GroupBaseFunc", func() {
			affect, errM := GroupFunc(func(v ...int64) (int64, error) {
				time.Sleep(time.Millisecond * 10)
				for _, item := range v {
					println(item)
				}
				return int64(len(v)), nil
			}, demos, opts...)
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
		opts := []Option{WithBatch(batch)}
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
			}, demos, opts...)
			convey.So(affect, convey.ShouldEqual, curTOtal)
			convey.So(len(errM), convey.ShouldEqual, 2)
		})
	})
}

func TestGroupByAI(t *testing.T) {
	// TC1: Empty input
	t.Run("empty slice", func(t *testing.T) {
		if res := Group([]int{}); len(res) != 0 {
			t.Errorf("Expected empty groups, got %v", res)
		}
	})

	// TC2: Total < default batch (default 100)
	t.Run("small slice", func(t *testing.T) {
		input := []int{1, 2, 3}
		res := Group(input)
		if len(res) != 3 || len(res[0]) > 1 {
			t.Errorf("Expected 3 group, got %v", res)
		}
	})

	// TC3: Custom batch size
	t.Run("custom batch", func(t *testing.T) {
		input := make([]int, 10)
		res := Group(input, WithBatch(3))
		expectedGroups := 4 // 3+3+3+1
		if len(res) != expectedGroups {
			t.Errorf("Expected %d groups, got %d", expectedGroups, len(res))
		}
	})

	// TC4: Minimum group constraint
	t.Run("min groups", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		res := Group(input, WithParallelismMax(3))
		if len(res) != 3 {
			t.Errorf("Expected 3 groups, got %d", len(res))
		}
	})
}

func TestGroupFuncByAI(t *testing.T) {
	// TC5: All success
	t.Run("all success", func(t *testing.T) {
		mockFn := func(v ...int) (int64, error) {
			return int64(len(v)), nil
		}

		input := []int{1, 2, 3, 4, 5}
		affect, errs := GroupFunc(mockFn, input, WithBatch(2))

		if affect != 5 {
			t.Errorf("Expected affect 5, got %d", affect)
		}
		if len(errs) != 0 {
			t.Errorf("Expected no errors, got %v", errs)
		}
	})

	// TC6: Partial errors
	t.Run("partial errors", func(t *testing.T) {
		mockFn := func(v ...int) (int64, error) {
			if v[0]%2 == 0 {
				return 0, errors.New("even error")
			}
			return 1, nil
		}

		input := []int{1, 2, 3, 4}
		_, errs := GroupFunc(mockFn, input, WithBatch(1))

		expectedErrors := 2 // Groups 2 and 4 (0-based index 1 and 3)
		if len(errs) != expectedErrors {
			t.Errorf("Expected %d errors, got %d", expectedErrors, len(errs))
		}
	})

	// TC7: Panic recovery
	t.Run("panic handling", func(t *testing.T) {
		mockFn := func(v ...int) (int64, error) {
			panic("test panic")
		}

		input := []int{1, 2}
		_, errs := GroupFunc(mockFn, input, WithBatch(1))

		if len(errs) != 2 {
			t.Errorf("Expected 2 errors, got %d", len(errs))
		}
		for _, err := range errs {
			if err.Error() != "err:test panic" {
				t.Errorf("Unexpected error message: %v", err)
			}
		}
	})

	// TC8: Mixed scenario
	t.Run("mixed results", func(t *testing.T) {
		mockFn := func(v ...int) (int64, error) {
			switch v[0] {
			case 1:
				return 1, nil
			case 2:
				return 0, errors.New("error2")
			case 3:
				panic("panic3")
			default:
				return 0, nil
			}
		}

		input := []int{1, 2, 3, 4}
		affect, errs := GroupFunc(mockFn, input, WithBatch(1))

		if affect != 1 {
			t.Errorf("Expected affect 1, got %d", affect)
		}
		if len(errs) != 2 {
			t.Errorf("Expected 2 errors, got %d", len(errs))
		}
	})
}
