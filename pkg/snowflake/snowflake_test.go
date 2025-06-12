package snowflake

import (
	"errors"
	"sync"
	"testing"
	"time"

	group "github.com/illidaris/aphrodite/pkg/group/v2"
	"github.com/smartystreets/goconvey/convey"
)

func TestNextIdFunc(t *testing.T) {
	convey.Convey("TestNextIdFunc", t, func() {
		convey.Convey("success", func() {
			opts := []Option{
				WithMachineID(func() int {
					return 111
				}),
				WithStartTime(
					time.Date(2021, 3, 4, 5, 6, 7, 11, time.UTC),
				),
			}
			idGen, _ := NextIdFunc(opts...)
			for i := 0; i < 10; i++ {
				id, err := idGen(nil)
				if err != nil {
					t.Fatalf("failed to generate id: %v", err)
				}
				println(DecomposeStr(id, opts...))
			}
			time.Sleep(time.Second * 1)
			for i := 0; i < 10; i++ {
				id, err := idGen(int64(i))
				if err != nil {
					t.Fatalf("failed to generate id: %v", err)
				}
				println(DecomposeStr(id, opts...))
			}
		})
		convey.Convey("machine id is error", func() {
			opts := []Option{
				WithMachineID(func() int {
					return 129
				}),
				WithStartTime(
					time.Date(2021, 3, 4, 5, 6, 7, 11, time.UTC),
				),
			}
			idGen, err := NextIdFunc(opts...)
			convey.So(err, convey.ShouldEqual, ErrInvalidMachineID)
			convey.So(idGen, convey.ShouldBeNil)
		})
	})
}

func TestNextIdFuncByGoroutine(t *testing.T) {
	convey.Convey("TestNextIdFunc", t, func() {
		opts := []Option{
			WithMachineID(func() int {
				return 123
			}),
			WithStartTime(
				time.Date(2021, 3, 4, 5, 6, 7, 11, time.UTC),
			),
			WithNowFunc(
				func() time.Time {
					return time.Date(2025, 6, 12, 5, 6, 7, 11, time.UTC)
				},
			),
		}
		idGen, err := NextIdFunc(opts...)
		convey.So(err, convey.ShouldBeNil)
		convey.Convey("success", func() {
			raws := []int{}
			m := sync.Map{}
			for i := 0; i < 10000; i++ {
				raws = append(raws, i)
			}
			affect, errs := group.GroupFunc(func(vs ...int) (int64, error) {
				subTotal := int64(0)
				for _, v := range vs {
					id, err := idGen(v)
					if err != nil {
						return 0, err
					}
					_, ok := m.Load(id)
					if ok {
						return 0, errors.New("duplicate")
					}
					m.Store(id, v)
					subTotal++
				}
				return subTotal, nil
			}, raws, group.WithBatch(1))
			convey.So(len(errs), convey.ShouldEqual, 0)
			convey.So(affect, convey.ShouldEqual, len(raws))
		})
	})
}

func BenchmarkNextIdFunc(b *testing.B) {
	opts := []Option{
		WithMachineID(func() int {
			return 111
		}),
	}
	idGen, _ := NextIdFunc(opts...)
	for i := 0; i < b.N; i++ {
		_, err := idGen(i)
		if err != nil {
			b.Fatalf("failed to generate id: %v", err)
		}
	}
}

func TestNextIdFuncByDate(t *testing.T) {
	convey.Convey("TestNextIdFuncByDate", t, func() {
		convey.Convey("success", func() {
			opts := []Option{
				WithMachineID(func() int {
					return 6
				}),
				WithNowFunc(func() time.Time {
					return time.Date(2025, 6, 10, 6, 7, 5, 0, time.UTC)
				}),
			}
			idGen, _ := NextIdFunc(opts...)
			for i := 0; i < 10000; i++ {
				id, err := idGen(i)
				if err != nil {
					t.Fatalf("failed to generate id: %v", err)
				}
				println(DecomposeStr(id, opts...))
			}
		})
	})
}
