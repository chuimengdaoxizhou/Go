package limiter

import (
	"sync"
	"time"
)

// FixedWindow 是固定窗口限流器。
// 每个固定时间窗口只允许指定数量的请求。
type FixedWindow struct {
	mu        sync.Mutex    // 互斥锁，保证并发安全
	count     int           // 当前窗口内请求数量
	limit     int           // 每个窗口允许的最大请求数
	startTime time.Time     // 当前窗口开始时间
	window    time.Duration // 窗口大小
}

// NewFixedWindow 创建一个新的固定窗口限流器。
func NewFixedWindow(limit int, window time.Duration) *FixedWindow {
	return &FixedWindow{
		limit:     limit,
		window:    window,
		startTime: time.Now(),
	}
}

// Allow 判断当前请求是否允许。
func (f *FixedWindow) Allow() bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 如果当前时间超过窗口时间，重置窗口
	if time.Since(f.startTime) > f.window {
		f.startTime = time.Now()
		f.count = 0
	}

	// 判断是否超过限制
	if f.count < f.limit {
		f.count++
		return true
	}
	return false
}
