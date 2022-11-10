// Package gocache 并发安全、支持lru、支持过期的缓存轮子
package gocache

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"git.code.oa.com/NGTest/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewGoCacheBuilder(t *testing.T) {
	key := "1"
	expectVal := interface{}(1)
	val := interface{}(2)
	expectErr := ErrKeyIsExpired
	size := 3
	name := "test_cache"
	Convey("验证从数据源获取:", t, func() {
		cache := NewGoCacheBuilder(WithExpire(time.Second), WithKeyCnt(size))
		v, err := cache.Get(context.TODO(), key, func() (interface{}, error) {
			return expectVal, nil
		})

		So(expectVal == v, ShouldBeTrue)
		So(err == nil, ShouldBeTrue)
	})

	Convey("验证缓存有效，不会透传数据源:", t, func() {
		cache := NewGoCacheBuilder(WithExpire(time.Second), WithKeyCnt(size), WithName(name))
		cache.cache.Set(key, expectVal)
		v, err := cache.Get(context.TODO(), key, func() (interface{}, error) {
			return val, fmt.Errorf("date error")
		})

		So(expectVal == v, ShouldBeTrue)
		So(err == nil, ShouldBeTrue)
	})

	Convey("验证数据源失败，用缓存兜底:", t, func() {
		cache := NewGoCacheBuilder(WithExpire(time.Millisecond), WithKeyCnt(size))
		cache.cache.Set(key, expectVal)
		time.Sleep(time.Millisecond)

		v, err := cache.Get(context.TODO(), key, func() (interface{}, error) {
			return val, expectErr
		})

		So(expectVal == v, ShouldBeTrue)
		So(expectErr == err, ShouldBeTrue)
	})

	Convey("验证 double check 逻辑:", t, func() {
		cache := NewGoCacheBuilder(WithExpire(time.Millisecond), WithKeyCnt(size))
		outputs := []gomonkey.OutputCell{
			{Values: gomonkey.Params{expectVal, ErrKeyIsExpired}},
			{Values: gomonkey.Params{expectVal, nil}},
		}
		gomonkey.ApplyMethodSeq(reflect.TypeOf(cache.cache), "Get", outputs)
		v, err := cache.Get(context.TODO(), key, func() (interface{}, error) {
			return val, nil
		})

		So(expectVal == v, ShouldBeTrue)
		So(err == nil, ShouldBeTrue)
	})
}

func BenchmarkGet(b *testing.B) {
	c := NewGoCacheBuilder(WithName("test_cache"), WithKeyCnt(20), WithExpire(100000))
	for i := 0; i < 1000000; i++ {
		c.cache.Set(strconv.Itoa(i), i)
	}
	b.ResetTimer()
	Convey("直接返回看get-put性能", b, func() {
		c.Get(context.TODO(), "1", func() (interface{}, error) {
			return 1, nil
		})
	})
}
