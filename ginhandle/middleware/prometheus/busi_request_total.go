package prometheus

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/illidaris/aphrodite/pkg/exception"
	"github.com/prometheus/client_golang/prometheus"
)

type BizCode string

const (
	METRICS_BIZ_REQUEST_TOTAL         = "biz_request_total"
	BIZ_CODE_ILLEGAL          BizCode = "Illegal"
	BIZ_CODE_BADPARAM         BizCode = "BadParam"
	BIZ_CODE_BUSI             BizCode = "Busi"
	BIZ_CODE_SUCCESS          BizCode = "Success"
	BIZ_CODE_UNKNOWN          BizCode = "Unknown"
)

type ctxKeyMetricsBizResponseCode struct{}

func WithMetricsBizCode(c *gin.Context, code BizCode, ex exception.Exception) *gin.Context {
	if ex == nil {
		if code == BIZ_CODE_BUSI {
			code = BIZ_CODE_SUCCESS
		} else {
			return c
		}
	}
	c.Request = c.Request.WithContext(WithMetricsBizResponseCode(c.Request.Context(), code))
	return c
}

func WithMetricsBizResponseCode(ctx context.Context, code BizCode) context.Context {
	val := BIZ_CODE_UNKNOWN
	switch code {
	case BIZ_CODE_ILLEGAL,
		BIZ_CODE_BADPARAM,
		BIZ_CODE_BUSI,
		BIZ_CODE_SUCCESS:
		val = code
	}
	return context.WithValue(ctx, ctxKeyMetricsBizResponseCode{}, val)
}

func GetMetricsBizResponseCode(ctx context.Context) (bool, BizCode) {
	v := ctx.Value(ctxKeyMetricsBizResponseCode{})
	if v == nil {
		return false, ""
	}
	val, ok := v.(BizCode)
	if !ok {
		return false, val
	}
	return true, val
}

var (
	BizRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: METRICS_BIZ_REQUEST_TOTAL,
			Help: "Total number of HTTP requests, with biz code",
		},
		[]string{"handler", "host", "method", "url", "biz_code"},
	)
)

func BizRequestTotalPrometheusMiddleware(fn RequestCounterURLLabelMappingFn) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		handler := c.HandlerName() // 获取处理函数名
		host := c.Request.Host
		method := c.Request.Method
		url := c.Request.URL.Path
		if fn != nil {
			url = fn(c)
		}
		exist, bizCode := GetMetricsBizResponseCode(c.Request.Context())
		if !exist {
			return
		}
		// 记录计数
		BizRequestTotal.WithLabelValues(
			handler,
			host,
			method,
			url,
			string(bizCode),
		).Inc()
	}
}
