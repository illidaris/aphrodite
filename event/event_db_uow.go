package event

import (
	"context"

	"github.com/illidaris/aphrodite/component/gormex"
	"github.com/illidaris/aphrodite/pkg/dependency"
)

var _ = dependency.IUnitOfWork(&EventTransactionImpl{})

func NewUnitOfWork(event dependency.IEventMessage) dependency.IUnitOfWork {
	return &EventTransactionImpl{
		event: event,
	}
}

type EventTransactionImpl struct {
	event dependency.IEventMessage
}

func (t *EventTransactionImpl) Execute(ctx context.Context, fs ...dependency.DbAction) (e error) {
	ent := t.event
	uow := gormex.NewUnitOfWork(ent.GetUOWID())
	action, locker := repo.InsertAction(ctx, ent)
	if action != nil {
		fs = append(fs, action)
	}
	err := uow.Execute(ctx, fs...)
	if err != nil {
		return err
	}
	err = publish(ctx, ent.GetTopic(), string(ent.GetKey()), ent.GetValue())
	if err != nil {
		return err
	}
	_, _ = repo.ClearByLocker(ctx, locker)
	return nil
}

// func Retry() {
// 	repo.WaitExecWithLock(context.Background(), 1, 1, 1, "test", 10*time.Second)
// }
