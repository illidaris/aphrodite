package event

import (
	"context"

	"github.com/illidaris/aphrodite/component/gormex"
	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/aphrodite/po"
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
	_, _ = repo.BaseDelete(ctx, ent, dependency.WithConds("id = ?", ent.Id))
	return nil
}
