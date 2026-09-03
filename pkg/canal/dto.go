package canal

import (
	"fmt"
	"strings"
)

// Action identifies the Elasticsearch bulk operation to perform.
type Action string

const (
	ActionCreate Action = "create" // 仅当文档不存在时创建（若存在则报错）
	ActionUpdate Action = "update" // 局部更新文档字段（Body 必须包裹在 {"doc":{}} 中）
	ActionDelete Action = "delete" // 删除文档（此时 Body 传入 nil）
	ActionIndex  Action = "index"  // 创建或全量替换文档（最常用）
)

// Entry describes one document operation emitted by a connector.
type Entry struct {
	Id    string
	Index string
	Act   Action
	Doc   map[string]any
}

// OperateStats 批量操作的累计统计信息。
type OperateStats struct {
	NumAdded    uint64 `json:"num_added"`    // 加入批量处理队列的文档数。
	NumFlushed  uint64 `json:"num_flushed"`  // 已刷新的批次（HTTP bulk 请求）数。
	NumFailed   uint64 `json:"num_failed"`   // 处理失败的文档数。
	NumIndexed  uint64 `json:"num_indexed"`  // 执行 index 操作的文档数。
	NumCreated  uint64 `json:"num_created"`  // 成功创建的文档数。
	NumUpdated  uint64 `json:"num_updated"`  // 成功更新的文档数。
	NumDeleted  uint64 `json:"num_deleted"`  // 成功删除的文档数。
	NumRequests uint64 `json:"num_requests"` // 发起的 bulk 请求总数。
}

// GetTitle returns the title used for periodic statistics notifications.
func (i OperateStats) GetTitle() string {
	return "[gocanal]定时数据Stats同步"
}

// GetMsg formats statistics as a separator-delimited notification message.
func (i OperateStats) GetMsg(splitKey string) string {
	words := []string{
		fmt.Sprintf("追加文档：%d", i.NumAdded),
		fmt.Sprintf("刷新文档：%d", i.NumFlushed),
		fmt.Sprintf("失败文档：%d", i.NumFailed),
		fmt.Sprintf("保存文档：%d", i.NumIndexed),
		fmt.Sprintf("创建文档：%d", i.NumCreated),
		fmt.Sprintf("更新文档：%d", i.NumUpdated),
		fmt.Sprintf("删除文档：%d", i.NumDeleted),
		fmt.Sprintf("请求文档：%d", i.NumRequests),
	}
	return strings.Join(words, splitKey)
}
