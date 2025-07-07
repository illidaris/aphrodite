package dto

import (
	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/aphrodite/pkg/netex"
)

var _ = dependency.IIP(&IPRequest{})

// BizRequest default business tag reqquest
type IPRequest struct {
	IP string `json:"ip" bson:"ip,omitempty" form:"ip" uri:"ip" url:"ip,omitempty"`
}

func (r *IPRequest) SetIP(v string) {
	r.IP = v
}

func (r IPRequest) GetIP() string {
	return r.IP
}

func (r IPRequest) GetIPInt() uint32 {
	if r.IP == "" {
		return 0
	}
	ip, _ := netex.IPv4ToInt(r.IP)
	return ip
}
