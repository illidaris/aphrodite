package event

import (
	"context"

	"github.com/illidaris/aphrodite/component/gormex"
	"github.com/illidaris/aphrodite/component/kafkaex"
	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/aphrodite/po"
)

var (
	repo    dependency.IMQProducerRepository[po.MqMessage]
	publish func(ctx context.Context, topic, key string, msg []byte) error
)

func Init(r dependency.IMQProducerRepository[po.MqMessage], p func(ctx context.Context, topic, key string, msg []byte) error) {
	repo = r
	publish = p
}

func InitDefault() {
	repo = &gormex.EventRepository[po.MqMessage]{}
	publish = kafkaex.GetKafkaManager().Publish
	Init(repo, publish)
}
