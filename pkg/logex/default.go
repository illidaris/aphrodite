package logex

import (
	"context"
	"fmt"

	"github.com/illidaris/aphrodite/pkg/dependency"
)

var _ = dependency.ILog(DefaultLogger{})

// DefaultLogger is a default logging struct.
type DefaultLogger struct{}

// Debug logs debug information.
// Parameters:
//
//	_ context.Context: Contextual information, not used in this implementation.
//	msg string: The message to be printed.
//	args ...interface{}: Variables to be formatted into the message.
func (l DefaultLogger) Debug(_ context.Context, msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}

// Info logs informational messages.
// Parameters are the same as the Debug function.
func (l DefaultLogger) Info(_ context.Context, msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}

// Warn logs warning messages.
// Parameters are the same as the Debug function.
func (l DefaultLogger) Warn(_ context.Context, msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}

// Error logs error messages.
// Parameters are the same as the Debug function.
func (l DefaultLogger) Error(_ context.Context, msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}
