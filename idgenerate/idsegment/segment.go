package idsegment

type Segment struct {
	Code   StatusCode `json:"code" form:"code" url:"code"`
	Cursor uint64     `json:"cursor" form:"cursor" url:"cursor"`
	MaxId  uint64     `json:"max" form:"max" url:"max"`
	MinId  uint64     `json:"min" form:"min" url:"min"`
}
