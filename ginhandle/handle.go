package ginhandle

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/illidaris/aphrodite/dto"
	"github.com/illidaris/aphrodite/ginhandle/middleware"
	"github.com/illidaris/aphrodite/ginhandle/middleware/prometheus"
	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/aphrodite/pkg/exception"
)

// BizGinExHandler 通用调用处理
func BizGinExHandler[Req dependency.IBindRequest, Resp any](request Req, exec func(context.Context, Req) (Resp, exception.Exception)) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		if exec == nil {
			middleware.AbortWithException(c, prometheus.BIZ_CODE_ILLEGAL, exception.ERR_BUSI.New("当前业务尚未启用"))
			return
		}
		if any(request) != nil {
			if ex := bindRequest(c, request); ex != nil {
				middleware.AbortWithException(c, prometheus.BIZ_CODE_BADPARAM, ex)
				return
			}
		}
		dependency.BizFrmCtx(ctx, request)
		dependency.IPFrmCtx(ctx, request)
		res, ex := exec(ctx, request)
		prometheus.WithMetricsBizCode(c, prometheus.BIZ_CODE_BUSI, ex)
		c.JSON(http.StatusOK, dto.NewResponse(res, ex))
	}
}

// GinOneHandler 通用调用处理
func GinOneHandler[Req, Resp any](exec func(context.Context, *Req) (Resp, exception.Exception)) func(c *gin.Context) {
	return func(c *gin.Context) {
		request := new(Req)
		ctx := c.Request.Context()
		if exec == nil {
			middleware.AbortWithException(c, prometheus.BIZ_CODE_ILLEGAL, exception.ERR_BUSI.New("当前业务尚未启用"))
			return
		}
		if ex := bindRequest(c, request); ex != nil {
			middleware.AbortWithException(c, prometheus.BIZ_CODE_BADPARAM, ex)
			return
		}
		res, ex := exec(ctx, request)
		prometheus.WithMetricsBizCode(c, prometheus.BIZ_CODE_BUSI, ex)
		c.JSON(http.StatusOK, dto.NewResponse(res, ex))
	}
}

// GinExHandler 通用调用处理
func GinExHandler[Req, Resp any](request *Req, exec func(context.Context, *Req) (Resp, exception.Exception), reqFuncs []func(context.Context, *Req) exception.Exception) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		if exec == nil {
			middleware.AbortWithException(c, prometheus.BIZ_CODE_ILLEGAL, exception.ERR_BUSI.New("当前业务尚未启用"))
			return
		}
		if request != nil {
			if ex := bindRequest(c, request); ex != nil {
				middleware.AbortWithException(c, prometheus.BIZ_CODE_BADPARAM, ex)
				return
			}
		}
		for _, f := range reqFuncs {
			ex := f(ctx, request)
			if ex != nil {
				middleware.AbortWithException(c, prometheus.BIZ_CODE_BUSI, ex)
				return
			}
		}
		res, ex := exec(ctx, request)
		prometheus.WithMetricsBizCode(c, prometheus.BIZ_CODE_BUSI, ex)
		c.JSON(http.StatusOK, dto.NewResponse(res, ex))
	}
}

// Deprecated: 弃用, 请使用 GinExHandler
func GinHandler[Req, Resp any](request Req, f func(context.Context, Req) (Resp, exception.Exception)) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		if ex := bindRequest(c, request); ex != nil {
			middleware.AbortWithException(c, prometheus.BIZ_CODE_BADPARAM, ex)
			return
		}
		res, ex := f(ctx, request)
		prometheus.WithMetricsBizCode(c, prometheus.BIZ_CODE_BUSI, ex)
		c.JSON(http.StatusOK, dto.NewResponse(res, ex))
	}
}

func bindRequest(c *gin.Context, request any) exception.Exception {
	if err := c.ShouldBind(request); err != nil {
		return exception.ERR_COMMON_BADPARAM.Wrap(err)
	}
	if err := c.ShouldBindUri(request); err != nil {
		return exception.ERR_COMMON_BADPARAM.Wrap(err)
	}
	return nil
}
