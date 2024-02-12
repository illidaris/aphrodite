package event

import (
	"context"

	"github.com/illidaris/aphrodite/component/gormex"
	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/aphrodite/po"
	"github.com/illidaris/core"
)

var _ = dependency.IUnitOfWork(EventTransactionImpl{})

func NewUnitOfWork(id string, event *po.MqMessage) dependency.IUnitOfWork {
	return &EventTransactionImpl{
		id:    id,
		event: event,
	}
}

type EventTransactionImpl struct {
	id    string
	event *po.MqMessage
}

func (t EventTransactionImpl) Execute(ctx context.Context, fs ...dependency.DbAction) (e error) {
	ent := t.event
	uow := gormex.NewUnitOfWork(t.id)
	action := repo.InsertAction(ctx, t.id, ent)
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
	go func(e *po.MqMessage) {
		newCtx := core.TraceID.SetString(context.Background(), e.TraceId)
		_, _ = repo.BaseDelete(newCtx, e,
			dependency.WithDataBase(e.Db))
	}(ent)
	return nil
}
