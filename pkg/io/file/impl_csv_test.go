package file

import (
	"bytes"
	"strings"
	"testing"
)

func TestCsvExporter(t *testing.T) {
	// 准备测试数据
	headers := [][]string{{"Name", "Age"}, {"Location", "Job"}}
	rows := [][]string{{"Alice", "25"}, {"New York", "Engineer"}, {"Bob", "30"}, {"San Francisco", "Doctor"}}

	// 执行测试
	bs := []byte{}
	var buf = bytes.NewBuffer(bs)
	err := CsvExporter()(buf, headers, rows...)
	if err != nil {
		t.Errorf("CsvExporter() error = %v", err)
	}
	strGot := buf.String()
	// 验证结果
	expected := "Name,Age\nLocation,Job\nAlice,25\nNew York,Engineer\nBob,30\nSan Francisco,Doctor\n"
	if !strings.EqualFold(strGot, expected) {
		t.Errorf("CsvExporter() got = %v, want %v", strGot, expected)
	}
}

func TestCsvImporter(t *testing.T) {
	// 准备测试数据
	input := "Name姓名,Age\nLocation,Job\nAlice,25\nNew York,Engineer\nBob,30\nSan Francisco,Doctor\n"
	reader := strings.NewReader(input)

	// 执行测试
	got, err := CsvImporter(nil)(reader)
	if err != nil {
		t.Errorf("CsvImporter() error = %v", err)
	}

	// 验证结果
	expected := [][]string{
		{"Name姓名", "Age"},
		{"Location", "Job"},
		{"Alice", "25"},
		{"New York", "Engineer"},
		{"Bob", "30"},
		{"San Francisco", "Doctor"},
	}
	if !equal(got, expected) {
		t.Errorf("CsvImporter() got = %v, want %v", got, expected)
	}
}

// Helper function to compare two 2D string slices
func equal(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
