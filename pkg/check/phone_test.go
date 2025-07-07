package check

import (
	"testing"
)

// TestIsValidChineseMobile 测试中国大陆手机号验证函数
func TestIsValidChineseMobile(t *testing.T) {
	tests := []struct {
		name   string // 测试用例名称
		mobile string // 输入手机号
		want   bool   // 期望结果
	}{
		// 有效手机号测试
		{"valid_13x", "13000000000", true}, // 第二位边界值3
		{"valid_15x", "15500000000", true}, // 第二位中间值5
		{"valid_19x", "19000000000", true}, // 第二位边界值9
		{"valid_19x", "19900000000", true}, // 第二位边界值9

		// 无效手机号测试
		{"invalid_prefix", "21000000000", false},           // 非1开头
		{"invalid_second_digit_low", "12000000000", false}, // 第二位<3
		{"non_numeric_chars", "138abcd1234", false},        // 包含字母
		{"too_short", "1380013800", false},                 // 10位数字
		{"too_long", "138001380000", false},                // 12位数字
		{"with_dash", "138-0013-8000", false},              // 含分隔符
		{"empty_string", "", false},                        // 空字符串
		{"all_zeros", "00000000000", false},                // 全零数字
		{"unicode_digits", "138００１３８０００", false},           // 全角数字
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidChineseMobile(tt.mobile)
			if got != tt.want {
				t.Errorf("IsValidChineseMobile(%q) = %v, want %v", tt.mobile, got, tt.want)
			}
		})
	}
}
