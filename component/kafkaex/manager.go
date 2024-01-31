package kafkaex

import (
	"context"
	"sync"

	"github.com/IBM/sarama"
)

var defaultManager *KafkaManager
var once sync.Once

// InitDefaultManager 初始化默认的 KafkaManager 实例，并返回其指针。
// 参数：
//   - user: 用户名
//   - password: 密码
//   - brokers: Kafka broker 的地址列表
//
// 返回值：
//   - *KafkaManager: KafkaManager 实例的指针
//   - error: 错误信息，如果没有错误发生则为 nil
func InitDefaultManager(user, password string, brokers ...string) (*KafkaManager, error) {
	// 使用 sync.Once 来保证函数只会执行一次
	once.Do(func() {
		// 创建一个新的 KafkaManager 实例，并传入 SASL 配置和 broker 地址列表
		m, err := NewKafkaManager(NewSASLConfig(user, password), brokers...)
		if err != nil {
			println("InitDefaultManager_NewKafkaManager", err.Error())
			return
		}
		// 创建 KafkaManager 实例的生产者
		err = m.NewProducer()
		if err != nil {
			println("InitDefaultManager_NewProducer", err.Error())
			return
		}
		// 将 KafkaManager 实例赋值给默认的实例变量
		defaultManager = m
	})
	// 返回 KafkaManager 实例的指针和 nil 错误
	return GetKafkaManager(), nil
}

// GetKafkaManager 返回一个指向 KafkaManager 的指针
func GetKafkaManager() *KafkaManager {
	return defaultManager
}

// NewKafkaManager 创建一个新的 KafkaManager 实例。
// 参数 config 是 Kafka 客户端配置。
// 参数 brokers 是 Kafka broker 的地址列表。
// 返回 KafkaManager 实例和可能的错误。
func NewKafkaManager(config *sarama.Config, brokers ...string) (*KafkaManager, error) {
	// 创建 Kafka 客户端
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, err
	}

	// 创建 KafkaManager 实例
	manager := &KafkaManager{
		brokers:       brokers,
		groups:        map[string]IConsumerGroup{},
		consumerClose: []func(){},
	}

	// 设置 Kafka 客户端和配置到 KafkaManager 实例
	manager.client = client
	manager.config = config

	return manager, nil
}

// KafkaManager kafka mq manager
type KafkaManager struct {
	brokers       []string
	config        *sarama.Config
	client        sarama.Client             // sarama kafka sdk client
	groups        map[string]IConsumerGroup // consumer group map
	consumerClose []func()                  // consumer close func
	producerSync  sarama.SyncProducer
}

func (m *KafkaManager) NewProducer() error {
	producer, err := sarama.NewSyncProducerFromClient(m.client)
	if err != nil {
		return err
	}
	m.producerSync = producer
	return nil
}

func (m *KafkaManager) Publish(ctx context.Context, topic, key string, msg []byte) error {
	mqMsg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(msg)}
	partition, offset, err := m.producerSync.SendMessage(mqMsg)
	if err != nil {
		return err
	}
	logger.Info(ctx, "sendmsg%d-%d: %s", partition, offset, string(msg))
	return nil
}

func (m *KafkaManager) NewConsumer(groupid string, handler ConsumeHandler, topics ...string) error {
	if _, ok := m.groups[groupid]; !ok {
		group, err := NewConsumerGroup(groupid, m.client)
		if err != nil {
			return err
		}
		m.groups[groupid] = group
	}
	group := m.groups[groupid]
	if group == nil {
		return ErrGroupNoFound
	}
	if err := group.CreateConsumer(handler, topics...); err != nil {
		return err
	}
	return nil
}

func (m *KafkaManager) ConsumersGo() {
	for _, g := range m.groups {
		if g == nil {
			continue
		}
		for _, c := range g.ConsumerMap() {
			if c == nil {
				continue
			}
			err := c.Go()
			if err != nil {
				logger.Error(context.TODO(), "group[%s]consumer[%s]go err %v", g.ID(), c.ID(), err)
			}
		}
	}
}

func (m *KafkaManager) ConsumerClose() {
	for _, g := range m.groups {
		if g == nil {
			continue
		}
		for _, c := range g.ConsumerMap() {
			if c == nil {
				continue
			}
			err := c.Close()
			if err != nil {
				logger.Error(context.TODO(), "group[%s]consumer[%s]close err %v", g.ID(), c.ID(), err)
			}
		}
	}
}

func (m *KafkaManager) GetConsumer(groupid, id string) (IConsumer, error) {
	if _, ok := m.groups[groupid]; !ok {
		return nil, ErrGroupNoFound
	}
	group := m.groups[groupid]
	if group == nil {
		return nil, ErrGroupNoFound
	}
	return group.GetConsumer(id), nil
}
