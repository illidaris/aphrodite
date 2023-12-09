package po

import (
	"encoding/json"
	"time"

	"github.com/illidaris/aphrodite/pkg/dependency"
)

var _ = dependency.ITask(&TaskQueueMessage{})

type TaskQueueMessage struct {
	dependency.EmptyPo
	IDAutoSection `gorm:"embedded"`
	RawBizSection `gorm:"embedded"`
	Category      uint32 `json:"category" gorm:"column:category;type:int;index(biz);comment:类别"` // 类别 // 1-导出任务
	Name          string `json:"name" gorm:"column:name;type:varchar(36);index(biz);comment:任务"` // 业务类型
	Args          string `json:"args" gorm:"column:args;type:text;comment:参数"`
	LockSection   `gorm:"embedded"`
	CreateBy      int64 `json:"createBy" gorm:"column:createBy;<-:create;index;type:bigint;comment:创建者"`     // 创建者
	CreateAt      int64 `json:"createAt" gorm:"column:createAt;<-:create;index;autoCreateTime;comment:创建时间"` // 创建时间
	UpdateBy      int64 `json:"updateBy" gorm:"column:updateBy;type:bigint;comment:修改者"`                     // 修改者
	UpdateAt      int64 `json:"updateAt" gorm:"column:updateAt;index;autoUpdateTime;comment:修改时间"`           // 修改时间
}

func (s TaskQueueMessage) TableName() string {
	return "task_queue_message"
}

func (s TaskQueueMessage) ID() any {
	return s.Id
}

func (s TaskQueueMessage) GetTimeout() time.Duration {
	return time.Duration(s.Timeout) * time.Second
}

func (s TaskQueueMessage) GetBizId() uint64 {
	return s.BizId
}

func (s TaskQueueMessage) GetCategory() uint32 {
	return s.Category
}

func (s TaskQueueMessage) GetName() string {
	return s.Name
}

func (p TaskQueueMessage) Database() string {
	return ""
}
func (p TaskQueueMessage) ToRow() []string {
	return []string{}
}

func (p *TaskQueueMessage) ToJson() string {
	bs, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(bs)
}
