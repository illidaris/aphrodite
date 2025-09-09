package middleware

import (
	"time"

	libCORS "github.com/gin-contrib/cors"
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

func CorsMiddleware(origins ...string) gin.HandlerFunc {
	return libCORS.New(libCORS.Config{
		AllowOrigins: origins,
		AllowMethods: []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "HEAD", "UPDATE"},
		AllowHeaders: []string{
			"Origin", "OsType", "Accept-Language", "X-Request-ID",
			"Content-Type", "Content-Length", "Accept-Encoding",
			"X-CSRF-Token", "Authorization", "X-Requested-With", "Cache-Control"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	})
}
