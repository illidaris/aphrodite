package prometheus

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/illidaris/aphrodite/pkg/exception"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	METRICS_BIZ_REQUEST_TOTAL = "biz_request_total"
)

type ctxKeyMetricsBizResponseCode struct{}

func WithMetricsBizResponseCodeFrmEx(c *gin.Context, ex exception.Exception) *gin.Context {
	if ex == nil {
		return c
	}
	c.Request = c.Request.WithContext(WithMetricsBizResponseCode(c.Request.Context(), int64(ex.Code())))
	return c
}

func WithMetricsBizResponseCode(ctx context.Context, code int64) context.Context {
	return context.WithValue(ctx, ctxKeyMetricsBizResponseCode{}, code)
}

func GetMetricsBizResponseCode(ctx context.Context) (bool, int64) {
	v := ctx.Value(ctxKeyMetricsBizResponseCode{})
	if v == nil {
		return false, 0
	}
	val, ok := v.(int64)
	if !ok {
		return false, val
	}
	return true, val
}

var (
	ginRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: METRICS_BIZ_REQUEST_TOTAL,
			Help: "Total number of HTTP requests, with business code",
		},
		[]string{"handler", "host", "method", "url", "business_code"},
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
		ginRequestsTotal.WithLabelValues(
			handler,
			host,
			method,
			url,
			strconv.Itoa(int(bizCode)),
		).Inc()
	}
}
