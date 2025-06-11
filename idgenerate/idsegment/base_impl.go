package idsegment

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/illidaris/aphrodite/pkg/dependency"

	"github.com/go-redis/redis/v8"
	"github.com/illidaris/aphrodite/component/gormex"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

var _ = ICache(IdRecordCache{})

type IdRecordCache struct {
	Client *redis.Client
}

func (i IdRecordCache) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	if i.Client == nil {
		return nil, errors.New("client is nil")
	}
	return i.Client.Eval(ctx, script, keys, args...).Result()
}
func (i IdRecordCache) EvalSha(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	if i.Client == nil {
		return nil, errors.New("client is nil")
	}
	return i.Client.EvalSha(ctx, script, keys, args...).Result()
}

var _ = IRepository(IdRecordRepository{})

type IdRecordRepository struct {
	gormex.BaseRepository[IdRecord]
	Options []dependency.BaseOptionFunc
}

func (i IdRecordRepository) Init() {

}

func (i IdRecordRepository) BlockNextSegment(ctx context.Context, key string, step int64, tryGenerate func() (*Segment, error)) (int64, int64, *Segment, error) {
	db := i.BuildFrmOptions(ctx, &IdRecord{}, i.Options...)
	tx := db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return 0, 0, nil, err
	}
	result := tx.Where("`key` = ?", key).Updates(map[string]interface{}{
		"value": gorm.Expr("`value` + ?", step),
	})
	if result.Error != nil {
		tx.Rollback()
		return 0, 0, nil, result.Error
	}
	if tryGenerate != nil {
		seg, err := tryGenerate()
		if err != nil {
			tx.Rollback()
			return 0, 0, nil, err
		}
		if seg != nil && seg.Code == StatusCodeNil {
			tx.Rollback()
			return 0, 0, seg, nil
		}
	}
	idRecord := &IdRecord{}
	idRecord.Key = key
	result = tx.First(idRecord)
	if result.Error != nil {
		tx.Rollback()
		return 0, 0, nil, result.Error
	}
	err := tx.Commit().Error
	return idRecord.Value - step + 1, idRecord.Value, nil, err
}

type IdRecord struct {
	Key      string `json:"key" gorm:"column:key;type:varchar(32);primaryKey;comment:业务主键"`
	BizId    int64  `json:"bizId" gorm:"column:bizId;type:bigint;comment:游戏ID"` // game id
	Value    int64  `json:"value" gorm:"column:value;type:bigint;default:0;comment:业务当前值"`
	CreateAt int64  `json:"createAt" gorm:"column:createAt;<-:create;autoCreateTime;comment:创建时间"` // 创建时间
	UpdateAt int64  `json:"updateAt" gorm:"column:updateAt;autoUpdateTime;comment:修改时间"`           // 修改时间
}

func (p IdRecord) TableName() string {
	return "id_record"
}

func (s IdRecord) ID() any {
	return s.Key
}

func (s IdRecord) Database() string {
	return ""
}

func (s IdRecord) DbSharding(keys ...any) string {
	if s.BizId > 0 {
		return cast.ToString(s.BizId)
	}
	if len(keys) == 0 {
		return ""
	}
	return cast.ToString(keys[0])
}

func (s IdRecord) ToJson() string {
	bs, err := json.Marshal(&s)
	if err != nil {
		return ""
	}
	return string(bs)
}

func (s IdRecord) ToRow() []string {
	return []string{}
}
