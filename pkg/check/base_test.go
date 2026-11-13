package check

import "testing"

func TestBaseValidators(t *testing.T) {
	tests := []struct {
		name string
		got  bool
		want bool
	}{
		{"IPv4 valid", IsIPv4("192.168.1.1"), true},
		{"IPv4 rejects leading zero", IsIPv4("192.168.001.1"), false},
		{"IPv6 valid", IsIPv6("2001:db8::1"), true},
		{"IPv6 rejects IPv4", IsIPv6("127.0.0.1"), false},
		{"IP valid", IsIP("::1"), true},
		{"IP invalid", IsIP("not an ip"), false},
		{"port lower bound", IsPort("0"), true},
		{"port upper bound", IsPort("65535"), true},
		{"port out of range", IsPort("65536"), false},
		{"port nonnumeric", IsPort("80/tcp"), false},
		{"international phone valid", IsInternationalPhone("+8613800138000"), true},
		{"international phone invalid", IsInternationalPhone("+0123"), false},
		{"gender valid", IsGender("female"), true},
		{"gender invalid", IsGender("Female"), false},
		{"length counts unicode characters", IsStringLengthInRange("你好a", 3, 3), true},
		{"length outside range", IsStringLengthInRange("你好", 3, 4), false},
		{"datetime valid", IsDateTime("2024-02-29", "2006-01-02"), true},
		{"datetime invalid", IsDateTime("2023-02-29", "2006-01-02"), false},
		{"UUID valid", IsUUID("550e8400-e29b-41d4-a716-446655440000"), true},
		{"UUID invalid", IsUUID("550e8400e29b41d4a716446655440000"), false},
		{"JSON valid", IsJSON(`{"ok":true}`), true},
		{"JSON invalid", IsJSON(`{"ok":}`), false},
		{"password strong", IsPasswordStrong("Abcdef1!"), true},
		{"password weak", IsPasswordStrong("abcdefgh"), false},
		{"MAC valid", IsMAC("00:11:22:33:44:55"), true},
		{"MAC invalid", IsMAC("00:11:22:33:44"), false},
		{"ID card valid", IsIDCard("11010519491231002X"), true},
		{"ID card wrong checksum", IsIDCard("110105194912310021"), false},
		{"MD5 valid", IsMD5("d41d8cd98f00b204e9800998ecf8427e"), true},
		{"MD5 invalid", IsMD5("xyz"), false},
		{"SHA1 valid", IsSHA1("0123456789abcdef0123456789abcdef01234567"), true},
		{"SHA1 invalid", IsSHA1("0123"), false},
		{"SHA256 valid", IsSHA256("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"), true},
		{"SHA256 invalid", IsSHA256("0123"), false},
		{"SHA512 valid", IsSHA512("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"), true},
		{"SHA512 invalid", IsSHA512("0123"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}
}
