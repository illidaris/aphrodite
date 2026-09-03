package canal

import (
	"context"
	"errors"
	"fmt"
	"time"

	pbe "github.com/withlin/canal-go/protocol/entry"

	"github.com/withlin/canal-go/client"
	"google.golang.org/protobuf/proto"
)

type IConnector interface {
	// Id 返回连接器实例标识。
	Id() string
	Stats(ctx context.Context) OperateStats
	Close(ctx context.Context) error
	Run(ctx context.Context) error
}

type ISyncConnector interface {
	IConnector
}

type SyncConnector struct {
	SyncConnectorOption
	CanalConnector *client.SimpleCanalConnector
}

func (c *SyncConnector) Id() string {
	return c.CanalInstance
}

func (c *SyncConnector) Stats(ctx context.Context) OperateStats {
	return c.Outer.Stats()
}

func (c *SyncConnector) Close(ctx context.Context) error {
	_ = c.Outer.Close(ctx)
	return c.CanalConnector.DisConnection()
}

func (c *SyncConnector) Run(ctx context.Context) error {
	// 建立 Canal 连接并订阅配置的表，退出时释放连接和输出目标。
	err := c.CanalConnector.Connect()
	if err != nil {
		return err
	}
	defer c.CanalConnector.DisConnection()
	// 订阅表
	err = c.CanalConnector.Subscribe(c.TableFilter)
	if err != nil {
		return err
	}

	if c.Outer == nil {
		return fmt.Errorf("outer is nil")
	}
	defer func() {
		if err := c.Outer.Close(ctx); err != nil {
			c.Log.Error(ctx, "[gocanal]Run_OuterClose %v", err)
		}
	}()

	for {
		if err := ctx.Err(); err != nil {
			return nil
		}
		unit := int32(2)
		timout := int64(c.Timeout.Seconds())
		message, err := c.CanalConnector.GetWithOutAck(c.Batch, &timout, &unit)
		if err != nil {
			return err
		}
		batchId := message.Id
		// 批次Id -1 或者 消息集 为空， 则休眠
		if batchId == -1 || len(message.Entries) <= 0 {
			time.Sleep(300 * time.Millisecond)
			continue
		}
		entries, err := c.Convert(ctx, message.Entries...)
		if err != nil {
			c.Log.Error(ctx, "[gocanal]Run_Convert canal entries: %v", err)
			return err
		}
		ok, err := c.Outer.Sync(ctx, entries...)
		if err != nil {
			c.Log.Error(ctx, "[gocanal]Run_Sync canal entries to Outer: %v", err)
			return err
		}
		if !ok {
			stats := c.Outer.Stats()
			okErr := fmt.Errorf("%d stats %s", batchId, stats.GetMsg(","))
			c.Log.Error(ctx, "[gocanal]Run_SyncFunc_NoOk %s", okErr)
			return okErr
		}
		if err := c.CanalConnector.Ack(batchId); err != nil {
			c.Log.Error(ctx, "[gocanal]Run_Ack batch_%d %v", batchId, err)
			return err
		}
	}
}

func (c *SyncConnector) Convert(ctx context.Context, entries ...pbe.Entry) ([]Entry, error) {
	// 将 Canal ROWDATA 事件解码为目标端可消费的 Entry 列表。
	result := []Entry{}
	for k := range entries {
		entry := &entries[k]
		if entry.GetEntryType() != pbe.EntryType_ROWDATA {
			continue
		}
		change := new(pbe.RowChange)
		if err := proto.Unmarshal(entry.GetStoreValue(), change); err != nil {
			return result, fmt.Errorf("decode row change: %w", err)
		}
		tableName := entry.GetHeader().GetTableName()
		index := c.Index
		if len(c.TableMap) > 0 {
			v, ok := c.TableMap[tableName]
			if ok {
				index = v
			}
		}
		for _, row := range change.GetRowDatas() {
			// c.Log.Debug(ctx, "Sync_Change %v %v", tableName, row.GetBeforeColumns(), row.GetAfterColumns())
			columns := row.GetAfterColumns()
			action := ActionIndex
			if change.GetEventType() == pbe.EventType_DELETE {
				columns = row.GetBeforeColumns()
				action = ActionDelete
			}
			id, doc := DefaultColsToDoc(columns)
			if id == "" {
				return result, errors.New("doc has no primary key")
			}
			c.Log.Debug(ctx, "Sync_Values %s %s %v", tableName, id, doc)
			result = append(result, Entry{
				Id:    id,
				Act:   action,
				Index: index,
				Doc:   doc,
			})
		}
	}
	return result, nil
}
