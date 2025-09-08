package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/illidaris/core"
	"github.com/illidaris/logger"
	"go.uber.org/zap"
)

// LoggerHandler record log
func LoggerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// content type
		contentType := c.GetHeader("Content-Type")
		sessionBirth := time.Now()
		// trace
		WithTrace(c, sessionBirth)
		// before
		c.Next()
		// after
		cost := time.Since(sessionBirth)
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		curCtx := c.Request.Context()
		curCtx = logger.NewContext(curCtx,
			zap.String(string(HTTPMethod), c.Request.Method),
			zap.String(string(HTTPContentType), contentType),
			zap.String(string(HTTPPath), path),
			zap.String(string(HTTPQuery), query),
			zap.String(string(HTTPClientIP), c.ClientIP()),
			zap.String(string(HTTPUserAgent), c.Request.UserAgent()),
			zap.Int(string(HTTPStatusCode), c.Writer.Status()),
			zap.Int64(core.Duration.String(), cost.Milliseconds()),
		)
		if ginErrs := c.Errors.ByType(gin.ErrorTypePrivate).String(); ginErrs != "" {
			logger.WarnCtx(curCtx, ginErrs)
		}
	}
}

type HTTPMetaData core.MetaData

const (
	HTTPStatusCode  HTTPMetaData = "statusCode"
	HTTPContentType HTTPMetaData = "contentType"
	HTTPMethod      HTTPMetaData = "httpMethod"
	HTTPPath        HTTPMetaData = "httpPath"
	HTTPQuery       HTTPMetaData = "httpQuery"
	HTTPClientIP    HTTPMetaData = "httpClientIp"
	HTTPUserAgent   HTTPMetaData = "httpUserAgent"
)

// WithTrace add trace log context
func WithTrace(c *gin.Context, birth time.Time) *gin.Context {
	// trace
	const rid = "X-Request-ID"
	traceID := c.GetHeader(rid)
	// session
	newUUID, _ := uuid.NewUUID()
	sID := newUUID.String()
	if traceID == "" {
		traceID = sID
	}
	sessionBirth := birth.UTC().UnixNano()
	// assembly trace & session
	ctx := c.Request.Context()
	ctx = core.TraceID.SetString(ctx, traceID) // set traceid  into ctx
	ctx = core.SessionID.SetString(ctx, sID)   // set session  into ctx
	ctx = logger.NewContext(ctx,
		zap.String(core.Action.String(), c.Request.URL.Path),
		zap.String(core.TraceID.String(), traceID),
		zap.String(core.SessionID.String(), sID),
		zap.Int64(core.SessionBirth.String(), sessionBirth))
	// instead of request
	c.Request = c.Request.WithContext(ctx)
	return c
}
