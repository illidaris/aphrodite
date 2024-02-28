package ginhandle

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/illidaris/aphrodite/dto"
	"github.com/illidaris/aphrodite/pkg/exception"
)

// GinHandler 通用调用处理
func GinHandler[Req, Resp any](request Req, f func(context.Context, Req) (Resp, exception.Exception)) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		execFunc := f
		if err := c.ShouldBind(request); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, dto.NewResponse(nil, exception.ERR_COMMON_BADPARAM.Wrap(err)))
			return
		}
		if err := c.ShouldBindUri(request); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, dto.NewResponse(nil, exception.ERR_COMMON_BADPARAM.Wrap(err)))
			return
		}
		res, ex := execFunc(ctx, request)
		c.JSON(http.StatusOK, dto.NewResponse(res, ex))
	}
}
