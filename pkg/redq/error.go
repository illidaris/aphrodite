package redq

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

var _ = asynq.ErrorHandler(&ZapErrorHandler{})

type ZapErrorHandler struct {
	core *zap.Logger
}

func NewZapErrHandler() *ZapErrorHandler {
	z := zap.L().WithOptions(zap.AddCallerSkip(2))
	return &ZapErrorHandler{
		core: z,
	}
}

func (h ZapErrorHandler) HandleError(ctx context.Context, task *asynq.Task, err error) {
	if task == nil {
		h.core.Error(fmt.Sprintf("消息消费失败：%v", err))
		return
	}
	h.core.Error(fmt.Sprintf("[%v]消息消费失败：%v,原文：%v", task.Type(), err, string(task.Payload())))
}
