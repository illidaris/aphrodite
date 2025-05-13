package dto

import (
	"net/url"

	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/spf13/cast"
)

var _ = dependency.IBizRequest(&BizRequest{}) // check impl

// BizRequest default business tag reqquest
type BizRequest struct {
	BizId int64 `json:"bizId" bson:"bizId" form:"bizId" uri:"bizId" url:"bizId,omitempty"`
}

func (r *BizRequest) GetBizId() int64 {
	return r.BizId
}

func (r *BizRequest) ToUrlQuery() url.Values {
	u := url.Values{}
	u.Add("bizId", cast.ToString(r.BizId))
	return u
}
