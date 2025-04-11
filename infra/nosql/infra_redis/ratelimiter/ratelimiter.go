package ratelimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// AlgorithmType 限流算法类型
type AlgorithmType = int

const (
	TokenBucket AlgorithmType = iota // 令牌桶算法
	LeakyBucket                      // 漏桶算法
)

// Option 配置选项
type Option = func(*RateLimiter)

// RateLimiter 限流器
type RateLimiter struct {
	client    *redis.Client
	keyPrefix string // 可为空, 非空时会加":"+key"
	rate      int64  // 速率(请求/秒)
	capacity  int64  // 桶容量

	algorithm   AlgorithmType // 限流算法
	redisScript *redis.Script // Lua脚本
}

// WithAlgorithm 设置限流算法
func WithAlgorithm(algorithm AlgorithmType) Option {
	return func(r *RateLimiter) { r.algorithm = algorithm }
}

// WithRate 设置速率(请求/秒)
func WithRate(rate int64) Option { return func(r *RateLimiter) { r.rate = rate } }

// WithCapacity 设置桶容量
func WithCapacity(capacity int64) Option { return func(r *RateLimiter) { r.capacity = capacity } }

// NewRateLimiter 创建限流器
func NewRateLimiter(client *redis.Client, keyPrefix string, opts ...Option) *RateLimiter {
	r := &RateLimiter{
		client:    client,
		keyPrefix: keyPrefix,
		algorithm: TokenBucket, // 默认令牌桶算法
		rate:      10,          // 默认10请求/秒
		capacity:  20,          // 默认桶容量20
	}

	for _, opt := range opts {
		opt(r)
	}

	// 加载Lua脚本
	r.loadScripts()

	return r
}

// loadScripts 加载Lua脚本
func (r *RateLimiter) loadScripts() {
	switch r.algorithm {
	case TokenBucket:
		r.redisScript = redis.NewScript(tokenBucketScript)
	case LeakyBucket:
		r.redisScript = redis.NewScript(leakyBucketScript)
	}
}

const tokenBucketScript = `
local key = KEYS[1]
local now = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local capacity = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])

local lastTime = tonumber(redis.call("HGET", key, "lastTime")) or now
local tokens = tonumber(redis.call("HGET", key, "tokens")) or capacity

-- 计算新增的令牌数
local elapsed = now - lastTime
local newTokens = math.floor(elapsed * rate / 1000)
tokens = math.min(tokens + newTokens, capacity)

-- 检查是否有足够的令牌
local allowed = tokens >= requested
if allowed then
    tokens = tokens - requested
    lastTime = now
end

-- 更新Redis
redis.call("HMSET", key, "lastTime", lastTime, "tokens", tokens)
redis.call("EXPIRE", key, math.ceil(capacity / rate) + 1)

return allowed and 1 or 0
`

// leakyBucketScript water: 指漏斗里面的水, 即剩余未处理的请求量
const leakyBucketScript = `
local key = KEYS[1]

local now = tonumber(ARGV[1]) -- 毫秒值
local rate = tonumber(ARGV[2])
local capacity = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])

local jan_1_2025 = 1735689600
if now <= 0 then
    now = redis.call("TIME")
    now = ((now[1] - jan_1_2025) * 1000) + (now[2] / 1000)
end

local bucket_state = redis.call("HMGET", key, "lastTime", "water")
local lastTime = tonumber(bucket_state[1]) or 0
local water = tonumber(bucket_state[2]) or 0

-- 计算漏出的水量
local elapsed = now - lastTime
local leaked = math.floor(elapsed * rate / 1000)
-- 漏桶内剩余的水量, 最小为0
water = math.max(water - leaked, 0)

-- 检查是否有足够的容量(有没有超过最大请求量)
local new_water = water + requested
if new_water <= capacity then
    redis.call("HMSET", key, "lastTime", now, "water", new_water)
    redis.call("EXPIRE", key, math.ceil(capacity / rate))
    -- return {1, water, 0}
    return 1
else
    local retry_after = (new_water - capacity) / rate
    -- return {0, water, retry_after}
    return 0
end
`

// UpdateRate 动态更新速率
func (r *RateLimiter) UpdateRate(newRate int64) {
	r.rate = newRate
}

// UpdateCapacity 动态更新容量
func (r *RateLimiter) UpdateCapacity(newCapacity int64) {
	r.capacity = newCapacity
}

// Allow 检查是否允许请求
func (r *RateLimiter) Allow(ctx context.Context, identifier string) (bool, error) {
	return r.AllowN(ctx, identifier, 1)
}

// AllowN 检查是否允许N个请求
func (r *RateLimiter) AllowN(ctx context.Context, identifier string, n int64) (bool, error) {
	key, now := key(r.keyPrefix, identifier), time.Now().UnixMilli()

	result, err := r.redisScript.Run(ctx, r.client, []string{key}, now, r.rate, r.capacity, n).Int()
	if err != nil {
		return false, err
	}

	return result == 1, nil
}

// Wait 等待直到允许请求
func (r *RateLimiter) Wait(ctx context.Context, identifier string, timeout time.Duration) error {
	for {
		allowed, err := r.Allow(ctx, identifier)
		if err != nil {
			return err
		}
		if allowed {
			return nil
		}

		if timeout <= 0 {
			timeout = time.Second / time.Duration(r.rate)
		}
		select {
		case <-time.After(timeout):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// GetStatus 获取当前限流状态
func (r *RateLimiter) GetStatus(ctx context.Context, identifier string) (map[string]any, error) {
	key := key(r.keyPrefix, identifier)

	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	status := make(map[string]any)
	for k, v := range result {
		status[k] = v
	}

	return status, nil
}

func key(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + ":" + key
}
