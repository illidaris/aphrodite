package table2struct

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

// Table2Obj 将二维字符串数组rows转换为指定dst类型的切片。opts为转换选项。
func Table2Objs(dst interface{}, rows [][]string, opts ...Table2StructOptionFunc) (err error) {
	// 初始化转换选项
	var (
		option  = newTable2StructOption(opts...)
		headMap = map[string]int{}
		header  = []string{}
	)
	dataValue := reflect.ValueOf(dst)
	// 验证dst是否为有效的指针类型和切片类型
	if dataValue.Kind() != reflect.Ptr || dataValue.Elem().Kind() != reflect.Slice {
		err = ErrTable2StructInvalid
		return
	}
	dataType := dataValue.Elem().Type().Elem()
	dataTypeKind := dataType.Kind()
	if dataTypeKind == reflect.Pointer {
		dataType = dataType.Elem()
	}
	// 遍历rows进行类型转换
	for rowIndex, row := range rows {
		// 处理表头
		if rowIndex == option.HeadIndex {
			for index, h := range row {
				if h != "" {
					headMap[h] = index
				}
			}
			header = row
		}
		// 跳过起始行之前的数据
		if rowIndex < option.StartRowIndex {
			continue
		}
		// 跳过无效行
		if len(row) >= len(headMap) && strings.Join(row[:len(headMap)], "") == "" {
			continue
		}
		// 达到行数限制时停止转换
		if option.Limit > 0 && rowIndex > option.Limit+option.StartRowIndex-1 {
			break
		}
		// 创建新的结构体实例
		newData := reflect.New(dataType).Elem()
		for _, v := range header {
			// 列不够，忽略
			colIndex := headMap[v]
			if colIndex > len(row)-1 {
				continue
			}
			cellValue := row[colIndex]
			cellValue = option.ValueConvert(v, cellValue)
			target := newData
			fieldNames := strings.Split(v, ".")
			if !option.Deep && len(fieldNames) > 1 {
				continue
			}
			for i, fieldName := range fieldNames {
				if fieldName == "" {
					return fmt.Errorf("field path:%s is invalid", v)
				}
				target = target.FieldByName(fieldName)
				if i == len(fieldNames)-1 {
					break
				}
				if target.Kind() == reflect.Pointer {
					// If the structure pointer is nil, create it.
					if target.IsNil() {
						target.Set(reflect.New(target.Type().Elem()).Elem().Addr())
					}
					target = reflect.ValueOf(target.Interface()).Elem()
				}
				if target.Kind() != reflect.Struct {
					return fmt.Errorf("field: %s is not struct", fieldName)
				}
			}
			SetValue(&target, target.Type(), cellValue)
		}
		// 将转换后的结构体实例添加到目标切片中
		if dataTypeKind == reflect.Pointer {
			dataValue.Elem().Set(reflect.Append(dataValue.Elem(), newData.Addr()))
		} else {
			dataValue.Elem().Set(reflect.Append(dataValue.Elem(), newData))
		}

	}
	return
}

func Objs2Table(dsts []interface{}, opts ...Table2StructOptionFunc) ([][]string, [][]string, error) {
	var (
		option  = newTable2StructOption(opts...)
		annoes  = []string{}
		heads   = []string{}
		headers = [][]string{}
		rows    = [][]string{}
	)
	for rowIndex, dst := range dsts {
		row := []string{}
		descs, err := Fields(dst, option.Deep, "")
		if err != nil {
			return headers, rows, err
		}
		for _, desc := range descs {
			if desc.IsExtend {
				continue
			}
			tag := desc.Id
			// 没有打标记就跳过
			if tag == "" || !option.FieldAllow(tag) {
				continue
			}
			if rowIndex == 0 {
				heads = append(heads, tag)
			}
			anno := desc.Tag.Get(option.AnnoTag)
			if rowIndex == 0 {
				comment := option.ParseAnno(tag, anno)
				annoes = append(annoes, comment)
			}
			valStr := option.ValueConvert(tag, cast.ToString(desc.Val))
			if option.IgnoreZero && valStr == "0" {
				valStr = ""
			}
			row = append(row, valStr)
		}
		rows = append(rows, row)
	}
	if len(annoes) > 0 && len(strings.Join(annoes, "")) > 0 {
		headers = append(headers, annoes)
	}
	if len(heads) > 0 {
		headers = append(headers, heads)
	}
	return headers, rows, nil
}
