package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/illidaris/aphrodite/pkg/exception"
	"github.com/illidaris/aphrodite/pkg/when"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSieveWhiteIdFrmCtx(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		want string
	}{
		{"missing", context.Background(), ""},
		{"string", WithSieveWhiteId(context.Background(), "user-1"), "user-1"},
		{"wrong type", context.WithValue(context.Background(), busiWhiteIdCtxKey{}, 42), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SieveWhiteIdFrmCtx(tt.ctx); got != tt.want {
				t.Errorf("SieveWhiteIdFrmCtx() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSieveConfigWhiteLists(t *testing.T) {
	cfg := SieveConfig{WhiteIds: []string{"id-1"}, WhiteIps: []string{"192.0.2.1"}, WhiteHeaders: []string{"allowed"}}
	tests := []struct {
		name string
		got  bool
		want bool
	}{
		{"id hit", cfg.IsInWhiteId("id-1"), true},
		{"id miss", cfg.IsInWhiteId("id-2"), false},
		{"ip hit", cfg.IsInWhiteIp("192.0.2.1"), true},
		{"ip miss", cfg.IsInWhiteIp("192.0.2.2"), false},
		{"header hit", cfg.IsInWhiteHeader("allowed"), true},
		{"header miss", cfg.IsInWhiteHeader("denied"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
	if !(SieveConfig{}).IsInWhiteId("any") || !(SieveConfig{}).IsInWhiteIp("any") || !(SieveConfig{}).IsInWhiteHeader("any") {
		t.Fatal("empty white lists must match all values")
	}
}

func TestSieveConfigIsInWhite(t *testing.T) {
	tests := []struct {
		name string
		cfg  SieveConfig
		id   string
		ip   string
		head string
		want bool
	}{
		{"or matches one", SieveConfig{CondType: WhiteCondTypeOr, WhiteIds: []string{"id-1"}, WhiteIps: []string{"ip-1"}, WhiteHeaders: []string{"header-1"}}, "id-1", "ip-2", "header-2", true},
		{"or no match", SieveConfig{CondType: WhiteCondTypeOr, WhiteIds: []string{"id-1"}, WhiteIps: []string{"ip-1"}, WhiteHeaders: []string{"header-1"}}, "id-2", "ip-2", "header-2", false},
		{"and matches all configured", SieveConfig{CondType: WhiteCondTypeAnd, WhiteIds: []string{"id-1"}, WhiteIps: []string{"ip-1"}}, "id-1", "ip-1", "ignored", true},
		{"and misses one", SieveConfig{CondType: WhiteCondTypeAnd, WhiteIds: []string{"id-1"}, WhiteIps: []string{"ip-1"}}, "id-1", "ip-2", "", false},
		{"and no conditions", SieveConfig{CondType: WhiteCondTypeAnd}, "id", "ip", "header", true},
		{"or no conditions", SieveConfig{CondType: WhiteCondTypeOr}, "id", "ip", "header", true},
		{"unknown strategy", SieveConfig{CondType: 99, WhiteIds: []string{"id-1"}}, "id-1", "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.IsInWhite(tt.id, tt.ip, tt.head); got != tt.want {
				t.Errorf("IsInWhite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSieveConfigCheck(t *testing.T) {
	now := time.Now()
	future := int(now.Add(time.Hour).Unix())
	past := int(now.Add(-time.Hour).Unix())
	baseCtx := context.Background()

	tests := []struct {
		name string
		cfg  SieveConfig
		ctx  context.Context
		ip   string
		head string
		want string
	}{
		{"disabled", SieveConfig{}, baseCtx, "", "", ""},
		{"enabled without limits", SieveConfig{Enable: true}, baseCtx, "", "", ""},
		{"before start", SieveConfig{Enable: true, NotBefore: future, WhiteIds: []string{"other"}, WhiteIps: []string{"other"}, WhiteHeaders: []string{"other"}}, baseCtx, "", "", "Forbidden Nbf"},
		{"expired", SieveConfig{Enable: true, Expire: past, WhiteIds: []string{"other"}, WhiteIps: []string{"other"}, WhiteHeaders: []string{"other"}}, baseCtx, "", "", "Forbidden Exp"},
		{"white id bypasses time", SieveConfig{Enable: true, NotBefore: future, WhiteIds: []string{"allowed"}}, WithSieveWhiteId(baseCtx, "allowed"), "", "", ""},
		{"white ip bypasses time", SieveConfig{Enable: true, NotBefore: future, WhiteIps: []string{"192.0.2.1"}}, baseCtx, "192.0.2.1", "", ""},
		{"white header bypasses time", SieveConfig{Enable: true, NotBefore: future, WhiteHeaders: []string{"allowed"}}, baseCtx, "", "allowed", ""},
		{"and strategy requires all configured values", SieveConfig{Enable: true, CondType: WhiteCondTypeAnd, NotBefore: future, WhiteIds: []string{"allowed"}, WhiteIps: []string{"192.0.2.1"}}, WithSieveWhiteId(baseCtx, "allowed"), "192.0.2.1", "", ""},
		{"or strategy matches one value", SieveConfig{Enable: true, CondType: WhiteCondTypeOr, NotBefore: future, WhiteIds: []string{"allowed"}, WhiteIps: []string{"192.0.2.1"}}, WithSieveWhiteId(baseCtx, "allowed"), "", "", ""},
		{"time boundary allowed", SieveConfig{Enable: true, NotBefore: 1, Expire: 1}, when.WithClockOffset(baseCtx, time.Unix(1, 0).Sub(time.Now())), "", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Check(tt.ctx, tt.ip, tt.head)
			if tt.want == "" {
				if err != nil {
					t.Fatalf("Check() error = %v, want nil", err)
				}
				return
			}
			if err == nil || err.Error() != tt.want {
				t.Errorf("Check() error = %v, want %q", err, tt.want)
			}
			if ex, ok := err.(exception.Exception); !ok || ex.Code() != int32(exception.ERR_BUSI) {
				t.Errorf("Check() error = %T, want business exception", err)
			}
		})
	}
}

func TestSieveMiddleware(t *testing.T) {
	now := time.Now()
	future := int(now.Add(time.Hour).Unix())
	past := int(now.Add(-time.Hour).Unix())

	tests := []struct {
		name   string
		cfg    SieveConfig
		id     string
		header string
		status int
	}{
		{"disabled", SieveConfig{}, "", "", http.StatusNoContent},
		{"enabled", SieveConfig{Enable: true}, "", "", http.StatusNoContent},
		{"before start", SieveConfig{Enable: true, NotBefore: future, WhiteIds: []string{"other-id"}, WhiteIps: []string{"203.0.113.1"}, WhiteHeaders: []string{"other-header"}}, "", "", http.StatusForbidden},
		{"expired", SieveConfig{Enable: true, Expire: past, WhiteIds: []string{"other-id"}, WhiteIps: []string{"203.0.113.1"}, WhiteHeaders: []string{"other-header"}}, "", "", http.StatusForbidden},
		{"white id bypasses time", SieveConfig{Enable: true, NotBefore: future, WhiteIds: []string{"allowed"}}, "allowed", "", http.StatusNoContent},
		{"white header bypasses time", SieveConfig{Enable: true, NotBefore: future, WhiteHeaders: []string{"yes"}}, "", "yes", http.StatusNoContent},
		{"or white id bypasses time", SieveConfig{Enable: true, NotBefore: future, WhiteIds: []string{"allowed"}}, "allowed", "", http.StatusNoContent},
		{"and whitelist hit bypasses time", SieveConfig{Enable: true, NotBefore: future, CondType: WhiteCondTypeAnd, WhiteIds: []string{"allowed"}}, "allowed", "", http.StatusNoContent},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := gin.New()
			engine.Use(SieveMiddleware("demo", func(context.Context, string) SieveConfig { return tt.cfg }))
			engine.GET("/", func(c *gin.Context) { c.Status(http.StatusNoContent) })
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = "198.51.100.2:1234"
			ctx := req.Context()
			if tt.id != "" {
				ctx = WithSieveWhiteId(ctx, tt.id)
			}
			req = req.WithContext(ctx)
			if tt.header != "" {
				req.Header.Set(HEAD_WHITE, tt.header)
			}
			rec := httptest.NewRecorder()
			engine.ServeHTTP(rec, req)
			if rec.Code != tt.status {
				t.Errorf("status = %d, want %d; body: %s", rec.Code, tt.status, rec.Body.String())
			}
		})
	}
}

func TestSieveMiddlewareTimeBoundary(t *testing.T) {
	// 时间戳等于 NotBefore/Expire 时均允许通过。
	for _, cfg := range []SieveConfig{{Enable: true, NotBefore: 1}, {Enable: true, Expire: 1}} {
		engine := gin.New()
		engine.Use(SieveMiddleware("demo", func(context.Context, string) SieveConfig { return cfg }))
		engine.GET("/", func(c *gin.Context) { c.Status(http.StatusNoContent) })
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req = req.WithContext(when.WithClockOffset(req.Context(), time.Unix(1, 0).Sub(time.Now())))
		rec := httptest.NewRecorder()
		engine.ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Errorf("boundary status = %d, want %d", rec.Code, http.StatusNoContent)
		}
	}
}
