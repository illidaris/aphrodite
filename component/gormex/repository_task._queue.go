package gormex

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/illidaris/aphrodite/dto"
	"github.com/illidaris/aphrodite/pkg/dependency"
	"gorm.io/gorm"
)

type TaskQueueRepository[T dependency.ITask] struct {
	BaseRepository[T]
}

// CountByCreator 请求者当前执行数量
func (r TaskQueueRepository[T]) CountByCreator(ctx context.Context, createBy, status int64) (int64, error) {
	return r.BaseCount(ctx, dependency.WithConds("createBy = ? and status = ?", createBy, status))
}

// WaitExecWithLock 需要锁定的记录
func (r TaskQueueRepository[T]) WaitExecWithLock(ctx context.Context, t T, batch int) (string, int64, error) {
	var (
		locker = uuid.NewString()
		page   = &dto.Page{PageIndex: 1, PageSize: int64(batch), Sorts: []string{"createAt|asc"}}
	)
	opts := []dependency.BaseOptionFunc{
		dependency.WithConds("expire < unix_timestamp() AND `bizId` = ? AND `category` = ? AND `name` =?",
			t.GetBizId(),
			t.GetCategory(),
			t.GetName()),
		dependency.WithPage(page),
	}
	result := r.BuildFrmOptions(ctx, new(T), opts...).Updates(map[string]interface{}{
		// 如果记录里设定了超时时间则采用改超时时间，当前时间均采用数据库时间，保障所有节点计算时间一致【防止节点时间不一致】
		"expire":  gorm.Expr(fmt.Sprintf("IF(`timeout` > 0, unix_timestamp() + `timeout` , unix_timestamp() + %d)", int64(t.GetTimeout().Seconds()))),
		"retries": gorm.Expr("retries + 1"),
		"locker":  locker,
	})
	return locker, result.RowsAffected, result.Error
}

// FindLockeds 找到被锁定的记录
func (r TaskQueueRepository[T]) FindLockeds(ctx context.Context, locker string) ([]T, error) {
	return r.BaseQuery(ctx, dependency.WithConds("locker = ?", locker))
}

// ReportExecResult 汇报执行结果
func (r TaskQueueRepository[T]) ReportExecResult(ctx context.Context, id int64, locker string, execResult string, execErr error) (int64, error) {
	if locker == "" {
		return 0, nil
	}
	var (
		t            = new(T)
		now          = time.Now()
		lastErrorStr = ""
	)
	if execErr != nil {
		lastErrorStr = execErr.Error()
		if len(lastErrorStr) > 255 {
			lastErrorStr = lastErrorStr[:255]
		}
	}
	opts := []dependency.BaseOptionFunc{
		dependency.WithConds("id = ? AND locker = ?", id, locker),
	}
	// 没有错误，表示执行成功
	if execErr == nil {
		result := r.BuildFrmOptions(ctx, t, opts...).Delete(t)
		return result.RowsAffected, result.Error
	}
	result := r.BuildFrmOptions(ctx, t, opts...).Updates(
		map[string]interface{}{
			"lastExecAt": now.Unix(),
			"lastError":  lastErrorStr,
			"locker":     "-",
			"expire":     0,
		})
	return result.RowsAffected, result.Error
}
