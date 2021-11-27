// Package gocache 并发安全、支持lru、支持过期的缓存轮子
package gocache

import (
	"fmt"
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
			getErrs = append(getErrs, KeyNotFoundErr)
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
			getErrs = append(getErrs, KeyNotFoundErr)
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
		expectErr := KeyIsExpiredErr
		val, err := cache.Get(1)
		So(interfaceEqual(expectVal, val), ShouldBeTrue)
		So(errorEqual(expectErr, err), ShouldBeTrue)
	})
}
