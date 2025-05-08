package limiter

import (
	"sync"
	"time"
)

// SlidingWindow 是滑动窗口限流器。
// 统计过去一段时间内的请求数量，避免突刺问题。
type SlidingWindow struct {
	mu         sync.Mutex    // 并发锁
	limit      int           // 时间窗口内允许的最大请求数
	interval   time.Duration // 时间窗口大小
	timestamps []time.Time   // 请求时间戳列表
}

// NewSlidingWindow 创建一个滑动窗口限流器。
func NewSlidingWindow(limit int, interval time.Duration) *SlidingWindow {
	return &SlidingWindow{
		limit:      limit,
		interval:   interval,
		timestamps: []time.Time{},
	}
}

// Allow 判断请求是否允许。
func (s *SlidingWindow) Allow() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-s.interval)

	// 移除过期的时间戳
	valid := s.timestamps[:0]
	for _, t := range s.timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	s.timestamps = valid

	if len(s.timestamps) < s.limit {
		s.timestamps = append(s.timestamps, now)
		return true
	}
	return false
}
