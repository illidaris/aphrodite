package dependency

import "time"

type IEntity interface {
	IPo
}

type ITask interface {
	IEntity
	GetTimeout() time.Duration
	GetBizId() int64
	GetCategory() int32
	GetName() string
}
