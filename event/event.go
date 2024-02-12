package event

import (
	"context"

	"github.com/illidaris/aphrodite/component/gormex"
	"github.com/illidaris/aphrodite/component/kafkaex"
	"github.com/illidaris/aphrodite/pkg/dependency"
)

var (
	repo    dependency.IMQProducerRepository[dependency.IEventMessage]
	publish func(ctx context.Context, topic, key string, msg []byte) error
)

func Init(r dependency.IMQProducerRepository[dependency.IEventMessage], p func(ctx context.Context, topic, key string, msg []byte) error) {
	repo = r
	publish = p
}

func InitDefault() {
	repo = &gormex.EventRepository[dependency.IEventMessage]{}
	publish = kafkaex.GetKafkaManager().Publish
	Init(repo, publish)
}
