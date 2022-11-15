package gocache

import (
	"container/list"
	"sync"
	"time"
)

// Cache 缓存核心结构
type Cache struct {
	lru     *list.List                    // 使用list，方便做lru
	data    map[interface{}]*list.Element // 数据域，方便查询（O1）
	barrier *sync.RWMutex                 // 保证并发安全
	opts    *CacheOpt                     // 相关参数选项
}

// CacheOpt Cache参数选项结构
type CacheOpt struct {
	expire time.Duration // 过期时间，使用 time.Duration方便复用time的封装
	keyCnt int           // key数量上限
}

// CacheOptHelper Cache参数选项结构helper
type CacheOptHelper func(opt *CacheOpt)

// elem 缓存元素
type elem struct {
	key    interface{} // 缓存key
	value  interface{} // 缓存value
	expire int64       // 过期时间，使用int64方便比较
}

// Get 获取数据，返回数据和对应的error。如果数据正常，则error返回nil
func (c *Cache) Get(key interface{}) (interface{}, error) {
	c.barrier.RLock()
	defer c.barrier.RUnlock()
	// 1 查询key是否存在，若key不存在，则直接返回
	val, found := c.data[key]
	if !found {
		return nil, ErrKeyNotFound
	}
	// 2 key存在，判断key是否过期，若过期则返回val和ErrKeyIsExpired，方便调用方处理
	if val.Value.(*elem).expire < time.Now().UnixNano() {
		return val.Value.(*elem).value, ErrKeyIsExpired
	}
	return val.Value.(*elem).value, nil
}

// Set 设置数据，设置数据。用error是否为nil来标识是否出错
func (c *Cache) Set(key, val interface{}) error {
	c.barrier.Lock()
	defer c.barrier.Unlock()
	// 1 判断key是否存在，若存在则直接移到最前面，更新该值，并返回nil
	if oldVal, ok := c.data[key]; ok {
		c.lru.MoveToFront(oldVal)
		oldVal.Value.(*elem).value = val
		oldVal.Value.(*elem).expire = time.Now().Add(c.opts.expire).UnixNano()
		return nil
	}
	// 2 先判断是否到达容量限制，若达到限制，则进行LRU剔除
	if c.opts.keyCnt > 0 && c.lru.Len() >= c.opts.keyCnt {
		if lruVal := c.lru.Back(); lruVal != nil {
			c.lru.Remove(lruVal)
			v := lruVal.Value.(*elem)
			delete(c.data, v.key)
		}
	}
	// 3 key不存在，则直接在最前面插入一个值，并更新map
	c.data[key] = c.lru.PushFront(&elem{key: key, value: val, expire: time.Now().Add(c.opts.expire).UnixNano()})
	return nil
}

// UpdateCallback 更新数据回调
type UpdateCallback func() (interface{}, error)

// Len 获取缓存key的数量
func (c *Cache) Len() int {
	return c.lru.Len()
}

// withExpire 设置超时时间
func withExpire(expire time.Duration) CacheOptHelper {
	return func(opt *CacheOpt) {
		opt.expire = expire
	}
}

// withKeyCnt 设置Key上限
func withKeyCnt(keyCnt int) CacheOptHelper {
	return func(opt *CacheOpt) {
		opt.keyCnt = keyCnt
	}
}

// newCache Cache构造函数
func newCache(opts ...CacheOptHelper) *Cache {
	// 1 构造默认的cache对象
	cache := &Cache{
		data:    make(map[interface{}]*list.Element),
		lru:     list.New(),
		barrier: new(sync.RWMutex),
		opts:    new(CacheOpt),
	}
	// 2 使用传入的opt修改cache对象的相关参数，并返回创建好的cache对象
	for i := range opts {
		opts[i](cache.opts)
	}
	return cache
}
