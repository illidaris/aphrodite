package netex

// IDistrict 国家/地区
type IDistrict interface {
	GetCounyryId() uint16 // 国家/地区
	GetProvinceId() uint8 // 省
	GetCityId() uint8     // 市
}

var _ = IDistrict(District{})

type District struct {
	Country  uint16 `json:"country"`
	Province uint8  `json:"province"`
	City     uint8  `json:"city"`
}

func (r District) GetCounyryId() uint16 {
	return r.Country
}
func (r District) GetProvinceId() uint8 {
	return r.Province
}
func (r District) GetCityId() uint8 {
	return r.City
}

// SplitCityId CityId转成国家Id/省Id/市Id
func SplitCityId(id uint32) (uint16, uint8, uint8) {
	return uint16(id >> 16), uint8(id >> 8), uint8(id)
}

// IDistrictToCityId 国家Id/省Id/市Id转CityId
func DistrictToCityId(country uint16, province, city uint8) uint32 {
	return uint32(country)<<16 | uint32(province)<<8 | uint32(city)
}

func DistrictUniqueFunc[T IDistrict]() func(s T) bool {
	m := map[uint32]struct{}{}
	return func(s T) bool {
		id := DistrictToCityId(s.GetCounyryId(), s.GetProvinceId(), s.GetCityId())
		_, ok := m[id]
		if !ok {
			m[id] = struct{}{}
		}
		return ok
	}
}
