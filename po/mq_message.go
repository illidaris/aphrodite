package po

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/illidaris/aphrodite/pkg/dependency"
	"github.com/illidaris/core"
)

var _ = dependency.ITask(&MqMessage{})
var _ = dependency.IEventMessage(&MqMessage{})

func NewMqMessage(
	ctx context.Context,
	bizId uint64,
	db string,
	category uint32,
	topic, key string,
	args any,
	timeout time.Duration,
) *MqMessage {
	bs, _ := json.Marshal(args)
	p := &MqMessage{}
	p.Locker = uuid.NewString()
	p.Expire = time.Now().Add(timeout).Unix()
	p.Timeout = int64(timeout.Seconds())
	p.BizId = bizId
	p.Db = db
	p.Name = topic
	p.Key = key
	p.Category = category
	p.Args = string(bs)
	p.TraceId = core.TraceID.GetString(ctx)
	return p
}

type MqMessage struct {
	TaskQueueMessage `gorm:"embedded"`
}

func (s MqMessage) TableName() string {
	return "aphrodite_mq_compensate"
}

func (s MqMessage) GetKey() []byte {
	return []byte(s.Key)
}

func (s MqMessage) GetValue() []byte {
	return []byte(s.Args)
}

func (s MqMessage) GetTopic() string {
	return s.Name
}

func (s MqMessage) ID() any {
	return s.Id
}

func (s MqMessage) GetTimeout() time.Duration {
	return time.Duration(s.Timeout) * time.Second
}

func (s MqMessage) GetBizId() uint64 {
	return s.BizId
}

func (s MqMessage) GetCategory() uint32 {
	return s.Category
}

func (s MqMessage) GetName() string {
	return s.Name
}

func (p MqMessage) Database() string {
	return ""
}
func (p MqMessage) ToRow() []string {
	return []string{}
}

func (p MqMessage) ToJson() string {
	bs, err := json.Marshal(&p)
	if err != nil {
		return ""
	}
	return string(bs)
}
