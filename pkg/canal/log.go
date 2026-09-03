package canal

import (
	"context"
	"fmt"
)

type ILogger interface {
	Debug(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
}

type DefaultLogger struct{}

func (l DefaultLogger) Debug(ctx context.Context, msg string, args ...any) {
	fmt.Printf(msg+"\n", args...)
}
func (l DefaultLogger) Info(ctx context.Context, msg string, args ...any) {
	fmt.Printf(msg+"\n", args...)
}
func (l DefaultLogger) Warn(ctx context.Context, msg string, args ...any) {
	fmt.Printf(msg+"\n", args...)
}
func (l DefaultLogger) Error(ctx context.Context, msg string, args ...any) {
	fmt.Printf(msg+"\n", args...)
}
