package kafkaex

import "fmt"

func BuildTopic(bizId int64, category int32, event string) string {
	return fmt.Sprintf("%d_%d_%s", bizId, category, event)
}
