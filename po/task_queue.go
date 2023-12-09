package po

type TaskQueueMessage struct {
	IDAutoSection `gorm:"embedded"`
	BizSection    `gorm:"embedded"`
	Category      int32  `json:"category" gorm:"column:category;type:int;index(biz);comment:类别"` // 类别 // 1-导出任务
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
