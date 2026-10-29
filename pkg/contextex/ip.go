package contextex

import (
	"context"

	coreContextex "github.com/illidaris/core/contextex"

	"github.com/illidaris/aphrodite/pkg/netex"
)

// WithIP stores the client IP address in ctx.
func WithIP(ctx context.Context, v string) context.Context {
	return coreContextex.WithIP(ctx, v)
}

// GetIP returns the client IP address stored in ctx.
func GetIP(ctx context.Context) string {
	return coreContextex.GetIP(ctx)
}

// GetIPInt parses the client IP as an IPv4 integer. It returns 0 for invalid
// addresses and IPv6 addresses.
func GetIPInt(ctx context.Context) uint32 {
	v, err := netex.IPv4ToInt(GetIP(ctx))
	if err != nil {
		return 0
	}
	return v
}
