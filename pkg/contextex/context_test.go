package contextex

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestID(t *testing.T) {
	convey.Convey("TestContextKeyID", t, func() {
		id := "aaaaaxxx"
		v := ElasticID.ID(id)
		convey.So(string(v), convey.ShouldEqual, "_aphrodite_es"+"_"+id)
	})
	convey.Convey("TestContextKeyID", t, func() {
		id := "asdasdxxx"
		v := DbTxID.ID(id)
		convey.So(string(v), convey.ShouldEqual, "_aphrodite_dbtx"+"_"+id)
	})
}
