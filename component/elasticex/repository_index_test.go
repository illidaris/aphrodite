package elasticex

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestBaseCreate(t *testing.T) {
	mockESX(func() {
		convey.Convey("TestBaseCreate", t, func() {
			ctx := context.Background()
			newRepo := BaseRepository[testStructShardingPo]{}
			affect, err := newRepo.BaseCreate(ctx, []*testStructShardingPo{{Id: "1"}})
			convey.So(err, convey.ShouldBeNil)
			convey.So(affect, convey.ShouldEqual, 1)
		})
	})
}

func TestBaseSave(t *testing.T) {
	mockESX(func() {
		convey.Convey("TestBaseSave", t, func() {
			ctx := context.Background()
			newRepo := BaseRepository[testStructShardingPo]{}
			affect, err := newRepo.BaseSave(ctx, []*testStructShardingPo{{Id: "1"}})
			convey.So(err, convey.ShouldBeNil)
			convey.So(affect, convey.ShouldEqual, 1)
		})
	})
}

func TestBaseUpdate(t *testing.T) {
	mockESX(func() {
		convey.Convey("TestBaseUpdate", t, func() {
			ctx := context.Background()
			newRepo := BaseRepository[testStructShardingPo]{}
			affect, err := newRepo.BaseUpdate(ctx, &testStructShardingPo{Id: "1", Code: "xx"})
			convey.So(err, convey.ShouldBeNil)
			convey.So(affect, convey.ShouldEqual, 1)
		})
	})
}

func TestBaseDelete(t *testing.T) {
	mockESX(func() {
		convey.Convey("TestBaseDelete", t, func() {
			ctx := context.Background()
			newRepo := BaseRepository[testStructShardingPo]{}
			affect, err := newRepo.BaseDelete(ctx, &testStructShardingPo{Id: "1"})
			convey.So(err, convey.ShouldBeNil)
			convey.So(affect, convey.ShouldEqual, 1)
		})
	})
}
