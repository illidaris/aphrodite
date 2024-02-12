package po

import (
	"encoding/json"
	"time"

	"github.com/illidaris/aphrodite/pkg/dependency"
)

var _ = dependency.ITask(&MqMessage[any]{})

type MqMessage[T any] struct {
	dependency.EmptyPo
	IDAutoSection `gorm:"embedded"`
	RawBizSection `gorm:"embedded"`
	Category      uint32 `json:"category" gorm:"column:category;type:int;index:biz;comment:类别"` // 类别
	Topic         string `json:"name" gorm:"column:name;type:varchar(36);index:biz;comment:任务"` // 业务类型
	Key           string `json:"key" gorm:"column:key;type:varchar(36);comment:分区ID"`           // 分区ID
	Args          T      `json:"args" gorm:"column:args;type:text;serializer:json;comment:参数"`
	TraceId       string `json:"traceId"  gorm:"column:traceId;type:varchar(36);default:0;comment:追踪链路ID"` // 关联traceId
	LockSection   `gorm:"embedded"`
	CreateAt      int64 `json:"createAt" gorm:"column:createAt;<-:create;index;autoCreateTime;comment:创建时间"` // 创建时间
	UpdateAt      int64 `json:"updateAt" gorm:"column:updateAt;index;autoUpdateTime;comment:修改时间"`           // 修改时间
}

func (s MqMessage[T]) TableName() string {
	return "aphrodite_mq_compensate"
}

func (s MqMessage[T]) ID() any {
	return s.Id
}

func (s MqMessage[T]) GetTimeout() time.Duration {
	return time.Duration(s.Timeout) * time.Second
}

func (s MqMessage[T]) GetBizId() uint64 {
	return s.BizId
}

func (s MqMessage[T]) GetCategory() uint32 {
	return s.Category
}

func (s MqMessage[T]) GetName() string {
	return s.Topic
}

func (p MqMessage[T]) Database() string {
	return ""
}
func (p MqMessage[T]) ToRow() []string {
	return []string{}
}

func (p MqMessage[T]) ToJson() string {
	bs, err := json.Marshal(&p)
	if err != nil {
		return ""
	}
	return string(bs)
}
