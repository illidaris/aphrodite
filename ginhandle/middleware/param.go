package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/illidaris/logger"
)

type ParamMiddlewareOption func(opt *paramMiddlewareOptions)

type paramMiddlewareOptions struct {
	RequestContentLengthMax  uint64
	ResponseContentLengthMax uint64
	AfterFunc                func(ctx context.Context, info *APIInfo)
}

func WithRequestContentLengthMax(max uint64) ParamMiddlewareOption {
	return func(opt *paramMiddlewareOptions) {
		opt.RequestContentLengthMax = max
	}
}

func WithResponseContentLengthMax(max uint64) ParamMiddlewareOption {
	return func(opt *paramMiddlewareOptions) {
		opt.ResponseContentLengthMax = max
	}
}

func WithAfterFunc(f func(ctx context.Context, info *APIInfo)) ParamMiddlewareOption {
	return func(opt *paramMiddlewareOptions) {
		opt.AfterFunc = f
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// ParamMiddleware 出参入参记录中间件
func ParamMiddleware(opts ...ParamMiddlewareOption) gin.HandlerFunc {
	opt := &paramMiddlewareOptions{}
	for _, f := range opts {
		f(opt)
	}
	return func(c *gin.Context) {
		now := time.Now()
		// max > 0  enable log response
		if opt.RequestContentLengthMax > 0 {
			if c.Request.ContentLength < int64(opt.RequestContentLengthMax) && c.ContentType() != binding.MIMEMultipartPOSTForm {
				bodyBytes, _ := io.ReadAll(c.Request.Body)
				if len(bodyBytes) > 0 {
					logger.InfoCtx(c.Request.Context(), "[Request]"+string(bodyBytes))
				}
				_ = c.Request.Body.Close() //  must close
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			} else {
				logger.InfoCtx(c.Request.Context(), fmt.Sprintf("request %d is too long", c.Request.ContentLength))
			}
		}
		// max > 0  enable log response
		if opt.ResponseContentLengthMax > 0 {
			w := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = w
			c.Next()
			l := w.body.Len()
			if l < int(opt.ResponseContentLengthMax) {
				responseBody := w.body.String()
				logger.InfoCtx(c.Request.Context(), "[Response]"+responseBody)
			} else {
				logger.InfoCtx(c.Request.Context(), fmt.Sprintf("response %d is too long", l))
			}
		} else {
			c.Next()
		}
		cost := time.Since(now)
		if opt.AfterFunc != nil {
			opt.AfterFunc(c.Request.Context(), &APIInfo{
				Method: c.Request.Method,
				Path:   c.Request.URL.Path,
				Cost:   cost.Milliseconds(),
				Status: c.Writer.Status(),
				Query:  c.Request.URL.RawQuery,
				IP:     c.ClientIP(),
			})
		}
	}
}
