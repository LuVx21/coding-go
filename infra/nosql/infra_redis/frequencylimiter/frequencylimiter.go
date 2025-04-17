package frequencylimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const lua_script = `
-- 原子递减: 从max开始, 每次递减down(如1), 直到min
-- 限频: 如每天仅能5次, redis.lua foo 1 0 5 86400
local key = KEYS[1];

local down = tonumber(ARGV[1]); -- 扣减数量
local min = tonumber(ARGV[2]) or math.mininteger; -- 递减下界
local max = tonumber(ARGV[3]); -- 递减上界
local exp = tonumber(ARGV[4]); -- 秒

local num = tonumber(redis.call('GET', key))
if num == nil then
    -- 如果没有设置过, 设置为max
    if max == nil then
        return 0
    end
    if exp == nil or exp <= 0 then
        redis.call('SET', key, max)
    else
        -- 设置初始值+过期时间
        redis.call('SET', key, max, 'EX', exp)
    end
    num = max
end

-- 如果当前值小于等于最小值, 直接返回0(不再递减)
if num <= min then
    return 0
end

if (num - down) >= min then
    redis.call('DECRBY', key, down)
    return 1
else
    return 0
end
`

type FrequencyLimiter struct {
	keyPrefix   string // 可为空, 非空时会加":"+key"
	client      *redis.Client
	redisScript *redis.Script // Lua脚本
}

func NewFrequencyLimiter(client *redis.Client, keyPrefix string) *FrequencyLimiter {
	return &FrequencyLimiter{
		keyPrefix:   keyPrefix,
		client:      client,
		redisScript: redis.NewScript(lua_script),
	}
}

// DecrTimesInDay 一天内能递减times个
func (r *FrequencyLimiter) DecrTimesInDay(ctx context.Context, identifier string, down int32, times int64) (bool, error) {
	return r.DecrInDay(ctx, identifier, down, 0, times)
}

// DecrTimesInDay 一天内允许max递减到min
func (r *FrequencyLimiter) DecrInDay(ctx context.Context, identifier string, down int32, min, max int64) (bool, error) {
	return r.Decr(ctx, identifier, down, min, max, time.Hour*24)
}

// Decr exp: <=0则不过期, max~min内递减
func (r *FrequencyLimiter) Decr(ctx context.Context, identifier string, down int32, min, max int64, exp time.Duration) (bool, error) {
	key := key(r.keyPrefix, identifier)
	result, err := r.redisScript.Run(ctx, r.client, []string{key}, down, min, max, exp.Seconds()).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}

func key(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + ":" + key
}
