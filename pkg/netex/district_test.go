package netex

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestSplitCityId(t *testing.T) {
	convey.Convey("TestSplitCityId", t, func() {
		convey.Convey("156:45:1", func() {
			countryId, provinceId, cityId := SplitCityId(10235137)
			println(countryId, provinceId, cityId)
			convey.So(countryId, convey.ShouldEqual, 156)
			convey.So(provinceId, convey.ShouldEqual, 45)
			convey.So(cityId, convey.ShouldEqual, 1)
		})
		convey.Convey("156:31:0", func() {
			countryId, provinceId, cityId := SplitCityId(10231552)
			println(countryId, provinceId, cityId)
			convey.So(countryId, convey.ShouldEqual, 156)
			convey.So(provinceId, convey.ShouldEqual, 31)
			convey.So(cityId, convey.ShouldEqual, 0)
		})
		convey.Convey("158:71:0", func() {
			countryId, provinceId, cityId := SplitCityId(10372864)
			println(countryId, provinceId, cityId)
			convey.So(countryId, convey.ShouldEqual, 158)
			convey.So(provinceId, convey.ShouldEqual, 71)
			convey.So(cityId, convey.ShouldEqual, 0)
		})
		convey.Convey("156:41:12", func() {
			countryId, provinceId, cityId := SplitCityId(10234124)
			println(countryId, provinceId, cityId)
			convey.So(countryId, convey.ShouldEqual, 156)
			convey.So(provinceId, convey.ShouldEqual, 41)
			convey.So(cityId, convey.ShouldEqual, 12)
		})
	})
}

func TestDistrictToCityId(t *testing.T) {
	tests := []struct {
		country  uint16
		province uint8
		city     uint8
		wantId   uint32
	}{
		{156, 31, 0, 10231552},
		{158, 71, 0, 10372864},
		{156, 45, 1, 10235137},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			gotId := DistrictToCityId(tt.country, tt.province, tt.city)
			if gotId != tt.wantId {
				t.Errorf("DistrictToCityId(%v, %v, %v) = %v; want %v", tt.country, tt.province, tt.city, gotId, tt.wantId)
			}
		})
	}
}
