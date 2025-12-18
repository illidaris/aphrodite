package gormex

import (
	"hash/crc32"

	"github.com/spf13/cast"
)

type BizSection struct {
	BizId int64 `json:"bizId" gorm:"column:bizId;type:bigint;comment:游戏ID"` // game id
}

func (i BizSection) Database() string {
	return ""
}

func (i BizSection) DbSharding(keys ...any) string {
	if len(keys) == 0 {
		return cast.ToString(i.BizId)
	}
	return cast.ToString(keys[0])
}

type OperationSection struct {
	Status   int32  `json:"status" gorm:"column:status;type:int;default:1;comment:状态"`                   // 状态
	CreateBy int64  `json:"createBy" gorm:"column:createBy;<-:create;index;type:bigint;comment:创建者"`     // 创建者
	CreateAt int64  `json:"createAt" gorm:"column:createAt;<-:create;index;autoCreateTime;comment:创建时间"` // 创建时间
	UpdateBy int64  `json:"updateBy" gorm:"column:updateBy;type:bigint;comment:修改者"`                     // 修改者
	UpdateAt int64  `json:"updateAt" gorm:"column:updateAt;index;autoUpdateTime;comment:修改时间"`           // 修改时间
	Describe string `json:"describe" gorm:"column:describe;type:varchar(255);comment:描述"`                // 描述
}

func TableIndex(key any, num uint32) uint32 {
	var id uint32
	switch k := key.(type) {
	case string:
		id = HashKey(k)
	case int:
		id = uint32(k)
	default:
		id = cast.ToUint32(k)
	}
	return id % num
}

func HashKey(key string) uint32 {
	if len(key) < 64 {
		//声明一个数组长度为64
		var srcatch [64]byte
		//拷贝数据到数组中
		copy(srcatch[:], key)
		//使用IEEE 多项式返回数据的CRC-32校验和
		return crc32.ChecksumIEEE(srcatch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}
