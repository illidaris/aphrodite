package dto

import (
	"time"

	"github.com/illidaris/aphrodite/pkg/dependency"
)

var _ = dependency.ISchedule(&ScheduleBase{})

type ScheduleBase struct {
	Batch   int64 `json:"batch" form:"batch" url:"batch"`       // 执行的批量
	Timeout int64 `json:"timeout" form:"timeout" url:"timeout"` // 执行超时时间(秒)
}

func (s ScheduleBase) GetBatch() int64 {
	return s.Batch
}
func (s ScheduleBase) GetTimeout() time.Duration {
	return time.Duration(s.Timeout) * time.Second
}
