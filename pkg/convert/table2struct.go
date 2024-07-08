package convert

// 包导入
import (
	"errors"
	"reflect"
	"strconv"

	"github.com/spf13/cast"
)

// 定义Table2Struct转换过程中的错误类型
var (
	ErrTable2StructInvalid = errors.New("无效的数据类型")
)

// table2StructOption定义了Table2Struct转换的配置选项
type table2StructOption struct {
	StructTag     string // 结构体标签，默认为"json"
	HeadIndex     int    // 表头索引，默认为0
	StartRowIndex int    // 起始行索引，默认为1，即第一行数据开始
	Limit         int    // 转换行数限制，默认为0，表示无限制
}

// newTable2StructOption根据提供的Table2StructOptionFuncs生成并返回table2StructOption实例
func newTable2StructOption(opts ...Table2StructOptionFunc) table2StructOption {
	opt := table2StructOption{
		StructTag:     "json",
		HeadIndex:     0,
		StartRowIndex: 1,
		Limit:         0,
	}
	for _, f := range opts {
		f(&opt)
	}
	return opt
}

// Table2StructOptionFunc定义了修改table2StructOption的函数类型
type Table2StructOptionFunc func(opt *table2StructOption)

// WithStructTag返回一个函数，用于设置table2StructOption的StructTag字段
func WithStructTag(v string) func(opt *table2StructOption) {
	return func(opt *table2StructOption) {
		opt.StructTag = v
	}
}

// WithHeadIndex返回一个函数，用于设置table2StructOption的HeadIndex字段
func WithHeadIndex(v int) func(opt *table2StructOption) {
	return func(opt *table2StructOption) {
		opt.HeadIndex = v
	}
}

// WithStartRowIndex返回一个函数，用于设置table2StructOption的StartRowIndex字段
func WithStartRowIndex(v int) func(opt *table2StructOption) {
	return func(opt *table2StructOption) {
		opt.StartRowIndex = v
	}
}

// WithLimit返回一个函数，用于设置table2StructOption的Limit字段
func WithLimit(v int) func(opt *table2StructOption) {
	return func(opt *table2StructOption) {
		opt.Limit = v
	}
}

// Table2Struct将二维字符串数组rows转换为指定dst类型的切片。opts为转换选项。
func Table2Struct(dst interface{}, rows [][]string, opts ...Table2StructOptionFunc) (err error) {
	// 初始化转换选项
	var (
		option  = newTable2StructOption(opts...)
		heads   = []string{}
		headMap = map[string]int{}
	)
	dataValue := reflect.ValueOf(dst)
	// 验证dst是否为有效的指针类型和切片类型
	if dataValue.Kind() != reflect.Ptr || dataValue.Elem().Kind() != reflect.Slice {
		err = ErrTable2StructInvalid
		return
	}
	dataType := dataValue.Elem().Type().Elem()
	// 遍历rows进行类型转换
	for rowIndex, row := range rows {
		// 处理表头
		if rowIndex == option.HeadIndex {
			heads = row
			for index, h := range heads {
				headMap[h] = index
			}
		}
		// 跳过起始行之前的数据
		if rowIndex < option.StartRowIndex {
			continue
		}
		// 达到行数限制时停止转换
		if rowIndex > option.Limit+option.StartRowIndex-1 {
			break
		}
		// 创建新的结构体实例
		newData := reflect.New(dataType).Elem()
		// 遍历结构体字段进行赋值
		for i := 0; i < dataType.NumField(); i++ {
			field := dataType.Field(i)
			tag := field.Tag.Get(option.StructTag)
			// 如果字段没有指定的标签，则跳过
			if tag == "" {
				continue
			}
			colIndex := headMap[tag]
			cellValue := row[colIndex]
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
				newData.Field(i).SetString(cellValue)
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
		dataValue.Elem().Set(reflect.Append(dataValue.Elem(), newData))
	}
	return
}