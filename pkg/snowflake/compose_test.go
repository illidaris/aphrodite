package snowflake

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestCompose(t *testing.T) {
	convey.Convey("TestCompose", t, func() {
		convey.Convey("success", func() {
			date := time.Date(2025, 10, 1, 1, 1, 1, 1, time.UTC)
			sequence := 1
			machineId := 1
			gene := 1

			id, err := Compose(date, sequence, machineId, gene)
			convey.So(err, convey.ShouldBeNil)
			convey.So(id, convey.ShouldEqual, 98947242657843216)

			vals := Decompose(id)
			convey.So((date.UnixNano()-defaultStartTime.UnixNano())/defaultTimeUnit, convey.ShouldEqual, vals[0])
			convey.So(sequence, convey.ShouldEqual, vals[1])
			convey.So(machineId, convey.ShouldEqual, vals[2])
			convey.So(gene, convey.ShouldEqual, vals[3])
		})
	})
}
