package kafkaex

import (
	"context"
	"sync"

	"github.com/IBM/sarama"
)

var _ = IConsumerGroup(&ConsumerGroup{})

type IConsumerGroup interface {
	ID() string
	CreateConsumer(id string, h ConsumeHandler, topics ...string) error
	GetConsumer(id string) IConsumer
	ConsumerMap() map[string]IConsumer
}

// NewConsumerGroup 创建一个新的ConsumerGroup实例。
// 参数：
//   - groupid：消费者组的ID
//   - client：Sarama客户端实例
//
// 返回值：
//   - *ConsumerGroup：新创建的ConsumerGroup实例
//   - error：如果创建失败，返回错误信息
func NewConsumerGroup(groupid string, client sarama.Client) (*ConsumerGroup, error) {
	g := &ConsumerGroup{
		id: groupid,
		rw: sync.RWMutex{},
	}
	group, err := sarama.NewConsumerGroupFromClient(g.id, client)
	if err != nil {
		return g, err
	}
	g.core = group
	return g, nil
}

type ConsumerGroup struct {
	rw          sync.RWMutex         // lock
	id          string               // consume group id
	core        sarama.ConsumerGroup // consume group core
	consumerMap map[string]IConsumer // consumer map
}

func (g *ConsumerGroup) ID() string {
	return g.id
}

func (g *ConsumerGroup) Close() {
	g.rw.Lock()
	defer g.rw.Unlock()
	if g.core == nil {
		return
	}
	if err := g.core.Close(); err != nil {
		logger.Error(context.TODO(), "ConsumerGroup[%s].CLose %v", g.id, err)
	}
}

func (g *ConsumerGroup) CreateConsumer(id string, h ConsumeHandler, topics ...string) error {
	g.rw.Lock()
	defer g.rw.Unlock()
	consumer := NewConsumer(id, g, h, topics...)
	if g.consumerMap == nil {
		g.consumerMap = map[string]IConsumer{}
	}
	if _, ok := g.consumerMap[consumer.id]; ok {
		return ErrConsumerExist
	}
	g.consumerMap[consumer.id] = consumer
	return nil
}

func (g *ConsumerGroup) GetConsumer(id string) IConsumer {
	g.rw.RLock()
	defer g.rw.RUnlock()
	if v, ok := g.consumerMap[id]; ok {
		return v
	}
	return nil
}

func (g *ConsumerGroup) ConsumerMap() map[string]IConsumer {
	return g.consumerMap
}
