package limiter

import "time"

// LeakyBucket 是漏桶限流器。
// 请求进入桶中，按固定速率流出，超出容量则丢弃。
type LeakyBucket struct {
	rate     int           // 每秒处理速率
	capacity int           // 桶容量
	queue    chan struct{} // 请求队列（模拟桶）
}

// NewLeakyBucket 创建一个漏桶限流器。
func NewLeakyBucket(rate int, capacity int) *LeakyBucket {
	lb := &LeakyBucket{
		rate:     rate,
		capacity: capacity,
		queue:    make(chan struct{}, capacity),
	}
	// 后台启动漏水 goroutine
	go lb.leak()
	return lb
}

// leak 模拟漏桶定速出水。
func (l *LeakyBucket) leak() {
	ticker := time.NewTicker(time.Second / time.Duration(l.rate))
	for range ticker.C {
		select {
		case <-l.queue:
			// 取出一个元素，表示请求被处理
		default:
			// 桶为空
		}
	}
}

// Allow 判断是否允许当前请求加入桶中。
func (l *LeakyBucket) Allow() bool {
	select {
	case l.queue <- struct{}{}:
		return true
	default:
		return false
	}
}
