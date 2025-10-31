package sender

const (
	LIMITER_BASE              = "_base"
	KEY_SECTION_CODE          = "vcode"
	KEY_SECTION_LOCKED        = "locked"
	KEY_SECTION_LIMITERS_BASE = "limiters_base"
	KEY_SECTION_LIMITERS_IP   = "limiters_ip"
	KEY_SECTION_LIMITERS_UID  = "limiters_uid"
)

// LUA_ST_INC 带判断的自增脚本，-1则为超限，否则返回当前值。
// KEYS 1-计数器KEY
// ARGV 1-上限,2-步长,3-有效期(秒)
const LUA_ST_INC = `
local temp 
local affect = redis.call('setnx',KEYS[1],0)
if (tonumber(affect))>0
then
redis.call('EXPIRE',KEYS[1],ARGV[3])
end
local num = redis.call('get',KEYS[1]) 
if tonumber(ARGV[1])>=(tonumber(num)+ARGV[2]) 
then
temp = redis.call('incrby',KEYS[1],ARGV[2]) 
return tostring(temp)
end 
return '-1'`
