package convert

import (
	"reflect"
	"testing"
	"time"

	"github.com/spf13/cast"
)

func TestTable2Struct(t *testing.T) {
	type StudentInfo struct {
		Name      string    `json:"name"`
		Age       uint16    `json:"age"`
		IsStudent bool      `json:"is_student"`
		CreateAt  time.Time `json:"create_at"`
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
				{"name", "age", "is_student", "create_at"},
				{"Alice", "25", "true", "2023-05-07 22:23:24"},
				{"Bob", "30", "false", "2023-05-08 22:24:24"},
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
				{"Alice", "25"},
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
