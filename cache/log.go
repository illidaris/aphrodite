package cache

import (
	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/aphrodite/pkg/logex"
)

// log is an instance of the ILog dependency for logging.
var log dependency.ILog

// SetLog sets the logging dependency.
// Parameters:
// l - represents the logging dependency interface, used to unify the logging method
// This function does not return any value.
func SetLog(l dependency.ILog) {
	log = l // Updates the global log object with the new logging dependency
}

// The logger function provides a reference to a logging instance implementing the ILog interface from the dependency package.
func logger() dependency.ILog {
	// Initializes the logger with the default implementation if it's currently nil.
	if log == nil {
		log = &logex.DefaultLogger{}
	}
	return log
}
