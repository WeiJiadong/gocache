// Package gocache 用于并发场景，支持自动化LRU，支持singlefight更新缓存
package gocache

import (
	"context"
	"time"

	"golang.org/x/sync/singleflight"
)

// GoCacheBuilder 缓存功能基础结构定义
type GoCacheBuilder struct {
	opts  *GoCacheBuilderOpt
	cache *Cache
	stat  *cacheStat
	sf    *singleflight.Group
}

// GoCacheBuilderOpt cache相关配置
type GoCacheBuilderOpt struct {
	expire       time.Duration // 过期时间
	keyCnt       int           // key数量，默认1个key
	name         string        // cache名
	statInterval time.Duration // cache状态打印周期
}

// UpdateHelper 更新功能helper
type UpdateHelper func() (val interface{}, err error)

// GoCacheBuilderOptHelper GoCacheBuilderOpt结构helper
type GoCacheBuilderOptHelper func(opts *GoCacheBuilderOpt)

// WithKeyCnt 设置key数量上限
func WithKeyCnt(keyCnt int) GoCacheBuilderOptHelper {
	return func(opts *GoCacheBuilderOpt) {
		opts.keyCnt = keyCnt
	}
}

// WithExpire 设置key过期时间
func WithExpire(expire time.Duration) GoCacheBuilderOptHelper {
	return func(opts *GoCacheBuilderOpt) {
		opts.expire = expire
	}
}

// WithName 设置更新函数
func WithName(name string) GoCacheBuilderOptHelper {
	return func(opts *GoCacheBuilderOpt) {
		opts.name = name
	}
}

// Get 获取val，若val存在且未过期，则更新至缓存，否则通过singleflight的方式更新至缓存
// 若更新失败，则使用旧数据兜底并返回对应的error给调用方
func (c *GoCacheBuilder) Get(ctx context.Context, key string, fn UpdateHelper) (interface{}, error) {
	// 1 先直接读缓存，若数据存在则直接返回，并且上报命中
	oldVal, err := c.cache.Get(key)
	if err == nil {
		c.stat.IncrementHit()
		return oldVal, nil
	}
	// 2 若数据不存在，或者数据过期，则将数据更新到缓存里，并且上报miss
	newVal, err, _ := c.sf.Do(key, func() (interface{}, error) {
		// double check 从缓存里再获取一把数据，尽量避免重复更新
		oldVal, err := c.cache.Get(key)
		if err == nil {
			c.stat.IncrementHit()
			return oldVal, nil
		}
		// 更新数据，若数据源报错，则返回旧数据和error，若正常则更新缓存并返回数据
		c.stat.IncrementMiss()
		val, err := fn()
		if err != nil {
			return oldVal, err
		}
		return val, c.cache.Set(key, val)
	})
	return newVal, err
}

// NewGoCacheBuilder GoCacheBuilder的构造函数
func NewGoCacheBuilder(opts ...GoCacheBuilderOptHelper) *GoCacheBuilder {
	// 1 构造默认的cache对象
	c := &GoCacheBuilder{
		opts: &GoCacheBuilderOpt{
			expire:       defaultExpire,
			keyCnt:       defaultKeyNum,
			name:         defaultCacheName,
			statInterval: defaultStatInterval,
		},
		sf: &singleflight.Group{},
	}
	// 2 使用传入的opt修改cache对象的相关参数，并返回创建好的cache对象
	for i := range opts {
		opts[i](c.opts)
	}
	// 3 根据参数填充cache对象和stat对象
	c.cache = newCache(withExpire(c.opts.expire), withKeyCnt(c.opts.keyCnt))
	c.stat = newCacheStat(c.opts.name, c.cache.Len)
	return c
}
