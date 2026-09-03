package canal

import (
	"context"
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
				if gotDoc[key] != want {
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
