package table2struct

import "strings"

// table2StructOption定义了Table2Struct转换的配置选项
type table2StructOption struct {
	StructTag        string                         // 结构体标签，默认为"json" ， 【新版不支持字段别名】
	AllowTagFields   []string                       // 允许导入或者导出的标签字段，不设置表示无限制
	IgnoreTagFields  []string                       // 忽略的字段列表
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
	Deep             bool                           // 是否深度遍历
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
	// 忽略
	if len(o.IgnoreTagFields) > 0 {
		for _, v := range o.IgnoreTagFields {
			if v == field {
				return false
			}
		}
	}
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
		IgnoreTagFields:  []string{},
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

// WithIgnoreTagFields 忽略导入或者导出的数据
func WithIgnoreTagFields(vs ...string) func(opt *table2StructOption) {
	return func(opt *table2StructOption) {
		opt.IgnoreTagFields = append(opt.IgnoreTagFields, vs...)
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

// WithDeep 启用深度遍历
func WithDeep() func(opt *table2StructOption) {
	return func(opt *table2StructOption) {
		opt.Deep = true
	}
}
