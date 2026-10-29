package contextex

import (
	"context"
	"testing"

	"github.com/illidaris/core"
	coreContextex "github.com/illidaris/core/contextex"
)

func TestTransferContext(t *testing.T) {
	type key string
	const selected key = "selected"
	const ignored key = "ignored"

	src := context.WithValue(context.Background(), selected, "copied")
	src = context.WithValue(src, ignored, "not copied")
	src = core.TraceID.Set(src, "trace-1")
	to := context.WithValue(context.Background(), ignored, "target value")

	got := TransferContext(src, to, selected, ignored)
	if got.Value(selected) != "copied" {
		t.Fatalf("selected value = %v, want copied", got.Value(selected))
	}
	if core.TraceID.Get(got) != "trace-1" {
		t.Errorf("trace ID = %q, want trace-1", core.TraceID.Get(got))
	}
	if got.Value(key("target-only")) != nil {
		t.Errorf("unexpected target-only value: %v", got.Value(key("target-only")))
	}
	if got.Value(ignored) != "not copied" {
		t.Errorf("source selected value should override target: %v", got.Value(ignored))
	}
}

func TestTransferContextPreservesTargetContext(t *testing.T) {
	type key string
	src := core.TraceID.Set(context.Background(), "source-trace")
	target, cancel := context.WithCancel(context.Background())
	defer cancel()
	target = context.WithValue(target, key("target"), true)

	got := TransferContext(src, target)
	if got.Value(key("target")) != true {
		t.Fatal("target value was not preserved")
	}
	cancel()
	if got.Err() != context.Canceled {
		t.Errorf("target cancellation was not preserved: Err() = %v", got.Err())
	}
}

func TestTransferContextPreservesTargetTraceWhenSourceHasNone(t *testing.T) {
	target := core.TraceID.Set(context.Background(), "target-trace")
	got := TransferContext(context.Background(), target)
	if core.TraceID.Get(got) != "target-trace" {
		t.Errorf("trace ID = %v, want target-trace", core.TraceID.Get(got))
	}
}

func TestTransferBackgroundDetachesCancellation(t *testing.T) {
	type key string
	const selected key = "selected"
	source, cancel := context.WithCancel(context.Background())
	source = coreContextex.WithTraceId(source, "trace-2")
	source = context.WithValue(source, selected, "copied")
	got := TransferBackground(source, selected)
	cancel()

	if got.Err() != nil {
		t.Errorf("background context inherited cancellation: %v", got.Err())
	}
	if got.Value(selected) != "copied" {
		t.Errorf("selected value = %v, want copied", got.Value(selected))
	}
	if coreContextex.TraceIdFrmCtx(got) != "trace-2" {
		t.Errorf("trace ID = %q, want trace-2", coreContextex.TraceIdFrmCtx(got))
	}
	if core.TraceID.Get(got) != "trace-2" {
		t.Errorf("core trace ID = %q, want trace-2", core.TraceID.Get(got))
	}
	if sessionID := coreContextex.SessionIdFrmCtx(got); sessionID == "" {
		t.Error("session ID is empty")
	}
}

func TestTransferBackgroundUsesCoreTraceID(t *testing.T) {
	source := core.TraceID.Set(context.Background(), "core-trace")
	got := TransferBackground(source)
	if coreContextex.TraceIdFrmCtx(got) != "core-trace" {
		t.Errorf("trace ID = %q, want core-trace", coreContextex.TraceIdFrmCtx(got))
	}
}

func TestTransferBackgroundWithoutTraceUsesSessionID(t *testing.T) {
	got := TransferBackground(context.Background())
	traceID := coreContextex.TraceIdFrmCtx(got)
	sessionID := coreContextex.SessionIdFrmCtx(got)
	if traceID == "" {
		t.Fatal("generated trace ID is empty")
	}
	if sessionID == "" {
		t.Fatal("generated session ID is empty")
	}
	if traceID != sessionID {
		t.Errorf("fallback trace ID = %q, want session ID %q", traceID, sessionID)
	}
}

func TestTransferBackgroundDoesNotInheritDeadline(t *testing.T) {
	type key string
	const selected key = "selected"
	source, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	source = context.WithValue(source, selected, "copied")
	got := TransferBackground(source, selected)

	if _, ok := got.Deadline(); ok {
		t.Error("background context unexpectedly has a deadline")
	}
	if got.Err() != nil {
		t.Errorf("background context inherited source error: %v", got.Err())
	}
	if got.Value(selected) != "copied" {
		t.Errorf("selected value = %v, want copied", got.Value(selected))
	}
}
