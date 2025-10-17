package table2struct

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cast"
)

func SetValue(target *reflect.Value, typ reflect.Type, cellValue string) {
	// 根据字段类型转换并赋值
	switch typ.Kind() {
	case reflect.Bool:
		v, _ := strconv.ParseBool(cellValue)
		target.SetBool(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, _ := strconv.ParseInt(cellValue, 10, 64)
		target.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, _ := strconv.ParseUint(cellValue, 10, 64)
		target.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v, _ := strconv.ParseFloat(cellValue, 64)
		target.SetFloat(v)
	case reflect.String:
		v := cellValue
		target.SetString(v)
	default:
		// 处理特殊类型，如time.Time和time.Duration
		switch typ.String() {
		case "time.Time":
			v := cast.ToTime(cellValue)
			target.Set(reflect.ValueOf(v))
		case "time.Duration":
			v := cast.ToDuration(cellValue)
			target.Set(reflect.ValueOf(v))
		}
	}
}

type FieldDesc struct {
	Id  string
	Pid string
	reflect.StructField
	Val      interface{}
	IsExtend bool
}

func Fields(obj interface{}, deep bool, pid string) ([]*FieldDesc, error) {
	descs := []*FieldDesc{}
	if obj == nil {
		return nil, errors.New("obj must not be nil")
	}
	typ := Type(obj)
	val := Value(obj)

	if kind := typ.Kind(); kind != reflect.Struct {
		return nil, errors.New("obj must be struct")
	}

	for i := 0; i < typ.NumField(); i++ {
		ct := typ.Field(i)
		cf := val.Field(i)
		desc := &FieldDesc{}
		desc.StructField = ct
		desc.Pid = pid
		keys := []string{}
		if pid != "" {
			keys = append(keys, pid)
		}
		keys = append(keys, ct.Name)
		desc.Id = strings.Join(keys, ".")
		desc.Val = cf.Interface()
		descs = append(descs, desc)
		if deep {
			if cf.Kind() == reflect.Pointer && !cf.IsNil() {
				cf = cf.Elem()
			}
			if cf.Kind() == reflect.Struct {
				desc.IsExtend = true
				subDescs, err := Fields(cf, deep, desc.Id)
				if err != nil {
					return nil, err
				}
				descs = append(descs, subDescs...)
				continue
			}
		}
	}
	return descs, nil
}

func NewInstance(obj interface{}) interface{} {
	if obj == nil {
		return nil
	}
	entity := reflect.ValueOf(obj)

	switch entity.Kind() {
	case reflect.Ptr:
		entity = reflect.New(entity.Elem().Type())
	case reflect.Chan:
		entity = reflect.MakeChan(entity.Type(), entity.Cap())
	case reflect.Map:
		entity = reflect.MakeMap(entity.Type())
	case reflect.Slice:
		entity = reflect.MakeSlice(entity.Type(), 0, entity.Cap())
	default:
		entity = reflect.New(entity.Type()).Elem()
	}

	return entity.Interface()
}

func Type(obj interface{}) reflect.Type {
	if obj == nil {
		return nil
	}
	if v, ok := obj.(reflect.Type); ok {
		return v
	}
	if v, ok := obj.(reflect.Value); ok {
		return v.Type()
	}
	if reflect.TypeOf(obj).Kind() == reflect.Pointer {
		return reflect.TypeOf(obj).Elem()
	}
	return reflect.TypeOf(obj)
}

func Value(obj interface{}) reflect.Value {
	var empty reflect.Value
	if obj == nil {
		return empty
	}
	if v, ok := obj.(reflect.Value); ok {
		return v
	}
	if reflect.TypeOf(obj).Kind() == reflect.Pointer {
		return reflect.ValueOf(obj).Elem()
	}
	return reflect.ValueOf(obj)
}
