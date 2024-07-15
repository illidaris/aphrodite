package dto

import "github.com/illidaris/aphrodite/pkg/dependency"

var _ = dependency.IRange(&Range{})

// Range range
type Range struct {
	Beg int64 `json:"beg" form:"beg" url:"beg"`
	End int64 `json:"end" form:"end" url:"end"`
}

func (i *Range) GetBeg() int64 {
	return i.Beg
}
func (i *Range) GetEnd() int64 {
	return i.End
}
