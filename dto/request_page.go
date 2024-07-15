package dto

import (
	"strings"

	"github.com/illidaris/aphrodite/pkg/convert"
	"github.com/illidaris/aphrodite/pkg/dependency"
)

var _ = dependency.ISortField(&SortField{}) // check impl

// SortField default sort field
type SortField struct {
	Field  string `json:"field" form:"field" url:"field"`    // order by col
	IsDesc bool   `json:"isDesc" form:"isDesc" url:"isDesc"` // asc or desc, default is asc
}

func (r *SortField) GetField() string {
	return r.Field
}

func (r *SortField) GetIsDesc() bool {
	return r.IsDesc
}

var _ = dependency.IPage(&Page{})

// Page default page request
type Page struct {
	PageIndex int64       `json:"page" form:"page" uri:"page" url:"page" binding:"required,gte=1"`                 // currect page no
	AfterId   interface{} `json:"afterId" form:"afterId" url:"afterId"`                                            // previous page last id, when sort by pk
	PageSize  int64       `json:"pageSize" form:"pageSize" uri:"pageSize" url:"pageSize" binding:"required,gte=1"` // page size
	Sorts     []string    `json:"sorts" form:"sorts" uri:"sorts" url:"sorts"`                                      // eg; field|desc
}

func (dto *Page) GetPageIndex() int64 {
	return dto.PageIndex
}

func (dto *Page) GetPageSize() int64 {
	return dto.PageSize
}

func (dto *Page) GetBegin() int64 {
	return (dto.PageIndex - 1) * dto.PageSize
}

func (dto *Page) GetSize() int64 {
	return dto.PageSize
}

func (dto *Page) GetAfterID() any {
	return dto.AfterId
}

func (dto *Page) GetSorts() []dependency.ISortField {
	s := []dependency.ISortField{}
	for _, v := range dto.Sorts {
		words := strings.Split(v, "|")
		if len(words) == 1 {
			w, ok := convert.DefFieldFilter(words[0])
			if !ok {
				continue
			}
			s = append(s, &SortField{Field: w})
		} else if len(words) > 1 {
			w, ok := convert.DefFieldFilter(words[0])
			if !ok {
				continue
			}
			s = append(s, &SortField{Field: w, IsDesc: words[1] == "desc"})
		}
	}
	return s
}
