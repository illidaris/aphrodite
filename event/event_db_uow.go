package event

import (
	"context"
	"errors"
	"time"

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
	if repo == nil {
		return errors.New("no impl repo")
	}
	action := repo.InsertAction(ctx, t.id, ent)
	if action != nil {
		fs = append(fs, action)
	}
	err := uow.Execute(ctx, fs...)
	if err != nil {
		return err
	}
	err = publish(ctx, ent.GetTopic(), string(ent.GetKey()), ent.GetValue())
	go func(e *po.MqMessage, publishErr error) {
		newCtx := core.TraceID.SetString(context.Background(), e.TraceId)
		if publishErr != nil {
			updateE := &po.MqMessage{}
			updateE.Id = e.Id
			updateE.BizId = e.BizId
			updateE.LastError = publishErr.Error()
			if len(updateE.LastError) > 255 {
				updateE.LastError = updateE.LastError[:255]
			}
			updateE.LastExecAt = time.Now().Unix()
			_, _ = repo.BaseUpdate(newCtx, updateE,
				dependency.WithDataBase(e.Db))
		} else {
			_, _ = repo.BaseDelete(newCtx, e,
				dependency.WithDataBase(e.Db))
		}
	}(ent, err)
	return nil
}
