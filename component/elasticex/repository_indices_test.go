package elasticex

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestBaseTableCreate(t *testing.T) {
	mockESX(func() {
		convey.Convey("TestBaseTableCreate ", t, func() {
			ctx := context.Background()
			newRepo := BaseRepository[testStructShardingPo]{}
			affect, err := newRepo.BaseTableCreate(ctx)
			convey.So(err, convey.ShouldBeNil)
			convey.So(affect, convey.ShouldEqual, 1)
		})
	})
}

func TestBaseTableExists(t *testing.T) {
	mockESX(func() {
		convey.Convey("BaseTableExists ", t, func() {
			ctx := context.Background()
			newRepo := BaseRepository[testStructShardingPo]{}
			exist, err := newRepo.BaseTableExists(ctx)
			convey.So(err, convey.ShouldBeNil)
			convey.So(exist, convey.ShouldBeTrue)
		})
	})
}
