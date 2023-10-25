package dependency

import "net/url"

// IBizRequest business tag reqquest
type IBizRequest interface {
	GetBizId() int64
	ToUrlQuery() url.Values
}
