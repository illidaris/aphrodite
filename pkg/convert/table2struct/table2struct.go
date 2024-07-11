package table2struct

// 包导入
import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cast"
)

// 定义Table2Struct转换过程中的错误类型
var (
	ErrTable2StructInvalid = errors.New("无效的数据类型")
)

// table2StructOption定义了Table2Struct转换的配置选项
type table2StructOption struct {
	StructTag        string                         // 结构体标签，默认为"json"
	AllowTagFields   []string                       // 允许导入或者导出的标签字段，不设置表示无限制
	FieldConvertFunc map[string]func(string) string // 字段转换函数，默认为空，不转换
	IgnoreZero       bool                           // 是否忽略零值，默认为false，不忽略
	AnnoTag          string                         // 注释标签，默认为"gorm"
	AnnoTagSplit     string                         // 注释标签分隔符，默认为";"
	AnnoTagKey       string                         // 注释标签键，默认为"comment"
	AnnoTagKeySplit  string                         // 注释标签键分隔符，默认为":"
	AnnoMap          map[string]string              // 注释标签键值对，默认为空
	HeadIndex        int                            // 表头索引，默认为0
	StartRowIndex    int                            // 起始行索引，默认为1，即第一行数据开始
	Limit            int                            // 转换行数限制，默认为0，表示无限制

}

// ParseAnno 解析注释
func (o table2StructOption) ParseAnno(tag, anno string) string {
	comment := ""
	if len(o.AnnoMap) > 0 {
		comment = o.AnnoMap[tag]
	}
	if comment == "" {
		kvs := strings.Split(anno, o.AnnoTagSplit)
		for _, v := range kvs {
			ks := strings.Split(v, o.AnnoTagKeySplit)
			if len(ks) > 1 && ks[0] == o.AnnoTagKey {
				comment = ks[1]
			}
		}
	}
	return comment
}

// Allow 过滤字段
func (o table2StructOption) FieldAllow(field string) bool {
	if len(o.AllowTagFields) == 0 {
		return true
	}
	for _, v := range o.AllowTagFields {
		if v == field {
			return true
		}
	}
	return false
}

// ValueConvert 值转化
func (o table2StructOption) ValueConvert(field, value string) string {
	f, ok := o.FieldConvertFunc[field]
	if !ok || f == nil {
		return value
	}
	return f(value)
}

// newTable2StructOption根据提供的Table2StructOptionFuncs生成并返回table2StructOption实例
func newTable2StructOption(opts ...Table2StructOptionFunc) table2StructOption {
	opt := table2StructOption{
		StructTag:        "json",
		AllowTagFields:   []string{},
		FieldConvertFunc: map[string]func(string) string{},
		AnnoTag:          "gorm",
		AnnoTagSplit:     ";",
		AnnoTagKey:       "comment",
		AnnoTagKeySplit:  ":",
		AnnoMap:          map[string]string{},
		HeadIndex:        0,
		StartRowIndex:    1,
		Limit:            0,
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

// WithAllowTagFields 允许导入或者导出的数据
func WithAllowTagFields(vs ...string) func(opt *table2StructOption) {
	return func(opt *table2StructOption) {
		opt.AllowTagFields = append(opt.AllowTagFields, vs...)
	}
}

// WithIgnoreZero 忽略0，"0" 转 ""
func WithIgnoreZero() func(opt *table2StructOption) {
	return func(opt *table2StructOption) {
		opt.IgnoreZero = true
	}
}

// WithFieldConvertFunc 字段值转化函数
func WithFieldConvertFunc(field string, f func(string) string) func(opt *table2StructOption) {
	return func(opt *table2StructOption) {
		opt.FieldConvertFunc[field] = f
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

// WithAnnoTag 注释标签，默认为"gorm"
func WithAnnoTag(v string) func(opt *table2StructOption) {
	return func(opt *table2StructOption) {
		opt.AnnoTag = v
	}
}

// AnnoTagSplit 注释标签分隔符，默认为";"
func WithAnnoTagSplit(v string) func(opt *table2StructOption) {
	return func(opt *table2StructOption) {
		opt.AnnoTagSplit = v
	}
}

// WithAnnoTagKey 注释标签键，默认为"comment"
func WithAnnoTagKey(v string) func(opt *table2StructOption) {
	return func(opt *table2StructOption) {
		opt.AnnoTagKey = v
	}
}

// WithAnnoTagKeySplit 注释标签键分隔符，默认为":"
func WithAnnoTagKeySplit(v string) func(opt *table2StructOption) {
	return func(opt *table2StructOption) {
		opt.AnnoTagKeySplit = v
	}
}

// WithAnnoMap 注释标签键值对，默认为空
func WithAnnoMap(m map[string]string) func(opt *table2StructOption) {
	return func(opt *table2StructOption) {
		for k, v := range m {
			opt.AnnoMap[k] = v
		}
	}
}

// Table2Struct将二维字符串数组rows转换为指定dst类型的切片。opts为转换选项。
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
		dataValue.Elem().Set(reflect.Append(dataValue.Elem(), newData))
	}
	return
}

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
			valStr := option.ValueConvert(tag, cast.ToString(val))
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
