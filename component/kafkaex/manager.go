package kafkaex

import (
	"context"

	"time"

	"github.com/IBM/sarama"
	"github.com/spf13/cast"
)

var m *KafkaManager

func InitManager(brokers ...string) *KafkaManager {
	m = NewKafkaManager(cast.ToStringSlice(brokers)...)
	m.Init()
	m.NewProducer()
	return m
}

func GetKafkaManager() *KafkaManager {
	return m
}

// NewKafkaManager new a kafka mq manager
func NewKafkaManager(brokers ...string) *KafkaManager {
	m = &KafkaManager{
		brokers:       brokers,
		groups:        map[string]IConsumerGroup{},
		consumerClose: []func(){},
	}
	return m
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

// init init
func (m *KafkaManager) Init() error {
	config := sarama.NewConfig()
	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.Net.SASL.User = "kafka"
	config.Net.SASL.Password = "pUuQNY9zG3NvObZxDwBhiHSBD6UxsQVx"

	// 同一个消费组中消费者订阅不同的topic,那么需要Sticky策略，该策略会使得分配更加均匀
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategySticky()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 0                   // 重新发送的次数
	config.Producer.Timeout = time.Millisecond * 10 // 等待 WaitForAck的时间
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewHashPartitioner

	// TODO: CONFIGS
	m.config = config
	client, err := sarama.NewClient(m.brokers, m.config)
	if err != nil {
		return err
	}
	m.client = client
	//client.LeastLoadedBroker().CreateTopics(&sarama.CreateTopicsRequest{})
	return nil
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
	if err := m.groups[groupid].CreateConsumer(handler, topics...); err != nil {
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
	return m.groups[groupid].GetConsumer(id), nil
}
