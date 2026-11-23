package check

import "testing"

func TestIsIPInWhitelist(t *testing.T) {
	whitelist := []string{
		" 10.0.0.1 ",
		"192.168.1.0/24",
		"172.16.10.10-172.16.10.20",
		"2001:db8::/32",
		"2001:db8:1::1-2001:db8:1::ff",
		"invalid",
	}

	tests := []struct {
		name string
		ip   string
		want bool
	}{
		// 独立 IP 精确命中。
		{"standalone IPv4", "10.0.0.1", true},
		// 输入 IP 的首尾空格应被忽略。
		{"standalone IPv4 with spaces", " 10.0.0.1 ", true},
		// CIDR 号段的网络起始地址应命中。
		{"CIDR IPv4 lower bound", "192.168.1.0", true},
		// CIDR 号段的最后一个地址应命中。
		{"CIDR IPv4 upper bound", "192.168.1.255", true},
		// 不属于 CIDR 号段的地址不应命中。
		{"CIDR IPv4 outside", "192.168.2.1", false},
		// 起止范围的起始地址应命中。
		{"range IPv4 lower bound", "172.16.10.10", true},
		// 起止范围的结束地址应命中。
		{"range IPv4 upper bound", "172.16.10.20", true},
		// 小于起止范围起始地址的值不应命中。
		{"range IPv4 below", "172.16.10.9", false},
		// 大于起止范围结束地址的值不应命中。
		{"range IPv4 above", "172.16.10.21", false},
		// IPv6 CIDR 号段内的地址应命中。
		{"CIDR IPv6", "2001:db8:ffff::1", true},
		// IPv6 起止范围内的地址应命中。
		{"range IPv6", "2001:db8:1::80", true},
		// 不属于任何 IPv6 白名单范围的地址不应命中。
		{"IPv6 outside", "2001:db9::1", false},
		// 非法 IP 输入应直接返回 false。
		{"invalid IP", "not-an-ip", false},
		// 空 IP 输入应直接返回 false。
		{"empty IP", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIPInWhitelist(tt.ip, whitelist); got != tt.want {
				t.Errorf("IsIPInWhitelist(%q, whitelist) = %v, want %v", tt.ip, got, tt.want)
			}
		})
	}
}

func TestIsIPInWhitelistInvalidEntriesAndEmptyList(t *testing.T) {
	whitelist := []string{
		"",                          // 空项应被忽略。
		"192.168.1.0/33",            // CIDR 前缀长度越界。
		"192.168.1.20-192.168.1.10", // 范围起始地址大于结束地址。
		"192.168.1.1-not-an-ip",     // 范围结束地址非法。
	}
	// 白名单中的无效项不能被误判为命中。
	if IsIPInWhitelist("192.168.1.1", whitelist) {
		t.Fatal("invalid whitelist entries should not match")
	}
	// 空白名单不能命中任何 IP。
	if IsIPInWhitelist("192.168.1.1", nil) {
		t.Fatal("empty whitelist should not match")
	}
}

var benchmarkIPInWhitelistResult bool

// BenchmarkIsIPInWhitelist 分别测量独立 IP、CIDR、起止范围和未命中场景的性能。
func BenchmarkIsIPInWhitelist(b *testing.B) {
	benchmarks := []struct {
		name      string
		ip        string
		whitelist []string
	}{
		// 白名单中仅有一个独立 IP，且首项即命中。
		{"StandaloneIPv4", "10.0.0.1", []string{"10.0.0.1"}},
		// 白名单中仅有一个 IPv4 CIDR 号段。
		{"CIDRIPv4", "192.168.1.100", []string{"192.168.1.0/24"}},
		// 白名单中仅有一个 IPv4 起止范围。
		{"RangeIPv4", "172.16.10.15", []string{"172.16.10.10-172.16.10.20"}},
		// 遍历完整白名单后仍未命中，用于测量最坏匹配路径。
		{"Miss", "10.255.255.255", []string{
			"10.0.0.1",
			"192.168.1.0/24",
			"172.16.10.10-172.16.10.20",
			"2001:db8::/32",
		}},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				benchmarkIPInWhitelistResult = IsIPInWhitelist(bm.ip, bm.whitelist)
			}
		})
	}
}
