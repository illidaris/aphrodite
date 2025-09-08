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

	"github.com/dubbogo/gost/log/logger"
	"github.com/gin-gonic/gin"
)

func NewGin(routeHandle func(*gin.Engine), mode string) *gin.Engine {
	gin.SetMode(mode)
	engine := gin.New()
	engine.Use(
		middleware.LoggerHandler(),
		middleware.RecoverHandler(),
		middleware.APIMiddleware(),
	)
	routeHandle(engine)
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
