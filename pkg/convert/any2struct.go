package convert

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// ConvertToStructByJson 将map转换为结构体 (性能比较差)
func ConvertToStructByJson(v interface{}, targetType reflect.Type) (interface{}, error) {
	bs, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	ptr := reflect.New(targetType).Interface()
	err = json.Unmarshal(bs, ptr)
	return ptr, err
}

// ConvertToStructByRef 将map转换为结构体 【尚未完成】
func ConvertToStructByRef(v interface{}, targetType reflect.Type) (interface{}, error) {
	val := reflect.ValueOf(v)
	kind := val.Kind()
	if kind == reflect.Map {
		result := reflect.New(targetType).Elem() // 创建一个新的结构体实例
		for _, key := range val.MapKeys() {
			field, _ := targetType.FieldByName(key.String()) // 获取结构体字段信息
			if field.IsExported() {                          // 确保字段是可导出的（即大写字母开头）
				fieldValue := val.MapIndex(key)                    // 获取字段值
				resultField := result.FieldByName(field.Name)      // 获取结构体中的字段值对应的反射值对象
				if resultField.IsValid() && resultField.CanSet() { // 确保可以设置该字段值
					resultField.Set(reflect.ValueOf(fieldValue.Interface())) // 设置字段值
				} else {
					return nil, fmt.Errorf("field %s not found or not settable", field.Name)
				}
			} else {
				return nil, fmt.Errorf("field %s is not valid or not exported", key.String())
			}
		}
		return result.Interface(), nil // 返回结构体实例的接口类型值
	} else {
		return nil, fmt.Errorf("value is not a map") // 不是map类型，无法转换
	}
}
