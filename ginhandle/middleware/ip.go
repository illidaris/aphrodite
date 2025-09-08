package middleware

import (
	"github.com/gin-gonic/gin"
	acContextex "github.com/illidaris/aphrodite/pkg/contextex"
)

func APIMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ip := c.ClientIP()
		if ip != "" {
			ctx = acContextex.WithIP(ctx, ip)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}
