package gocache

import (
	"fmt"
	"sync/atomic"
	"time"
)

type cacheStat struct {
	name         string        // 缓存名
	hit          uint64        // 命中次数
	miss         uint64        // miss次数
	sizeCallback func() int    // cache大小
	statInterval time.Duration // stat 周期统计
}

type lenCallback func() int

// newCacheStat cacheStat构造函数
func newCacheStat(name string, fn lenCallback) *cacheStat {
	// 创建cacheStat对象
	st := &cacheStat{
		name:         name,
		statInterval: defaultStatInterval,
		sizeCallback: fn,
	}
	// 启动后台监控协程
	go st.StatLoop()
	return st
}

// IncrementHit 增加命中次数
func (cs *cacheStat) IncrementHit() {
	atomic.AddUint64(&cs.hit, 1)
}

// IncrementMiss 增加未命中次数
func (cs *cacheStat) IncrementMiss() {
	atomic.AddUint64(&cs.miss, 1)
}

// StatLoop stat statInterval一个周期循环
func (cs *cacheStat) StatLoop() {
	// 1 创建一个指定周期的定时器
	ticker := time.NewTicker(cs.statInterval)
	defer ticker.Stop()
	// 2 周期执行指标打印
	for range ticker.C {
		// 获取指标值，并原子更新指标变量
		hit := atomic.SwapUint64(&cs.hit, 0)
		miss := atomic.SwapUint64(&cs.miss, 0)
		total := hit + miss
		// 无访问就跳过打印
		if total == 0 {
			continue
		}
		percent := 100 * float32(hit) / float32(total)
		fmt.Printf("cache(%s) - qpm: %d, hit_ratio: %.1f%%, elements: %d, hit: %d, miss: %d",
			cs.name, total, percent, cs.sizeCallback(), hit, miss)
	}
}
