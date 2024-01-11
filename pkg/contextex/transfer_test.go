package contextex

import (
	"context"
	"testing"

	"github.com/illidaris/core"
	"github.com/smartystreets/goconvey/convey"
)

func TestTransferContext(t *testing.T) {
	convey.Convey("TransferContext", t, func() {
		convey.Convey("TransferContext_Success", func() {
			a := 1
			keyMap := map[any]any{
				"key": "kkk",
				1:     "123",
				&a:    555,
			}
			rawCtx := core.TraceID.Set(context.Background(), "x123x")
			for k, v := range keyMap {
				rawCtx = context.WithValue(rawCtx, k, v)
			}
			keys := []any{}
			for k := range keyMap {
				keys = append(keys, k)
			}
			toCtx := TransferContext(rawCtx, context.Background(), keys...)
			for k, v := range keyMap {
				res := toCtx.Value(k)
				convey.So(res, convey.ShouldEqual, v)
			}
		})
	})
}

func TestTransferBackground(t *testing.T) {
	convey.Convey("TransferBackground", t, func() {
		convey.Convey("TransferBackground_Success", func() {
			a := 1
			keyMap := map[any]any{
				"key": "kkk",
				1:     "123",
				&a:    555,
			}
			rawCtx := core.TraceID.Set(context.Background(), "x123x")
			for k, v := range keyMap {
				rawCtx = context.WithValue(rawCtx, k, v)
			}
			keys := []any{}
			for k := range keyMap {
				keys = append(keys, k)
			}
			toCtx := TransferContext(rawCtx, context.Background(), keys...)
			for k, v := range keyMap {
				res := toCtx.Value(k)
				convey.So(res, convey.ShouldEqual, v)
			}
			convey.So(core.TraceID.Get(toCtx), convey.ShouldEqual, "x123x")
		})
	})
}
