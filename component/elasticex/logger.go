package elasticex

import (
	"context"
	"fmt"

	"github.com/illidaris/aphrodite/pkg/logex"
	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

var _ = elastic.Logger(&BaseLogger{})

const (
	SqlLogFormat = "elapsed:%dms,affect:%d,err:%s,sql:%s"
)

func NewLogger() elastic.Logger {
	l := &BaseLogger{}
	l.core = zap.L().WithOptions(zap.AddCallerSkip(3))
	return l
}

type BaseLogger struct {
	core *zap.Logger
}

func (l *BaseLogger) Printf(msg string, args ...interface{}) {
	l.core.Info(fmt.Sprintf(msg, args...), logex.FieldsFromCtx(context.TODO())...)
}
