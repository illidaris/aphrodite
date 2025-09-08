package middleware

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/illidaris/core"
	"github.com/illidaris/logger"
	"go.uber.org/zap"
)

// RecoverHandler recover from panic
func RecoverHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				fmt.Println(err)
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					var se *os.SyscallError
					if ok := errors.Is(ne.Err, se); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.WithContext(c.Request.Context()).Error(c.Request.URL.Path, zap.Any(core.Error.String(), err), zap.String("HttpRequest", string(httpRequest)))
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: error check
					c.Abort()
					return
				}

				logger.WithContext(c.Request.Context()).Error("recover from panic",
					zap.Any(core.Error.String(), err),
					zap.String("HttpRequest", string(httpRequest)),
					zap.String("stack", string(debug.Stack())))

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
