package dependency

import "time"

type IEntity interface {
	IPo
}

type ITask interface {
	IEntity
	GetTimeout() time.Duration
	GetBizId() uint64
	GetCategory() uint32
	GetName() string
}
