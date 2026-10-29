package contextex

import (
	"context"
	"testing"
)

func TestLogIgnore(t *testing.T) {
	ctx := context.Background()
	if IsLogIgnoreFrmCtx(ctx) {
		t.Fatal("new context unexpectedly ignores logs")
	}
	if !IsLogIgnoreFrmCtx(WithLogIgnore(ctx)) {
		t.Fatal("WithLogIgnore did not mark the context")
	}
}

func TestIPHelpers(t *testing.T) {
	ctx := WithIP(context.Background(), "192.168.1.10")
	if got := GetIP(ctx); got != "192.168.1.10" {
		t.Errorf("GetIP() = %q", got)
	}
	if got := GetIPInt(ctx); got != 3232235786 {
		t.Errorf("GetIPInt() = %d, want 3232235786", got)
	}
	if got := GetIPInt(WithIP(context.Background(), "not-an-ip")); got != 0 {
		t.Errorf("GetIPInt(invalid) = %d, want 0", got)
	}
	if got := GetIPInt(WithIP(context.Background(), "2001:db8::1")); got != 0 {
		t.Errorf("GetIPInt(IPv6) = %d, want 0", got)
	}
}

func TestBizIDHelpers(t *testing.T) {
	ctx := WithBizId(context.Background(), 42)
	if got := GetBizId(ctx); got != 42 {
		t.Errorf("GetBizId() = %d, want 42", got)
	}
	if got := GetBizId(context.Background()); got != 0 {
		t.Errorf("GetBizId(empty) = %d, want 0", got)
	}
}
