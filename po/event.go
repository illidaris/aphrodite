package po

type EventArgs struct {
	BizSection
	Id    string // 消息唯一ID
	Topic string // 事件主题
	Key   string // 事件Key
	Value any    // 事件主体
}
