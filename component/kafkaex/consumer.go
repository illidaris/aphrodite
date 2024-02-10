package kafkaex

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

// IConsumer consumer interface
type IConsumer interface {
	sarama.ConsumerGroupHandler
	ID() string
	Go() error
	Close() error
}

type ConsumeHandler func(context.Context, *Message) (ReceiptStatus, error)

// IConsumer impl check
var _ = IConsumer(&Consumer{})

// NewConsumer new consumer
func NewConsumer(id string, group *ConsumerGroup, handler ConsumeHandler, topics ...string) *Consumer {
	// uuid.NewDCEPerson()
	c := &Consumer{
		id:       id,
		execFunc: handler,
		topics:   topics,
	}
	c.group = group
	return c
}

// Consumer consumer impl
type Consumer struct {
	id        string
	runFlag   int32 // 0-stop 1-running
	group     *ConsumerGroup
	topics    []string
	ready     chan bool
	closeFunc func()
	execFunc  ConsumeHandler
}

func (c *Consumer) ID() string {
	return c.id
}

// Close close consume goroutine
func (c *Consumer) Close() error {
	if c.closeFunc != nil {
		c.closeFunc()
	} else if c.runFlag == 0 {
		return ErrConsumerNotRun
	}
	return nil
}

// Go async consume message
func (c *Consumer) Go() error {
	// prevent concurrency exec
	if !atomic.CompareAndSwapInt32(&c.runFlag, 0, 1) {
		return fmt.Errorf("%s is already running", c.id)
	}
	// new perennial context
	ctx, cancel := context.WithCancel(context.Background())
	c.ready = make(chan bool)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := c.group.core.Consume(ctx, c.topics, c); err != nil {
				// 当setup失败的时候，error会返回到这里
				logger.Error(ctx, "Error from consumer: %v", err)
				return
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				logger.Println(ctx.Err())
				return
			}
			c.ready = make(chan bool)
		}
	}()
	<-c.ready
	logger.Info(ctx, "Sarama consumer up and running!...")
	// 保证在系统退出时，通道里面的消息被消费
	c.closeFunc = func() {
		logger.Println("kafka close")
		cancel()
		wg.Wait()
		atomic.CompareAndSwapInt32(&c.runFlag, 1, 0)
	}
	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	ctx := session.Context()
	if ctx == nil || ctx.Done() == nil {
		return ErrCtxNil
	}
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				logger.Printf("message channel was closed")
				return nil
			}
			var headers string
			if bs, err := json.Marshal(message.Headers); err != nil {
				headers = string(bs)
			}
			status, err := c.execFunc(ctx, &Message{
				Id:         uuid.NewString(),
				Headers:    headers,
				ConsumerId: c.id,
				Key:        message.Key,
				Offset:     message.Offset,
				Partition:  message.Partition,
				Topic:      message.Topic,
				Value:      message.Value,
				Ts:         message.Timestamp.Unix(),
				BlockTs:    message.BlockTimestamp.Unix(),
			})
			if err != nil {
				logger.Error(ctx, "exec %d %v", status, err)
				// 某些错误需要放到死信队列  有些放到重试队列
			} else if status == ReceiptSuccess || status == ReceiptAlreadyDo {
				session.MarkMessage(message, "")
			}
			logger.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-ctx.Done():
			return nil
		}
	}
}
