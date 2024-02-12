package dependency

type IMessage interface {
	GetTopic() string
	GetKey() []byte
	GetValue() []byte
}

type IEventMessage interface {
	ITask
	IMessage
}
