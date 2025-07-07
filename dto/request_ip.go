package dto

import "github.com/illidaris/aphrodite/pkg/dependency"

var _ = dependency.IIP(&IPRequest{})

// BizRequest default business tag reqquest
type IPRequest struct {
	IP string `json:"ip" bson:"ip,omitempty" form:"ip" uri:"ip" url:"ip,omitempty"`
}

func (r *IPRequest) SetIP(v string) {
	r.IP = v
}
