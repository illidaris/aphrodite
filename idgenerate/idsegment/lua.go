package idsegment

import (
	"errors"
	"strings"

	"github.com/spf13/cast"
)

type StatusCode int32

const (
	StatusCodeNil            = -iota // 成功
	StatusCodeUnknown                // 未知错误
	StatusCodeHUninit                // 没有初始化
	StatusCodeOverflow               // 超出上限
	StatusCodeBadParam               // 参数错误
	StatusCodeCacheIdUpdated         // 缓存中的Id已经更新过
)

func (i StatusCode) ToError() error {
	switch i {
	case StatusCodeNil:
		return nil
	case StatusCodeHUninit:
		return errors.New("idgenerate: id segment has not been initialized")
	case StatusCodeOverflow:
		return errors.New("idgenerate: id segment overflow")
	case StatusCodeBadParam:
		return errors.New("idgenerate: bad param")
	case StatusCodeCacheIdUpdated:
		return errors.New("idgenerate: cache id has been updated")
	default:
		return errors.New("idgenerate: unknown error")
	}
}

func parseLuaResult(res interface{}) *Segment {
	seg := &Segment{}
	resultStr := cast.ToString(res)
	words := strings.Split(resultStr, "|")
	for index, word := range words {
		switch index {
		case 0:
			seg.Code = StatusCode(cast.ToInt32(word))
		case 1:
			seg.MinId = cast.ToInt64(word)
		case 2:
			seg.MaxId = cast.ToInt64(word)
		case 3:
			seg.Cursor = cast.ToInt64(word)
		}
	}
	if len(words) < 1 {
		seg.Code = StatusCodeUnknown
		return seg
	}
	return seg
}

// LUASCRIPT_HINCR 带判断的自增脚本
// -1则未知错误,-2则为没有初始化，-3则超出上限，否则返回`起始值|上一个值|最大值|当前值`。
const LUASCRIPT_HINCR = `
local maxv = redis.call('HGET',KEYS[1],'max')
if maxv then else return '-2' end
local prev = redis.call('HGET',KEYS[1],'cur') 
if prev then else return '-2' end
local step
local code = 0
if ARGV[1] + tonumber(prev) < tonumber(maxv)
then
step = tonumber(ARGV[1])
else
step = tonumber(maxv) - tonumber(prev) 
end
if step < 1 
then
return string.format("-3|%d|%d|%d",tonumber(prev),tonumber(maxv),tonumber(prev))
end
local temp = redis.call('HINCRBY',KEYS[1],'cur',step)
if temp
then
return string.format("0|%d|%d|%d",tonumber(prev),tonumber(maxv),tonumber(temp))
else 
return string.format("-1|%d|%d|%d",tonumber(prev),tonumber(maxv),tonumber(prev))
end`

// LUASCRIPT_HREPL 重新设置的上限与当前值
// -1则未知错误,-2则为没有初始化,-4则参数错误，否则返回当前值。
const LUASCRIPT_HREPL = `
if tonumber(ARGV[1]) >= tonumber(ARGV[2]) then return '-4' end
local maxv = redis.call('hget',KEYS[1],'max')
local curv = redis.call('hget',KEYS[1],'cur') 
if ((not maxv) or (not curv))
then
local res = redis.call('HMSET',KEYS[1],'cur',ARGV[1],'max',ARGV[2])
if res then return string.format("0|%d|%d",tonumber(ARGV[1]),tonumber(ARGV[2])) else return '-1' end
end
if (tonumber(maxv) > tonumber(ARGV[1])) then return '-5' end
local res = redis.call('HMSET',KEYS[1],'cur',ARGV[1],'max',ARGV[2])
if res then return string.format("0|%d|%d",tonumber(ARGV[1]),tonumber(ARGV[2])) else return '-1' end`
