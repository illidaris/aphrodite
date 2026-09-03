package canal

import (
	"context"
	"fmt"
	"sync/atomic"

	pbe "github.com/withlin/canal-go/protocol/entry"
)

// ColsToKVsHandle 将 Canal 列集合转换为文档 ID 和文档字段。
type ColsToKVsHandle func(columns []*pbe.Column) (string, map[string]any)

// DefaultColsToDoc 使用列名作为字段名，并返回首个主键列的值作为文档 ID。
func DefaultColsToDoc(columns []*pbe.Column) (string, map[string]any) {
	doc := make(map[string]any, len(columns))
	var id string
	for _, column := range columns {
		value := column.GetValue()
		doc[column.GetName()] = value
		if column.GetIsKey() && id == "" {
			id = value
		}
	}
	return id, doc
}

var _ = ColsToKVsHandle(DefaultColsToDoc)

// KVsCbHandle 定义键值数据变更后的回调函数签名。
type KVsCbHandle func(context.Context, string, string, map[string]string)

// ValueChangeCbHandle 定义字段值变更后的回调函数签名。
type ValueChangeCbHandle func(context.Context, string, []*pbe.Column, []*pbe.Column)

// IOuter 是 Canal 数据同步目标的抽象接口。
type IOuter interface {
	Stats() OperateStats
	Close(context.Context) error
	Sync(ctx context.Context, entries ...Entry) (bool, error)
	Check(ctx context.Context, key string) error
	SyncStruct(ctx context.Context, key, index, mapping string) error
}

// BaseOuter 保存同步目标通用的表映射和日志配置。
type BaseOuter struct {
	TableMap map[string]string
	Log      ILogger
}

var _ = IOuter(&StdOuter{})

// StdOuter 是一个将同步内容输出到标准输出的简单同步目标，常用于示例和测试。
type StdOuter struct {
	BaseOuter
	Counter atomic.Uint64
}

func (i *StdOuter) Stats() OperateStats {
	total := i.Counter.Load()
	return OperateStats{
		NumAdded:    total,
		NumFlushed:  total,
		NumRequests: total,
	}
}

// Close 关闭标准输出同步目标；该实现无需释放资源。
func (i *StdOuter) Close(context.Context) error {
	return nil
}

// Sync 输出每条同步记录并累加已处理记录计数。
func (i *StdOuter) Sync(ctx context.Context, entries ...Entry) (bool, error) {
	for _, v := range entries {
		i.Counter.Add(1)
		fmt.Printf("[%s]%s %v", v.Act, v.Id, v.Doc)
	}
	return true, nil
}

// Check 检查目标中的键是否存在；标准输出目标始终返回成功。
func (i *StdOuter) Check(ctx context.Context, key string) error {
	return nil
}

// SyncStruct 接收结构映射配置；标准输出目标不执行实际映射。
func (i *StdOuter) SyncStruct(ctx context.Context, key, index, mapping string) error {
	return nil
}
