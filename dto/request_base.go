package dto

import (
	"net/url"

	"github.com/google/go-querystring/query"
)

type DetailRequest struct {
	BizRequest
	IPRequest
	Id      int64  `json:"id" bson:"_id,omitempty" form:"id" uri:"id" url:"id,omitempty"`
	Code    string `json:"code" bson:"code,omitempty" form:"code" uri:"code" url:"code,omitempty"`
	Name    string `json:"name" bson:"name,omitempty" form:"name" uri:"name" url:"name,omitempty"`
	Keyword string `json:"keyword" bson:"keyword,omitempty" form:"keyword" uri:"keyword" url:"keyword,omitempty"`
}

func (r DetailRequest) ToUrlQuery() url.Values {
	u, _ := query.Values(r)
	return u
}

func (r DetailRequest) GetUrlQuery() url.Values {
	return r.ToUrlQuery()
}

func (r DetailRequest) Encode() ([]byte, error) {
	return nil, nil
}
