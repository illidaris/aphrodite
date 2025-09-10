package ginhandle

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/illidaris/aphrodite/ginhandle/middleware"
	prometheusEx "github.com/illidaris/aphrodite/ginhandle/middleware/prometheus"
	"github.com/illidaris/logger"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-gonic/gin"
	_ "go.uber.org/automaxprocs"
)

// GinHandleOptions 运行参数
type GinHandleOptions struct {
	Mode                string
	Collectors          []prometheus.Collector
	ParamMiddleware     bool
	ParamMiddlewareOpts []middleware.ParamMiddlewareOption
	HealthCheck         bool
	MetricCheck         bool
	PrometheusConfig    prometheusEx.Config
	GinInnerHandle      func(*gin.Engine)
}

// Gin模式
func WithMode(mode string) GinHandleOption {
	return func(opts *GinHandleOptions) {
		opts.Mode = mode
	}
}

// 启用健康检查
func WithHealthCheck(v bool) GinHandleOption {
	return func(opts *GinHandleOptions) {
		opts.HealthCheck = v
	}
}

// 启用指标检查
func WithMetricCheck(v bool) GinHandleOption {
	return func(opts *GinHandleOptions) {
		opts.MetricCheck = v
	}
}

// Prometheus收集器
func WithCollectors(cs ...prometheus.Collector) GinHandleOption {
	return func(opts *GinHandleOptions) {
		opts.Collectors = append(opts.Collectors, cs...)
	}
}

// 启用 prometheus
func WithPrometheusConfig(v prometheusEx.Config) GinHandleOption {
	return func(opts *GinHandleOptions) {
		opts.PrometheusConfig = v
	}
}

// 内置 handle
func WithGinInnerHandle(f func(*gin.Engine)) GinHandleOption {
	return func(opts *GinHandleOptions) {
		opts.GinInnerHandle = f
	}
}

// 参数中间件
func WithParamMiddleware(v bool, ps ...middleware.ParamMiddlewareOption) GinHandleOption {
	return func(opts *GinHandleOptions) {
		opts.ParamMiddleware = v
		opts.ParamMiddlewareOpts = append(opts.ParamMiddlewareOpts, ps...)
	}
}

type GinHandleOption func(*GinHandleOptions)

func NewGinHandleOption(opts ...GinHandleOption) *GinHandleOptions {
	o := &GinHandleOptions{
		Mode:             gin.ReleaseMode,
		Collectors:       []prometheus.Collector{},
		ParamMiddleware:  true,
		HealthCheck:      true,
		MetricCheck:      true,
		PrometheusConfig: prometheusEx.Config{},
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// Gin默认初始化
func NewGin(opts ...GinHandleOption) *gin.Engine {
	opt := NewGinHandleOption(opts...)
	gin.SetMode(opt.Mode)
	engine := gin.New()
	engine.Use(
		middleware.LoggerHandler(),
		middleware.RecoverHandler(),
		middleware.APIMiddleware(),
	)
	if opt.HealthCheck {
		engine.HEAD("/health", func(c *gin.Context) { c.Status(http.StatusOK) })
		engine.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })
	}
	if opt.GinInnerHandle != nil {
		opt.GinInnerHandle(engine)
	}
	if opt.MetricCheck {
		reg := prometheus.NewRegistry()
		prometheus.DefaultRegisterer = reg
		prometheus.DefaultGatherer = reg
		if len(opt.Collectors) > 0 {
			reg.MustRegister(opt.Collectors...)
		}
		p := prometheusEx.NewWithConfig(opt.PrometheusConfig)
		p.Use(engine)
	}
	if opt.ParamMiddleware {
		engine.Use(middleware.ParamMiddleware(opt.ParamMiddlewareOpts...))
	}
	return engine
}

func GracefulRun(ctx context.Context, e http.Handler, addr string, timeout time.Duration) {
	// bind ip&port
	srv := &http.Server{
		Addr:    addr,
		Handler: e,
	}

	errCh := make(chan error, 1)
	defer close(errCh)
	// listen
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	notifyCtx, stop := signal.NotifyContext(ctx, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case s := <-notifyCtx.Done():
		stop()
		logger.Info(fmt.Sprintf("Shutdown: Receive Sign(%s)", s))
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		if err := srv.Shutdown(timeoutCtx); err != nil {
			logger.Error(fmt.Sprintf("Shutdown: %s", err))
		}
		logger.Info("Shutdown: exit")
		break
	case err := <-errCh:
		logger.Error(fmt.Sprintf("Listen: Receive Error %s", err))
		break
	}
}

func GracefulRunWithAop(ctx context.Context, e http.Handler, addr string, timeout time.Duration, before func(port int), after func()) {
	// bind ip&port
	srv := &http.Server{
		Handler: e,
	}
	BaseGracefulRunWithAop(ctx, srv, addr, timeout, before, after)
}

func BaseGracefulRunWithAop(ctx context.Context, srv *http.Server, addr string, timeout time.Duration, before func(port int), after func()) {
	defer after()
	errCh := make(chan error, 1)
	defer close(errCh)
	// listen
	go func() {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			errCh <- err
			return
		}
		port := ln.Addr().(*net.TCPAddr).Port
		before(port)
		// service connections
		if err := srv.Serve(ln); err != nil {
			errCh <- err
		}
	}()

	notifyCtx, stop := signal.NotifyContext(ctx, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case s := <-notifyCtx.Done():
		stop()
		logger.Info(fmt.Sprintf("Shutdown: Receive Sign(%s)", s))
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		if err := srv.Shutdown(timeoutCtx); err != nil {
			logger.Error(fmt.Sprintf("Shutdown: %s", err))
		}
		logger.Info("Shutdown: exit")
		break
	case err := <-errCh:
		logger.Error(fmt.Sprintf("Listen: Receive Error %s", err))
		break
	}
}
