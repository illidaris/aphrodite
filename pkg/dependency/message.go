package dependency

import "time"

type IMessage interface {
	GetTopic() string
	GetKey() []byte
	GetValue() []byte
}

type IEventMessage interface {
	GetBizId() uint64
	GetTimeout() time.Duration
	GetUOWID() string
	IMessage
}
