package middleware

import (
	"context"
	"testing"

	"github.com/illidaris/aphrodite/pkg/check"
)

func TestSieveConfigIsInWhite(t *testing.T) {
	tests := []struct {
		name         string
		cfg          SieveConfig
		id, ip, head string
		want         bool
	}{
		{
			name: "or matches id",
			cfg: SieveConfig{
				CondType: WhiteCondTypeOr,
				WhiteIds: []string{"u1"},
				WhiteIps: []string{"10.0.0.0/24"},
			},
			id:   "u1",
			ip:   "10.0.1.1",
			want: true,
		},
		{
			name: "or matches ip cidr",
			cfg: SieveConfig{
				CondType: WhiteCondTypeOr,
				WhiteIds: []string{"u1"},
				WhiteIps: []string{"10.0.0.0/24"},
			},
			id:   "u2",
			ip:   "10.0.0.8",
			want: true,
		},
		{
			name: "or matches ip range",
			cfg: SieveConfig{
				CondType: WhiteCondTypeOr,
				WhiteIps: []string{"10.0.0.10-10.0.0.20"},
			},
			ip:   "10.0.0.15",
			want: true,
		},
		{
			name: "or matches header",
			cfg: SieveConfig{
				CondType:     WhiteCondTypeOr,
				WhiteHeaders: []string{"token"},
			},
			head: "token",
			want: true,
		},
		{
			name: "or matches none",
			cfg: SieveConfig{
				CondType: WhiteCondTypeOr,
				WhiteIds: []string{"u1"},
				WhiteIps: []string{"10.0.0.1"},
			},
			id:   "u2",
			ip:   "10.0.0.2",
			want: false,
		},
		{
			name: "or ignores invalid ip",
			cfg: SieveConfig{
				CondType: WhiteCondTypeOr,
				WhiteIps: []string{"10.0.0.1"},
			},
			ip:   "not-an-ip",
			want: false,
		},
		{
			name: "and matches all configured checks",
			cfg: SieveConfig{
				CondType:     WhiteCondTypeAnd,
				WhiteIds:     []string{"u1"},
				WhiteIps:     []string{"10.0.0.1"},
				WhiteHeaders: []string{"token"},
			},
			id:   "u1",
			ip:   "10.0.0.1",
			head: "token",
			want: true,
		},
		{
			name: "and misses one configured check",
			cfg: SieveConfig{
				CondType: WhiteCondTypeAnd,
				WhiteIds: []string{"u1"},
				WhiteIps: []string{"10.0.0.1"},
			},
			id:   "u2",
			ip:   "10.0.0.1",
			want: false,
		},
		{
			name: "and ignores empty whitelist",
			cfg: SieveConfig{
				CondType: WhiteCondTypeAnd,
				WhiteIds: []string{"u1"},
			},
			id:   "u1",
			ip:   "10.0.0.1",
			head: "token",
			want: true,
		},
		{
			name: "no whitelist configured",
			cfg:  SieveConfig{CondType: WhiteCondTypeOr},
			id:   "u1",
			ip:   "10.0.0.1",
			head: "token",
			want: false,
		},
		{
			name: "no whitelist configured with and",
			cfg:  SieveConfig{CondType: WhiteCondTypeAnd},
			id:   "u1",
			ip:   "10.0.0.1",
			head: "token",
			want: false,
		},
		{
			name: "unknown condition type",
			cfg: SieveConfig{
				CondType: WhiteCondType(99),
				WhiteIds: []string{"u1"},
			},
			id:   "u1",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.IsInWhite(tt.id, tt.ip, tt.head); got != tt.want {
				t.Fatalf("IsInWhite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWhiteCheckValid(t *testing.T) {
	tests := []struct {
		name  string
		check WhiteCheck[string]
		want  bool
	}{
		{
			name:  "empty whitelist",
			check: WhiteCheck[string]{Value: "u1"},
			want:  false,
		},
		{
			name:  "default matcher hits",
			check: WhiteCheck[string]{Value: "u1", Whites: []string{"u1", "u2"}},
			want:  true,
		},
		{
			name:  "default matcher misses",
			check: WhiteCheck[string]{Value: "u3", Whites: []string{"u1", "u2"}},
			want:  false,
		},
		{
			name: "custom matcher hits",
			check: WhiteCheck[string]{
				Value:  "10.0.0.8",
				Whites: []string{"10.0.0.0/24"},
				Check: func(whites []string, value string) bool {
					return check.IsIPInWhitelist(value, whites)
				},
			},
			want: true,
		},
		{
			name: "custom matcher misses",
			check: WhiteCheck[string]{
				Value:  "10.0.1.8",
				Whites: []string{"10.0.0.0/24"},
				Check: func(whites []string, value string) bool {
					return check.IsIPInWhitelist(value, whites)
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.check.Valid(); got != tt.want {
				t.Fatalf("Valid() = %v, want %v", got, tt.want)
			}
		})
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

	wrongTypeCtx := context.WithValue(context.Background(), busiWhiteIdCtxKey{}, 1)
	if got := SieveWhiteIdFrmCtx(wrongTypeCtx); got != "" {
		t.Fatalf("wrong type value = %q, want empty", got)
	}
}
