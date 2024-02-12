package dependency

import "time"

type IEntity interface {
	IPo
}

type IBaseTask interface {
	GetTimeout() time.Duration
	GetBizId() uint64
	GetCategory() uint32
	GetName() string
}

type ITask interface {
	IEntity
	IBaseTask
}
