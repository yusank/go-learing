package main

import (
	"log"
	"math"
	"sync"
	"time"
)

// SlidingWindowCounter 为滑动窗口计数器
type SlidingWindowCounter struct {
	debug  bool
	bucket map[int64]float64
	lock   *sync.RWMutex
	size   int64
}

// NewSlidingWindowCounter 创建可配置大小的滑动窗口
func NewSlidingWindowCounter(size int64, debug bool) *SlidingWindowCounter {
	if size <= 0 {
		return nil
	}

	return &SlidingWindowCounter{
		debug:  debug,
		bucket: make(map[int64]float64),
		lock:   new(sync.RWMutex),
		size:   size,
	}
}

// Incr 计数
func (n *SlidingWindowCounter) Incr(i float64) {
	if i <= 0 {
		return
	}

	n.lock.Lock()
	defer n.lock.Unlock()

	now := time.Now().Unix()
	n.bucket[now] = n.bucket[now] + i
	n.removeOld(now - int64(n.size))
}

// removeOld 清理旧数据
func (n *SlidingWindowCounter) removeOld(before int64) {
	for ts := range n.bucket {
		if ts <= before {
			delete(n.bucket, ts)
		}
	}
}

// Sum 求和
func (n *SlidingWindowCounter) Sum() float64 {
	n.lock.RLock()
	defer n.lock.RUnlock()

	var s float64
	for _, score := range n.bucket {
		s += score
	}

	n.log("cur window length:", len(n.bucket))
	return s
}

// Avg 求平均
func (n *SlidingWindowCounter) Avg() float64 {
	return n.Sum() / float64(n.size)
}

// Max 返回传入时间now至 now-size 时间窗口内的最大值
// 可以传入当前时间 也可以传入相对早一点的时间（在滑动窗口范围内）
func (n *SlidingWindowCounter) Max(now time.Time) float64 {
	n.lock.RLock()
	defer n.lock.RUnlock()

	var (
		start = now.Unix() - n.size
		end   = now.Unix()
		max   float64
	)
	for ts := range n.bucket {
		// 在时间范围内且大于 max
		if ts >= start && ts <= end && n.bucket[ts] > max {
			max = n.bucket[ts]
		}
	}

	return max
}

// Min 返回传入时间now至 now-size 时间窗口内的最小值
// 可以传入当前时间 也可以传入相对早一点的时间（在滑动窗口范围内）
func (n *SlidingWindowCounter) Min(now time.Time) float64 {
	n.lock.RLock()
	defer n.lock.RUnlock()

	var (
		start = now.Unix() - n.size
		end   = now.Unix()
		min   = math.MaxFloat64
	)

	for ts := range n.bucket {
		// 在时间范围内且大于 max
		if ts >= start && ts <= end && n.bucket[ts] < min {
			min = n.bucket[ts]
		}
	}

	// 做一次特殊处理
	if min == math.MaxFloat64 {
		min = 0
	}

	return min
}

func (n *SlidingWindowCounter) log(a ...interface{}) {
	if n.debug {
		log.Println(a...)
	}
}
