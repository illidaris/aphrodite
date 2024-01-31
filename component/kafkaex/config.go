package kafkaex

import (
	"time"

	"github.com/IBM/sarama"
)

// NewSASLConfig 返回一个使用给定用户名和密码创建的SASL配置。
func NewSASLConfig(user, password string) *sarama.Config {
	// 创建一个默认的配置
	config := DefaultConfig()
	// 启用SASL
	config.Net.SASL.Enable = true
	// 设置SASL机制为明文
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	// 设置SASL用户名
	config.Net.SASL.User = user
	// 设置SASL密码
	config.Net.SASL.Password = password
	// 返回配置
	return config
}

// DefaultConfig 返回一个默认的Sarama配置对象
func DefaultConfig() *sarama.Config {
	config := sarama.NewConfig()

	// 设置消费组的重平衡策略为Sticky，使得分配更加均匀
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategySticky()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3             // 重新发送的次数
	config.Producer.Timeout = time.Second * 3 // 等待 WaitForAck的时间
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewHashPartitioner
	return config
}
