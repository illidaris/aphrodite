package convert

import (
	"testing"
	"unicode/utf8"

	"github.com/smartystreets/goconvey/convey"
)

// TestJson
func TestJson(t *testing.T) {
	convey.Convey("TestJson", t, func() {
		m := map[string]interface{}{"a": "x", "b": 1}
		ft := "{\"a\":\"x\",\"b\":1}"
		convey.Convey("empty", func() {
			raw := Json(m)
			convey.So(raw, convey.ShouldEqual, ft)
		})
	})
}

func TestTruncateMySQLVarchar(t *testing.T) {
	convey.Convey("TestTruncateMySQLVarchar", t, func() {
		tests := []struct {
			name   string
			raw    string
			length int
			want   string
		}{
			{
				name:   "ascii shorter than length",
				raw:    "abc123",
				length: 10,
				want:   "abc123",
			},
			{
				name:   "ascii truncated",
				raw:    "abc123",
				length: 3,
				want:   "abc",
			},
			{
				name:   "chinese truncated by utf8mb4 bytes",
				raw:    "中文测试abc",
				length: 10,
				want:   "中文测",
			},
			{
				name:   "mixed digits letters symbols and emoji",
				raw:    "中A1!🙂文B2@🚀",
				length: 13,
				want:   "中A1!🙂文",
			},
			{
				name:   "emoji boundary with four bytes",
				raw:    "🙂🚀✨abc",
				length: 8,
				want:   "🙂🚀",
			},
			{
				name:   "does not split emoji when byte budget is too small",
				raw:    "a🙂b",
				length: 4,
				want:   "a",
			},
			{
				name:   "zero length",
				raw:    "abc",
				length: 0,
				want:   "",
			},
			{
				name:   "negative length",
				raw:    "abc",
				length: -1,
				want:   "",
			},
			{
				name:   "empty string",
				raw:    "",
				length: 10,
				want:   "",
			},
		}

		for _, tt := range tests {
			tt := tt
			convey.Convey(tt.name, func() {
				got := TruncateMySQLVarchar(tt.raw, tt.length)
				convey.So(got, convey.ShouldEqual, tt.want)
				convey.So(len(got) <= max(tt.length, 0), convey.ShouldBeTrue)
				convey.So(utf8.ValidString(got), convey.ShouldBeTrue)
			})
		}
	})
}
