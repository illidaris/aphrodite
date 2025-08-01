package api

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
)

var ContextLength = Uint("OLLAMA_CONTEXT_LENGTH", 4096)

func Uint(key string, defaultValue uint) func() uint {
	return func() uint {
		if s := Var(key); s != "" {
			if n, err := strconv.ParseUint(s, 10, 64); err != nil {
				slog.Warn("invalid environment variable, using default", "key", key, "value", s, "default", defaultValue)
			} else {
				return uint(n)
			}
		}

		return defaultValue
	}
}

// Var returns an environment variable stripped of leading and trailing quotes or spaces
func Var(key string) string {
	return strings.Trim(strings.TrimSpace(os.Getenv(key)), "\"'")
}
