package idsnow

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/smartystreets/goconvey/convey"
)

func TestMachineIdDistribute(t *testing.T) {
	convey.Convey("测试机器分配", t, func() {
		db, mock := redismock.NewClientMock()

		mock.ExpectLPush("ap.idsnow.mids", 1, 2, 3).SetVal(1)
		affect, err := db.LPush(context.Background(), "ap.idsnow.mids", 1, 2, 3).Result()
		convey.So(err, convey.ShouldBeNil)
		convey.So(affect, convey.ShouldEqual, 1)

		err = mock.ExpectationsWereMet()
		convey.So(err, convey.ShouldBeNil)
	})
}
