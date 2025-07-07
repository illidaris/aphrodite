package ginhandle

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/illidaris/aphrodite/dto"
	"github.com/illidaris/aphrodite/pkg/exception"
)

// GinExHandler 通用调用处理
func GinExHandler[Req, Resp any](request *Req, exec func(context.Context, *Req) (Resp, exception.Exception), reqFuncs []func(context.Context, *Req)) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		if exec == nil {
			c.AbortWithStatusJSON(http.StatusOK, dto.NewResponse(nil, exception.ERR_BUSI.New("当前业务尚未启用")))
			return
		}
		if request != nil {
			if err := c.ShouldBind(request); err != nil {
				c.AbortWithStatusJSON(http.StatusOK, dto.NewResponse(nil, exception.ERR_COMMON_BADPARAM.Wrap(err)))
				return
			}
			if err := c.ShouldBindUri(request); err != nil {
				c.AbortWithStatusJSON(http.StatusOK, dto.NewResponse(nil, exception.ERR_COMMON_BADPARAM.Wrap(err)))
				return
			}
		}
		for _, f := range reqFuncs {
			f(ctx, request)
		}
		res, ex := exec(ctx, request)
		c.JSON(http.StatusOK, dto.NewResponse(res, ex))
	}
}

// Deprecated: 弃用, 请使用 GinExHandler
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
