package table2struct

// 包导入
import (
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cast"
)

// Deprecated: Table2Objs 代替 Table2Struct将二维字符串数组rows转换为指定dst类型的切片。opts为转换选项。
func Table2Struct(dst interface{}, rows [][]string, opts ...Table2StructOptionFunc) (err error) {
	// 初始化转换选项
	var (
		option  = newTable2StructOption(opts...)
		headMap = map[string]int{}
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
		// 遍历结构体字段进行赋值
		for i := 0; i < dataType.NumField(); i++ {
			field := dataType.Field(i)
			tag := field.Tag.Get(option.StructTag)
			// 如果字段没有指定的标签，则跳过
			if tag == "" || !option.FieldAllow(tag) {
				continue
			}
			colIndex, ok := headMap[tag]
			if !ok {
				continue
			}
			// 列不够，忽略
			if colIndex > len(row)-1 {
				continue
			}
			cellValue := row[colIndex]
			cellValue = option.ValueConvert(tag, cellValue)
			// 根据字段类型转换并赋值
			switch field.Type.Kind() {
			case reflect.Bool:
				v, _ := strconv.ParseBool(cellValue)
				newData.Field(i).SetBool(v)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				v, _ := strconv.ParseInt(cellValue, 10, 64)
				newData.Field(i).SetInt(v)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				v, _ := strconv.ParseUint(cellValue, 10, 64)
				newData.Field(i).SetUint(v)
			case reflect.Float32, reflect.Float64:
				v, _ := strconv.ParseFloat(cellValue, 64)
				newData.Field(i).SetFloat(v)
			case reflect.String:
				v := cellValue
				newData.Field(i).SetString(v)
			default:
				// 处理特殊类型，如time.Time和time.Duration
				switch field.Type.String() {
				case "time.Time":
					v := cast.ToTime(cellValue)
					newData.Field(i).Set(reflect.ValueOf(v))
				case "time.Duration":
					v := cast.ToDuration(cellValue)
					newData.Field(i).Set(reflect.ValueOf(v))
				}
			}
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

// Deprecated: Objs2Table 代替
func Struct2Table(dsts []interface{}, opts ...Table2StructOptionFunc) ([][]string, [][]string, error) {
	var (
		option  = newTable2StructOption(opts...)
		annoes  = []string{}
		heads   = []string{}
		headers = [][]string{}
		rows    = [][]string{}
	)
	for rowIndex, dst := range dsts {
		row := []string{}
		dataType := reflect.TypeOf(dst)
		dataValue := reflect.ValueOf(dst)
		if dataType.Kind() == reflect.Pointer {
			dataType = dataType.Elem()
			dataValue = dataValue.Elem()
		}
		for i := 0; i < dataType.NumField(); i++ {
			field := dataType.Field(i)
			tag := field.Tag.Get(option.StructTag)
			// 没有打标记就跳过
			if tag == "" || !option.FieldAllow(tag) {
				continue
			}
			if rowIndex == 0 {
				heads = append(heads, tag)
			}
			anno := field.Tag.Get(option.AnnoTag)
			if rowIndex == 0 {
				comment := option.ParseAnno(tag, anno)
				annoes = append(annoes, comment)
			}
			val := dataValue.Field(i).Interface()
			valStr := option.ValueConvert(tag, val)
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
