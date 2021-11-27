package gocache

import "time"

// WithExpire 设置超时时间
func WithExpire(expire time.Duration) CacheOptHelper {
	return func(opt *CacheOpt) {
		opt.expire = expire
	}
}

// WithKeyCnt 设置Key上限
func WithKeyCnt(keyCnt int) CacheOptHelper {
	return func(opt *CacheOpt) {
		opt.keyCnt = keyCnt
	}
}
