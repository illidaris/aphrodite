package po

import (
	"encoding/json"
	"time"

	"github.com/illidaris/aphrodite/pkg/dependency"
)

var _ = dependency.ITask(&MqMessage{})
var _ = dependency.IEventMessage(&MqMessage{})

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
