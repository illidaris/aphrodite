package kafkaex

import (
	"context"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/illidaris/aphrodite/pkg/convert"
)

var (
	defaultManager *KafkaManager
	// onceInitManager is a sync.Once variable used to ensure that the init() function is only called once.
	onceInitManager       sync.Once
	onceInitSyncProducer  sync.Once
	onceInitAsyncProducer sync.Once
)

// InitDefaultManager 初始化默认的 KafkaManager 实例，并返回其指针。
// 参数：
//   - user: 用户名
//   - password: 密码
//   - brokers: Kafka broker 的地址列表
//
// 返回值：
//   - *KafkaManager: KafkaManager 实例的指针
//   - error: 错误信息，如果没有错误发生则为 nil
func InitDefaultManager(opts ...OptionsFunc) (*KafkaManager, error) {
	// 使用 sync.Once 来保证函数只会执行一次
	onceInitManager.Do(func() {
		// 创建一个新的 KafkaManager 实例，并传入 SASL 配置和 broker 地址列表
		m, err := NewKafkaManager(opts...)
		if err != nil {
			if logger != nil {
				logger.Error(context.Background(), err.Error())
			}
			println("InitDefaultManager_NewKafkaManager", err.Error())
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
func NewKafkaManager(opts ...OptionsFunc) (*KafkaManager, error) {
	// 创建 KafkaManager 实例
	manager := &KafkaManager{
		groups:        map[string]IConsumerGroup{},
		consumerClose: []func(){},
	}

	for _, o := range opts {
		o(&manager.Options)
	}
	// 设置 Kafka 客户端和配置到 KafkaManager 实例
	return manager, nil
}

// KafkaManager kafka mq manager
type KafkaManager struct {
	Options
	groups        map[string]IConsumerGroup // consumer group map
	consumerClose []func()                  // consumer close func
	producerSync  sarama.SyncProducer
	producerAsync sarama.AsyncProducer
}

// GetSyncProducer returns a synchronized Kafka producer.
func (m *KafkaManager) GetSyncProducer() sarama.SyncProducer {
	// Use a sync.Once to ensure that the producer is initialized only once.
	onceInitSyncProducer.Do(func() {
		// Create a new synchronized Kafka producer from the existing client.
		config := NewSASLConfig(m.App, m.User, m.Pwd)
		client, err := sarama.NewClient(m.Addrs, config)
		if err != nil {
			logger.Error(context.TODO(), "GetSyncProducer_NewClient err %v", err)
			return
		}
		producer, err := sarama.NewSyncProducerFromClient(client)
		if err != nil {
			// Print an error message if the producer creation fails.
			println("InitDefaultManager_NewProducer", err.Error())
			return
		}
		// Store the producer in the KafkaManager instance.
		m.producerSync = producer
	})
	// Return the synchronized Kafka producer.
	return m.producerSync
}

// GetASyncProducer returns an asynchronous Kafka producer.
func (m *KafkaManager) GetASyncProducer() sarama.AsyncProducer {
	// Ensure that the async producer is initialized only once.
	onceInitAsyncProducer.Do(func() {
		// Create a new async producer from the Kafka client.
		config := NewSASLConfig(m.App, m.User, m.Pwd)
		client, err := sarama.NewClient(m.Addrs, config)
		if err != nil {
			logger.Error(context.TODO(), "GetASyncProducer_NewClient err %v", err)
			return
		}
		producer, err := sarama.NewAsyncProducerFromClient(client)
		if err != nil {
			// Print an error message if the producer creation fails.
			println("InitDefaultManager_NewProducer_Async", err.Error())
			return
		}
		// Start consuming messages from the producer.
		producer.Input()
		// Store the created producer in the KafkaManager instance.
		m.producerAsync = producer
	})
	// Return the KafkaManager's async producer.
	return m.producerAsync
}

// Publish is a function that publishes a message to a Kafka topic.
// It takes a context, topic name, key, and message as input and returns an error.
func (m *KafkaManager) Publish(ctx context.Context, topic, key string, msg []byte) error {
	var (
		mqMsg = &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(key),
			Value: sarama.ByteEncoder(msg),
		}
		producer = m.GetSyncProducer()
	)
	if producer == nil {
		return ErrProducerNoFound
	}
	beg := time.Now()
	partition, offset, err := producer.SendMessage(mqMsg)
	logger.Info(ctx, "sendmsg%d-%d: %,%dms,%v", partition, offset, string(msg), time.Since(beg).Milliseconds(), err)
	if err != nil {
		return err
	}
	return nil
}

// NewConsumer creates a new consumer for the specified group ID and handler.
// It returns an error if the group ID is not found or if there is an error creating the consumer group.
func (m *KafkaManager) NewConsumer(id, groupid string, handler ConsumeHandler, topics ...string) error {
	if _, ok := m.groups[groupid]; !ok {
		config := NewSASLConfig(m.App, m.User, m.Pwd)
		client, err := sarama.NewClient(m.Addrs, config)
		if err != nil {
			logger.Error(context.TODO(), "NewConsumer_NewClient%s_%s err %v", groupid, id, err)
			return err
		}
		group, err := NewConsumerGroup(groupid, client)
		if err != nil {
			logger.Error(context.TODO(), "NewConsumer_NewConsumerGroup%s_%s err %v", groupid, id, err)
			return err
		}
		m.groups[groupid] = group
	}
	group := m.groups[groupid]
	if group == nil {
		return ErrGroupNoFound
	}
	if id == "" {
		id = convert.RandomID()
	}
	if err := group.CreateConsumer(id, handler, topics...); err != nil {
		return err
	}
	return nil
}

// ConsumersGo starts all consumers in the KafkaManager.
// It loops through all groups and consumers, and starts each consumer's goroutine.
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

// ConsumerClose closes all consumers in the KafkaManager.
// It loops through all groups and consumers, and closes each consumer.
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

// GetConsumer returns the consumer with the specified ID from the specified group ID.
// It returns an error if the group ID or consumer ID is not found.
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
