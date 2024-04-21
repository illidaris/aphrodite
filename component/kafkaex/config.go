package kafkaex

import (
	"time"

	"github.com/IBM/sarama"
)

// NewSASLConfig 返回一个使用给定用户名和密码创建的SASL配置。
func NewSASLConfig(app, user, password string) *sarama.Config {
	// 创建一个默认的配置
	config := DefaultConfig()
	config.ClientID = app
	// 启用SASL
	config.Net.SASL.Enable = true
	// 设置SASL机制为明文
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	// 设置SASL用户名
	config.Net.SASL.User = user
	// 设置SASL密码
	config.Net.SASL.Password = password
	// 允许没有收到acks而可以同时发送的最大batch数
	// config.Net.MaxOpenRequests = 5
	// 返回配置
	return config
}

// DefaultConfig 返回一个默认的Sarama配置对象
func DefaultConfig() *sarama.Config {
	config := sarama.NewConfig()
	// 设置消费组的重平衡策略为Sticky，使得分配更加均匀
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	// 生产者设置
	config.Producer.RequiredAcks = sarama.WaitForAll // 设置生产者发送消息时的确认策略为所有副本都确认，0-无需应答 1-本地确认 -1-全部确认
	config.Producer.Timeout = time.Second * 5        // 等待 WaitForAck的时间
	// config.Producer.MaxMessageBytes = 1000000 // 这个参数必须要小于broker中的`message.max.bytes`
	config.Producer.Partitioner = sarama.NewHashPartitioner // 分区策略
	config.Producer.Retry.Max = 3                           // 重新发送的次数
	config.Producer.Return.Successes = true
	return config
}
