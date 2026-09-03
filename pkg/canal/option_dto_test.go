package canal

import (
	"reflect"
	"testing"
	"time"
)

func TestNewSyncConnectorOptionDefaultsAndOverrides(t *testing.T) {
	log := DefaultLogger{}
	tableMap := map[string]string{"src.users": "users"}
	got := NewSyncConnectorOption(
		WithSyncCanalIp("127.0.0.1"), WithSyncCanalPort(11112),
		WithSyncBatch(100), WithSyncTimeout(2*time.Second),
		WithSyncTableMap(tableMap), WithSyncLogger(log), nil,
	)
	if got.CanalIp != "127.0.0.1" || got.CanalPort != 11112 || got.Batch != 100 || got.Timeout != 2*time.Second {
		t.Fatalf("同步选项未应用覆盖值: %+v", got)
	}
	if got.TableFilter != ".*\\..*" || got.CanalUser != "admin" || got.CanalSoTimeout != 60000 || got.CanalIdleTimeout != 3600000 {
		t.Fatalf("同步选项默认值错误: %+v", got)
	}
	if !reflect.DeepEqual(got.TableMap, tableMap) || got.Log == nil {
		t.Fatal("同步选项未保留映射或日志器")
	}
}

func TestNewMigrateConnectorOptionDefaultsAndOverrides(t *testing.T) {
	got := NewMigrateConnectorOption(WithMigrateDbIp("db"), WithMigrateDbPort(3307), WithMigrateCursorType(1), WithMigrateTableArgs([]int{1, 2}), nil)
	if got.DbIp != "db" || got.DbPort != 3307 || got.CursorType != 1 || !reflect.DeepEqual(got.TableArgs, []int{1, 2}) {
		t.Fatalf("迁移选项未应用覆盖值: %+v", got)
	}
	if got.CursorField != "id" || got.CursorPos != "0" || got.TableFilter != "%d" || got.Log == nil {
		t.Fatalf("迁移选项默认值错误: %+v", got)
	}
}

func TestOperateStatsMessage(t *testing.T) {
	stats := OperateStats{NumAdded: 1, NumFlushed: 2, NumFailed: 3, NumIndexed: 4, NumCreated: 5, NumUpdated: 6, NumDeleted: 7, NumRequests: 8}
	if got := stats.GetTitle(); got != "[gocanal]定时数据Stats同步" {
		t.Fatalf("标题 = %q", got)
	}
	want := "追加文档：1|刷新文档：2|失败文档：3|保存文档：4|创建文档：5|更新文档：6|删除文档：7|请求文档：8"
	if got := stats.GetMsg("|"); got != want {
		t.Fatalf("消息 = %q，期望 %q", got, want)
	}
}
