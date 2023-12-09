package po

import "github.com/spf13/cast"

// IDAutoSection
type IDAutoSection struct {
	Id uint64 `json:"id" gorm:"column:id;type:bigint;primaryKey;comment:唯一ID"` // 主键ID
}

func (p IDAutoSection) ID() any {
	return p.Id
}

// RawBizSection
type RawBizSection struct {
	BizId uint64 `json:"bizId" gorm:"column:bizId;type:bigint;index:biz;comment:业务ID"` // biz id
}

func (s RawBizSection) Database() string {
	return ""
}

// BizSection
type BizSection struct {
	RawBizSection
}

func (s BizSection) DbSharding(keys ...any) string {
	if len(keys) == 0 {
		return cast.ToString(s.BizId)
	}
	return cast.ToString(keys[0])
}

// OperationSection
type OperationSection struct {
	Status   int32  `json:"status" gorm:"column:status;type:int;default:1;comment:状态"`                   // 状态 0-默认 1-未发布 2-预发布 3-发布中 4-已结束
	CreateBy int64  `json:"createBy" gorm:"column:createBy;<-:create;index;type:bigint;comment:创建者"`     // 创建者
	CreateAt int64  `json:"createAt" gorm:"column:createAt;<-:create;index;autoCreateTime;comment:创建时间"` // 创建时间
	UpdateBy int64  `json:"updateBy" gorm:"column:updateBy;type:bigint;comment:修改者"`                     // 修改者
	UpdateAt int64  `json:"updateAt" gorm:"column:updateAt;index;autoUpdateTime;comment:修改时间"`           // 修改时间
	Describe string `json:"describe" gorm:"column:describe;type:varchar(255);comment:描述"`                // 描述
}

// LockSection
type LockSection struct {
	Locker     string `json:"locaker"  gorm:"column:locker;type:varchar(36);default:0;index;comment:防竞争锁"`   // 防竞争锁
	Expire     int64  `json:"expire"  gorm:"column:expire;type:bigint;default:0;index;comment:锁有效期"`         // 锁有效期
	Timeout    int64  `json:"timeout" gorm:"column:timeout;<-:create;type:bigint;default:0;comment:任务超时（秒）"` // 任务超时时间
	LastError  string `json:"lastError"  gorm:"column:lastError;type:varchar(255);comment:最后失败原因"`           // 最后失败原因
	LastExecAt int64  `json:"lastExecAt"  gorm:"column:lastExecAt;type:bigint;comment:最后执行时间"`               // 最后执行时间
	Retries    int32  `json:"retries"  gorm:"column:retries;type:int;default:0;index;comment:执行次数"`          // 重试次数
}
