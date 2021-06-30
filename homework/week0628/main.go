package main

import (
	"context"
	"log"
	"math/rand"
	"time"
)

func main() {
	// 5 个格子的滑动窗口
	swc := NewSlidingWindowCounter(5, true)

	ctx, cancel := context.WithCancel(context.Background())
	// 负责随机写
	go randomAdd(ctx, swc)

	// 统计 15 秒
	var timer = time.NewTimer(15 * time.Second)

	// 每隔一秒打印计数
loop:
	for {
		select {
		case <-timer.C:
			break loop
		default:
			time.Sleep(time.Second)
		}

		now := time.Now()
		log.Printf("[%d] sum:%v avg:%v max:%v min:%v \n", now.Unix(), swc.Sum(), swc.Avg(), swc.Max(now), swc.Min(now))
	}
	cancel()
	// 停止随机写 后再打印一次
	time.Sleep(time.Second * 2)
	now := time.Now()
	log.Printf("before exist [%d]sum:%v avg:%v max:%v min:%v", now.Unix(), swc.Sum(), swc.Avg(), swc.Max(now), swc.Min(now))
}

func randomAdd(ctx context.Context, swc *SlidingWindowCounter) {
	rand.Seed(time.Now().UnixNano())
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		time.Sleep(randomDuration())
		swc.Incr(1)
	}
}

// 随机返回 0-50 毫秒
func randomDuration() time.Duration {
	// [0,50)
	return time.Millisecond * time.Duration(rand.Int63n(50))
}
