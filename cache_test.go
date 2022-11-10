package gocache

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLen(t *testing.T) {
	expectVal, size := 0, 3
	Convey("验证len是否正确:", t, func() {
		cache := NewGoCacheBuilder(WithExpire(time.Second), WithKeyCnt(size))
		So(expectVal == cache.cache.Len(), ShouldBeTrue)
	})
}

func TestGet(t *testing.T) {
	key := "1"
	expectErr := ErrKeyNotFound
	size := 3
	Convey("验证不存在key:", t, func() {
		cache := NewGoCacheBuilder(WithExpire(time.Second), WithKeyCnt(size))
		val, err := cache.cache.Get(key)
		So(val == nil, ShouldBeTrue)
		So(expectErr == err, ShouldBeTrue)
	})
}

func TestSet(t *testing.T) {
	keys := []string{"1", "2", "3", "1", "4"}
	expectVals := []interface{}{"4", "1", "3"}
	size := 3
	Convey("验证LRU:", t, func() {
		cache := NewGoCacheBuilder(WithExpire(time.Second), WithKeyCnt(size))
		for i := range keys {
			cache.cache.Set(keys[i], keys[i])
		}
		ok, i := true, 0
		for ptr := cache.cache.lru.Front(); ptr != nil; ptr = ptr.Next() {
			if ptr.Value.(*elem).value != expectVals[i] {
				fmt.Println(ptr.Value.(*elem).value, expectVals[i], i)
				ok = false
				break
			}
			i++
		}
		So(ok, ShouldBeTrue)
	})
}
