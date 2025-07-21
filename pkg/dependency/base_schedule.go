package dependency

import "time"

type ISchedule interface {
	GetBatch() int64
	GetTimeout() time.Duration
}
