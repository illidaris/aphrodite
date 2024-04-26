package convert

import (
	"regexp"
	"text/template"
)

// FieldFilterLevel 定义了字段过滤的级别
type FieldFilterLevel int32

// FieldFilterLevel 的常量定义
const (
	FieldFilterLevelDefault FieldFilterLevel = iota // 默认过滤级别
	FieldFilterLevelEncode                          // 编码过滤级别
	FieldFilterLevelAssign                          // 赋值过滤级别
)

const REGEXP_FIELD = `^[a-zA-Z0-9_]*$` // 定义合法字段名的正则表达式

// 默认字段过滤级别
var defaultFieldFilterLevel FieldFilterLevel

// 允许的内部字段名集合
var innerAllowedFields = map[string]struct{}{
	"id":        {},
	"createAt":  {},
	"create_at": {},
	"modifyAt":  {},
	"modify_at": {},
	"updateAt":  {},
	"update_at": {},
	"sort":      {},
}

// SetdefaultFieldFilterLevel 设置默认的字段过滤级别，并添加额外的允许字段
func SetdefaultFieldFilterLevel(level FieldFilterLevel, fields ...string) {
	defaultFieldFilterLevel = level
	AddAllowFields(fields...)
}

// DefFieldFilter 使用默认过滤级别对字段进行过滤
func DefFieldFilter(s string) (string, bool) {
	return FieldFilter(s, defaultFieldFilterLevel)
}

// FieldFilter 根据指定的过滤级别和额外允许的字段对字段进行过滤
func FieldFilter(s string, level FieldFilterLevel, fields ...string) (string, bool) {
	if !IsField(s) {
		return "", false
	}
	switch level {
	case FieldFilterLevelEncode:
		return template.HTMLEscapeString(s), true // 对字段进行HTML编码
	case FieldFilterLevelAssign:
		allows := map[string]struct{}{}
		for f := range innerAllowedFields {
			allows[f] = struct{}{}
		}
		for _, v := range fields {
			allows[v] = struct{}{}
		}
		if !IsAllowFields(s, allows) {
			return "", false
		}
	}
	return s, true
}

// IsField 检查字符串是否为合法的字段名
func IsField(s string) bool {
	ok, _ := MatchString(s, REGEXP_FIELD)
	return ok
}

// MatchString 使用正则表达式检查字符串是否匹配
func MatchString(s string, exp string) (bool, error) {
	return regexp.MatchString(exp, s)
}

// AddAllowFields 添加额外允许的字段名到集合中
func AddAllowFields(fields ...string) {
	for _, v := range fields {
		innerAllowedFields[v] = struct{}{}
	}
}

// IsAllowFields 检查字段是否在允许的字段集合中
func IsAllowFields(field string, allowMap map[string]struct{}) bool {
	_, ok := allowMap[field]
	return ok
}

// AddAllowSortField 添加额外允许的排序字段名到集合中（与 AddAllowFields 功能重叠，可能需重构）
func AddAllowSortField(fields ...string) {
	for _, v := range fields {
		innerAllowedFields[v] = struct{}{}
	}
}
