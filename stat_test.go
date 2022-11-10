package gocache

import (
	"sync/atomic"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStatLoop(t *testing.T) {
	expectVal, len := uint64(0), 3
	stat := newCacheStat("test_cache", func() int {
		return len
	})
	Convey("数据上报验证:", t, func() {
		stat.IncrementHit()
		stat.IncrementMiss()
		stat.statInterval = time.Millisecond
		go stat.StatLoop()
		time.Sleep(time.Second)
		So(expectVal == atomic.LoadUint64(&stat.hit), ShouldBeTrue)
		So(expectVal == atomic.LoadUint64(&stat.miss), ShouldBeTrue)
	})
}
