package gormex

import (
	"context"
	"errors"
	"reflect"

	"gorm.io/gorm/schema"
)

var _ = schema.SerializerInterface(&EncryptedSerializer{})

// EncryptedSerializer 在 GORM 读写字段时调用加解密函数。
// 空值和空字符串会按原样保存，避免对无内容的数据执行加解密。
type EncryptedSerializer struct {
	EncryptFunc func(context.Context, string) (string, error)
	DecryptFunc func(context.Context, string) (string, error)
}

// Value 将字段值转换为数据库值；非字符串类型按原有约定返回 nil。
func (s EncryptedSerializer) Value(ctx context.Context,
	field *schema.Field,
	dst reflect.Value,
	fieldValue interface{}) (interface{}, error) {
	if fieldValue == nil {
		return nil, nil
	}
	plain, ok := fieldValue.(string)
	if !ok {
		return nil, nil
	}
	if plain == "" {
		return "", nil
	}
	if s.EncryptFunc == nil {
		return nil, errors.New("gormex: EncryptFunc 未配置")
	}
	return s.EncryptFunc(ctx, plain)
}

// Scan 将数据库中的密文解密后写入目标字段。
func (s EncryptedSerializer) Scan(ctx context.Context,
	field *schema.Field,
	dst reflect.Value,
	dbValue interface{}) error {
	if dbValue == nil {
		return nil
	}
	var cipherText string
	switch v := dbValue.(type) {
	case string:
		cipherText = v
	case []byte:
		cipherText = string(v)
	default:
		return nil
	}
	if cipherText == "" {
		field.Set(ctx, dst, "")
		return nil
	}
	if s.DecryptFunc == nil {
		return errors.New("gormex: DecryptFunc 未配置")
	}
	plain, err := s.DecryptFunc(ctx, cipherText)
	if err != nil {
		return err
	}
	field.Set(ctx, dst, plain)
	return nil
}
