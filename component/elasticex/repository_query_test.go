package elasticex

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestBaseGet(t *testing.T) {
	mockESX(func() {
		convey.Convey("BaseGet", t, func() {
			ctx := context.Background()
			newRepo := BaseRepository[testStructShardingPo]{}
			res, err := newRepo.BaseGet(ctx)
			convey.So(err, convey.ShouldBeNil)
			convey.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestBaseCount(t *testing.T) {
	mockESX(func() {
		convey.Convey("BaseCount", t, func() {
			ctx := context.Background()
			newRepo := BaseRepository[testStructShardingPo]{}
			affect, err := newRepo.BaseCount(ctx)
			convey.So(err, convey.ShouldBeNil)
			convey.So(affect, convey.ShouldEqual, 1)
		})
	})
}

func TestBaseQuery(t *testing.T) {
	mockESX(func() {
		convey.Convey("BaseQuery ", t, func() {
			ctx := context.Background()
			newRepo := BaseRepository[testStructShardingPo]{}
			ps, err := newRepo.BaseQuery(ctx)
			convey.So(err, convey.ShouldBeNil)
			convey.So(len(ps), convey.ShouldEqual, 1)
		})
	})
}

func TestBaseQueryWithCount(t *testing.T) {
	mockESX(func() {
		convey.Convey("BaseQueryWithCount ", t, func() {
			ctx := context.Background()
			newRepo := BaseRepository[testStructShardingPo]{}
			ps, affect, err := newRepo.BaseQueryWithCount(ctx)
			convey.So(err, convey.ShouldBeNil)
			convey.So(affect, convey.ShouldEqual, 1)
			convey.So(len(ps), convey.ShouldEqual, 1)
		})
	})
}
