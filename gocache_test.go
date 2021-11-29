// Package gocache 并发安全、支持lru、支持过期的缓存轮子
package gocache

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func interfaceSliceEqual(sets, gets []interface{}) bool {
	for i := range sets {
		if sets[i] != gets[i] {
			fmt.Println("sets:", sets)
			fmt.Println("gets", gets)
			return false
		}
	}
	return true
}

func errorSliceEqual(sets, gets []error) bool {
	for i := range sets {
		if sets[i] != gets[i] {
			return false
		}
	}
	return true
}

func interfaceEqual(set, get interface{}) bool {
	return set == get
}

func errorEqual(set, get error) bool {
	return set == get
}

func genKey(gid int) string {
	return strconv.Itoa(gid)
}

func TestNew(t *testing.T) {
	Convey("简单的set get:", t, func() {
		cache := New(WithExpire(5*time.Second), WithKeyCnt(3))

		sets := make([]interface{}, 0, 3)
		for i := 0; i < 3; i++ {
			sets = append(sets, i)
			cache.Set(i, i)
		}
		gets := make([]interface{}, 0, 3)
		for i := 0; i < 3; i++ {
			val, err := cache.Get(i)
			if err == nil {
				gets = append(gets, val)
			}
		}

		So(interfaceSliceEqual(sets, gets), ShouldBeTrue)
	})

	Convey("lru验证新增:", t, func() {
		cache := New(WithExpire(5*time.Second), WithKeyCnt(3))
		for i := 0; i < 10; i++ {
			cache.Set(i, i)
		}
		expectErrs := make([]error, 0, 7)
		expectVals := []interface{}{7, 8, 9}
		getVals := make([]interface{}, 0, 3)
		getErrs := make([]error, 0, 7)
		for i := 0; i < 7; i++ {
			getErrs = append(getErrs, ErrKeyNotFound)
		}
		for i := 0; i < 10; i++ {
			val, err := cache.Get(i)
			if err == nil {
				getVals = append(getVals, val)
			} else {
				expectErrs = append(expectErrs, err)
			}
		}

		So(interfaceSliceEqual(expectVals, getVals), ShouldBeTrue)
		So(errorSliceEqual(expectErrs, getErrs), ShouldBeTrue)
	})

	Convey("lru验证修改:", t, func() {
		cache := New(WithExpire(5*time.Second), WithKeyCnt(3))
		for i := 0; i < 10; i++ {
			cache.Set(i, i)
		}
		cache.Set(7, 7)
		cache.Set(0, 0)
		expectErrs := make([]error, 0, 7)
		expectVals := []interface{}{0, 7, 9}
		getVals := make([]interface{}, 0, 3)
		getErrs := make([]error, 0, 7)
		for i := 0; i < 7; i++ {
			getErrs = append(getErrs, ErrKeyNotFound)
		}
		for i := 0; i < 10; i++ {
			val, err := cache.Get(i)
			if err == nil {
				getVals = append(getVals, val)
			} else {
				expectErrs = append(expectErrs, err)
			}
		}

		So(interfaceSliceEqual(expectVals, getVals), ShouldBeTrue)
		So(errorSliceEqual(expectErrs, getErrs), ShouldBeTrue)
	})

	Convey("验证过期:", t, func() {
		cache := New(WithExpire(5*time.Second), WithKeyCnt(3))
		cache.Set(1, 1)
		time.Sleep(6 * time.Second)
		expectVal := interface{}(1)
		expectErr := ErrKeyIsExpired
		val, err := cache.Get(1)

		So(interfaceEqual(expectVal, val), ShouldBeTrue)
		So(errorEqual(expectErr, err), ShouldBeTrue)
	})

	Convey("验证GetAndSet正常获取:", t, func() {
		cache := New(WithExpire(5*time.Second), WithKeyCnt(3))
		val, err := cache.GetAndSet(context.TODO(), genKey(1), func() (val interface{}, err error) {
			return 1, nil
		})
		expectVal := interface{}(1)

		So(interfaceEqual(expectVal, val), ShouldBeTrue)
		So(errorEqual(nil, err), ShouldBeTrue)
	})

	Convey("验证GetAndSet缓存有效，不会透传数据源:", t, func() {
		cache := New(WithExpire(time.Second), WithKeyCnt(3))
		cache.Set("1", 2)
		val, err := cache.GetAndSet(context.TODO(), genKey(1), func() (val interface{}, err error) {
			return 1, fmt.Errorf("date error")
		})
		expectVal := interface{}(2)

		So(interfaceEqual(expectVal, val), ShouldBeTrue)
		So(errorEqual(nil, err), ShouldBeTrue)
	})

	Convey("验证GetAndSet数据源失败，用缓存兜底:", t, func() {
		cache := New(WithExpire(time.Second), WithKeyCnt(3))
		cache.Set("1", 2)
		time.Sleep(2 * time.Second)
		expectErr := fmt.Errorf("date error")
		val, err := cache.GetAndSet(context.TODO(), genKey(1), func() (val interface{}, err error) {
			return 1, expectErr
		})
		expectVal := interface{}(2)

		So(interfaceEqual(expectVal, val), ShouldBeTrue)
		So(errorEqual(expectErr, err), ShouldBeTrue)
	})

	// Convey("验证GetAndSet double check 逻辑:", t, func() {
	// 	cache := New(WithExpire(time.Second), WithKeyCnt(3))
	// 	flag := true
	// 	get := func() (interface{}, error) {
	// 		if flag {
	// 			flag = false
	// 			return nil, ErrKeyNotFound
	// 		}
	// 		return interface{}(1), nil
	// 	}
	// 	gomonkey.ApplyFunc(cache.Get, get)
	// 	expectVal := interface{}(1)
	// 	val, err := cache.GetAndSet(context.TODO(), genKey(1), func() (val interface{}, err error) {
	// 		return 2, nil
	// 	})

	// 	So(interfaceEqual(expectVal, val), ShouldBeTrue)
	// 	So(errorEqual(nil, err), ShouldBeTrue)
	// })
}
