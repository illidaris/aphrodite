package check

import (
	"net/netip"
	"strings"
)

// IsIPInWhitelist 判断 ip 是否命中 IP 白名单。
// 白名单项支持独立 IPv4/IPv6 地址、CIDR 号段（如 192.168.1.0/24）
// 以及起止地址范围（如 192.168.1.10-192.168.1.20）。
// 无效的 ip 或白名单项会被忽略。
func IsIPInWhitelist(ip string, whitelist []string) bool {
	target, err := netip.ParseAddr(strings.TrimSpace(ip))
	if err != nil {
		return false
	}

	for _, item := range whitelist {
		if isIPInWhitelistEntry(target, strings.TrimSpace(item)) {
			return true
		}
	}
	return false
}

func isIPInWhitelistEntry(ip netip.Addr, entry string) bool {
	if entry == "" {
		return false
	}

	// 独立 IP
	if allowed, err := netip.ParseAddr(entry); err == nil {
		return ip == allowed
	}

	// CIDR 号段
	if prefix, err := netip.ParsePrefix(entry); err == nil {
		return prefix.Contains(ip)
	}

	// 起止 IP 范围
	startText, endText, ok := strings.Cut(entry, "-")
	if !ok {
		return false
	}
	start, err := netip.ParseAddr(strings.TrimSpace(startText))
	if err != nil {
		return false
	}
	end, err := netip.ParseAddr(strings.TrimSpace(endText))
	if err != nil || start.BitLen() != end.BitLen() || start.Compare(end) > 0 {
		return false
	}
	return start.Compare(ip) <= 0 && ip.Compare(end) <= 0
}
