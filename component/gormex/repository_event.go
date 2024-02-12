package gormex

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/illidaris/aphrodite/dto"
	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/aphrodite/po"
	"github.com/illidaris/core"
	"gorm.io/gorm"
)

type EventRepository struct {
	BaseRepository[po.MqMessage[any]]
}

func (r EventRepository) InsertAction(ctx context.Context, message dependency.IEventMessage) (func(context.Context) error, string) {
	p := &po.MqMessage[any]{}
	p.TraceId = core.TraceID.GetString(ctx)
	p.Locker = uuid.NewString()
	p.BizId = message.GetBizId()
	p.Topic = message.GetTopic()
	p.Key = string(message.GetKey())
	p.Args = message
	p.Expire = time.Now().Add(message.GetTimeout()).Unix()
	p.Timeout = int64(message.GetTimeout().Seconds())
	return func(ctx context.Context) error {
		_, err := r.BaseCreate(
			ctx,
			[]*po.MqMessage[any]{p},
			dependency.WithDbShardingKey(message.GetBizId()),
			dependency.WithIgnore(true),
		)
		return err
	}, p.Locker
}

// WaitExecWithLock 需要锁定的记录
func (r EventRepository) WaitExecWithLock(ctx context.Context, bizId, category, batch int, name string, timeout time.Duration) (string, int64, error) {
	var (
		locker = uuid.NewString()
		page   = &dto.Page{PageIndex: 1, PageSize: int64(batch), Sorts: []string{"createAt|asc"}}
		p      = &po.MqMessage[any]{}
	)
	p.BizId = uint64(bizId)
	opts := []dependency.BaseOptionFunc{
		dependency.WithConds("expire < unix_timestamp() AND `bizId` = ? AND `category` = ? AND `name` =?",
			bizId,
			category,
			name),
		dependency.WithPage(page),
	}
	result := r.BuildFrmOptions(ctx, p, opts...).Updates(map[string]interface{}{
		// 如果记录里设定了超时时间则采用改超时时间，当前时间均采用数据库时间，保障所有节点计算时间一致【防止节点时间不一致】
		"expire":  gorm.Expr(fmt.Sprintf("IF(`timeout` > 0, unix_timestamp() + `timeout` , unix_timestamp() + %d)", int64(timeout.Seconds()))),
		"retries": gorm.Expr("retries + 1"),
		"locker":  locker,
	})
	return locker, result.RowsAffected, result.Error
}

// FindLockeds 找到被锁定的记录
func (r EventRepository) FindLockeds(ctx context.Context, locker string) ([]dependency.ITask, error) {
	task := []dependency.ITask{}
	ps, err := r.BaseQuery(ctx, dependency.WithConds("locker = ?", locker))
	for _, p := range ps {
		v := p
		task = append(task, &v)
	}
	return task, err
}

func (r EventRepository) Clear(ctx context.Context, id string) (int64, error) {
	return r.BaseDelete(ctx, &po.MqMessage[any]{}, dependency.WithConds("id = ?", id))
}

func (r EventRepository) ClearByLocker(ctx context.Context, locker string) (int64, error) {
	newCtx := core.TraceID.SetString(context.Background(), core.TraceID.GetString(ctx))
	return r.BaseDelete(newCtx, &po.MqMessage[any]{}, dependency.WithConds("locker = ?", locker))
}
