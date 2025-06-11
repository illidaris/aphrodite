package idsegment

type Segment struct {
	Code   StatusCode `json:"code" form:"code" url:"code"`
	Cursor int64      `json:"cursor" form:"cursor" url:"cursor"`
	MaxId  int64      `json:"max" form:"max" url:"max"`
	MinId  int64      `json:"min" form:"min" url:"min"`
}
