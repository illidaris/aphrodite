package kafkaex

type ReceiptStatus int32

const (
	ReceiptSuccess   ReceiptStatus = iota // 处理成功
	ReceiptAlreadyDo                      // 已经处理，重复消息
	ReceiptErrUnKnow                      // 未知错误
	ReceiptErrParse                       // 解析失败
	ReceiptRetryMax                       // 重试达到上限
)
