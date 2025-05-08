package limiter

import (
	"sync"
	"time"
)

// TokenBucket 是令牌桶限流器。
// 每隔一段时间往桶里添加令牌，请求拿令牌才能通过。
type TokenBucket struct {
	rate       int        // 每秒生成令牌数
	capacity   int        // 桶最大容量
	tokens     int        // 当前令牌数量
	lastRefill time.Time  // 上次生成令牌时间
	mu         sync.Mutex // 锁
}

// NewTokenBucket 创建一个令牌桶限流器。
func NewTokenBucket(rate int, capacity int) *TokenBucket {
	return &TokenBucket{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity,
		lastRefill: time.Now(),
	}
}

// Allow 判断请求是否能获取到令牌。
func (t *TokenBucket) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(t.lastRefill).Seconds()
	newTokens := int(elapsed * float64(t.rate))

	// 补充令牌
	if newTokens > 0 {
		t.tokens = min(t.capacity, t.tokens+newTokens)
		t.lastRefill = now
	}

	// 消耗令牌
	if t.tokens > 0 {
		t.tokens--
		return true
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
