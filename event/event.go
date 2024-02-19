package event

import (
	"context"
	"errors"

	"github.com/illidaris/aphrodite/component/gormex"
	"github.com/illidaris/aphrodite/component/kafkaex"
	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/aphrodite/po"
)

var (
	repo    dependency.IMQProducerRepository[po.MqMessage]
	publish func(ctx context.Context, topic, key string, msg []byte) error
)

func defaultPublish(ctx context.Context, topic, key string, msg []byte) error {
	return errors.New("no impl message queue")
}

func Init(r dependency.IMQProducerRepository[po.MqMessage], p func(ctx context.Context, topic, key string, msg []byte) error) {
	repo = r
	publish = p

}

func InitDefault() {
	repo = &gormex.EventRepository[po.MqMessage]{}
	publish = defaultPublish
	if kafkaex.GetKafkaManager() != nil {
		publish = kafkaex.GetKafkaManager().Publish
	}
	Init(repo, publish)
}
