package convert

import (
	"reflect"
	"testing"
	"time"

	"github.com/spf13/cast"
)

func TestTable2Struct(t *testing.T) {
	type StudentInfo struct {
		Id        int       `json:"id"`
		Name      string    `json:"name"`
		Age       uint16    `json:"age"`
		IsStudent bool      `json:"is_student"`
		CreateAt  time.Time `json:"create_at"`
		Desc      string    `json:"desc"`
	}
	// 定义测试用例
	tests := []struct {
		name    string
		rows    [][]string
		opts    []Table2StructOptionFunc
		want    interface{}
		wantErr bool
	}{
		{
			name: "Valid conversion",
			rows: [][]string{
				{"name", "age", "is_student", "create_at", "", "", ""},
				{"Alice", "25", "true", "2023-05-07 22:23:24", "", "", ""},
				{"Bob", "30", "false", "2023-05-08 22:24:24", "", "", ""},
				{"", "", "", "", "", "x", ""},
				{"", "", "", "", "Y", "x", ""},
				{"", "", "", "", "", "", "", "xx"},
			},
			opts: []Table2StructOptionFunc{
				WithStructTag("json"),
				WithHeadIndex(0),
				WithStartRowIndex(1),
				WithLimit(2),
			},
			want: []StudentInfo{
				{Name: "Alice", Age: 25, IsStudent: true, CreateAt: cast.ToTime("2023-05-07 22:23:24")},
				{Name: "Bob", Age: 30, IsStudent: false, CreateAt: cast.ToTime("2023-05-08 22:24:24")},
			},
			wantErr: false,
		},
		{
			name: "Invalid dst type",
			rows: [][]string{
				{"Name", "Age"},
			},
			opts: []Table2StructOptionFunc{
				WithStructTag("json"),
				WithHeadIndex(0),
				WithStartRowIndex(1),
			},
			want:    []StudentInfo{},
			wantErr: true,
		},
		// 可以添加更多的测试用例
	}

	// 遍历测试用例
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 调用Table2Struct函数
			got := []StudentInfo{}
			err := Table2Struct(&got, tt.rows, tt.opts...)
			if (err != nil) && !tt.wantErr {
				t.Errorf("Table2Struct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// 比较结果
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Table2Struct() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestStruct2Table tests the Struct2Table function.
func TestStruct2Table(t *testing.T) {
	// Define a test struct type
	type MyStruct struct {
		Name  string `table2struct:"name"`
		Age   int    `table2struct:"age"`
		Email string `table2struct:"email"`
	}

	// Create a slice of MyStruct instances
	dst := []interface{}{
		MyStruct{Name: "Alice", Email: "alice@example.com", Age: 25},
		MyStruct{Age: 30, Email: "bob@example.com", Name: "Bob"},
	}

	// Call Struct2Table function
	heads, rows, err := Struct2Table(dst, WithStructTag("table2struct"))
	if err != nil {
		t.Errorf("Struct2Table returned an error: %v", err)
	}

	// Define expected results
	expectedHeads := []string{"name", "age", "email"}
	expectedRows := [][]string{
		{"Alice", "25", "alice@example.com"},
		{"Bob", "30", "bob@example.com"},
	}

	// Check if the heads match the expected heads
	if !reflect.DeepEqual(heads, expectedHeads) {
		t.Errorf("Expected heads: %v, but got: %v", expectedHeads, heads)
	}

	// Check if the rows match the expected rows
	if !reflect.DeepEqual(rows, expectedRows) {
		t.Errorf("Expected rows: %v, but got: %v", expectedRows, rows)
	}
}
