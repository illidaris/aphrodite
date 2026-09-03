package canal

import (
	"context"
	"reflect"
	"testing"

	pbe "github.com/withlin/canal-go/protocol/entry"
)

func TestDefaultColsToDoc(t *testing.T) {
	tests := []struct {
		name    string
		cols    []*pbe.Column
		wantID  string
		wantDoc map[string]any
	}{
		{name: "空列", cols: nil, wantDoc: map[string]any{}},
		{name: "首个主键", cols: []*pbe.Column{
			{Name: "id", Value: "7", IsKey: true},
			{Name: "name", Value: "张三"},
		}, wantID: "7", wantDoc: map[string]any{"id": "7", "name": "张三"}},
		{name: "重复主键取首个", cols: []*pbe.Column{
			{Name: "id", Value: "7", IsKey: true},
			{Name: "id2", Value: "8", IsKey: true},
		}, wantID: "7", wantDoc: map[string]any{"id": "7", "id2": "8"}},
		{name: "按 MySQL 类型转换", cols: []*pbe.Column{
			{Name: "signed", Value: "-7", MysqlType: "int(11)"},
			{Name: "unsigned", Value: "7", MysqlType: "bigint unsigned"},
			{Name: "score", Value: "3.14", MysqlType: "decimal(10,2)"},
			{Name: "enabled", Value: "true", MysqlType: "boolean"},
			{Name: "metadata", Value: `{"active":true}`, MysqlType: "json"},
		}, wantDoc: map[string]any{
			"signed":   int64(-7),
			"unsigned": uint64(7),
			"score":    3.14,
			"enabled":  true,
			"metadata": map[string]any{"active": true},
		}},
		{name: "空值和无法转换的值", cols: []*pbe.Column{
			nil,
			{Name: "deleted_at", Value: "", MysqlType: "datetime", IsNullPresent: &pbe.Column_IsNull{IsNull: true}},
			{Name: "invalid_number", Value: "unknown", MysqlType: "int"},
		}, wantDoc: map[string]any{"deleted_at": nil, "invalid_number": "unknown"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotDoc := DefaultColsToDoc(tt.cols)
			if gotID != tt.wantID {
				t.Errorf("文档 ID = %q，期望 %q", gotID, tt.wantID)
			}
			if len(gotDoc) != len(tt.wantDoc) {
				t.Fatalf("文档字段数 = %d，期望 %d", len(gotDoc), len(tt.wantDoc))
			}
			for key, want := range tt.wantDoc {
				if !reflect.DeepEqual(gotDoc[key], want) {
					t.Errorf("字段 %q = %v，期望 %v", key, gotDoc[key], want)
				}
			}
		})
	}
}

func TestStdOuterSyncAndStats(t *testing.T) {
	outer := &StdOuter{}
	ok, err := outer.Sync(context.Background(), Entry{Id: "1"}, Entry{Id: "2"})
	if err != nil || !ok {
		t.Fatalf("Sync() = (%v, %v)，期望 (true, nil)", ok, err)
	}
	stats := outer.Stats()
	if stats.NumAdded != 2 || stats.NumFlushed != 2 || stats.NumRequests != 2 {
		t.Fatalf("统计信息 = %+v，期望三项均为 2", stats)
	}
}
