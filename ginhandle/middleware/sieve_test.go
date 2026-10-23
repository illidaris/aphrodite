package middleware

import (
	"context"
	"testing"
)

func TestSieveConfigIsInWhite(t *testing.T) {
	tests := []struct {
		name         string
		cfg          SieveConfig
		id, ip, head string
		want         bool
	}{
		{"or matches one", SieveConfig{CondType: WhiteCondTypeOr, WhiteIds: []string{"u1"}, WhiteIps: []string{"10.0.0.1"}}, "other", "10.0.0.1", "", true},
		{"or matches none", SieveConfig{CondType: WhiteCondTypeOr, WhiteIds: []string{"u1"}, WhiteIps: []string{"10.0.0.1"}}, "other", "10.0.0.2", "", false},
		{"and matches all", SieveConfig{CondType: WhiteCondTypeAnd, WhiteIds: []string{"u1"}, WhiteIps: []string{"10.0.0.1"}}, "u1", "10.0.0.1", "", true},
		{"and misses one", SieveConfig{CondType: WhiteCondTypeAnd, WhiteIds: []string{"u1"}, WhiteIps: []string{"10.0.0.1"}}, "u2", "10.0.0.1", "", false},
		{"empty lists allow", SieveConfig{CondType: WhiteCondTypeOr}, "u1", "10.0.0.1", "x", true},
		{"unknown strategy denies", SieveConfig{CondType: WhiteCondType(99), WhiteIds: []string{"u1"}}, "u1", "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.IsInWhite(tt.id, tt.ip, tt.head); got != tt.want {
				t.Fatalf("IsInWhite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSieveConfigWhitelistHelpers(t *testing.T) {
	cfg := SieveConfig{WhiteIds: []string{"u1"}, WhiteIps: []string{"127.0.0.1"}, WhiteHeaders: []string{"token"}}
	tests := map[string]struct {
		check                func() (bool, bool)
		wantCheck, wantMatch bool
	}{
		"id":        {func() (bool, bool) { return cfg.IsInWhiteId("u1") }, true, true},
		"ip miss":   {func() (bool, bool) { return cfg.IsInWhiteIp("x") }, true, false},
		"header":    {func() (bool, bool) { return cfg.IsInWhiteHeader("token") }, true, true},
		"empty cfg": {func() (bool, bool) { return (SieveConfig{}).IsInWhiteId("u1") }, false, false},
	}
	for name, tt := range tests {
		gotCheck, gotMatch := tt.check()
		if gotCheck != tt.wantCheck || gotMatch != tt.wantMatch {
			t.Errorf("%s = [%v %v], want [%v %v]", name, gotCheck, gotMatch, tt.wantCheck, tt.wantMatch)
		}
	}
}

func TestSieveWhiteIDContext(t *testing.T) {
	ctx := WithSieveWhiteId(context.Background(), "u1")
	if got := SieveWhiteIdFrmCtx(ctx); got != "u1" {
		t.Fatalf("SieveWhiteIdFrmCtx() = %q, want %q", got, "u1")
	}
	if got := SieveWhiteIdFrmCtx(context.Background()); got != "" {
		t.Fatalf("missing value = %q, want empty", got)
	}
}
