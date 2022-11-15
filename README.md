# gocache
[![Badge](https://img.shields.io/badge/link-996.icu-%23FF4D5B.svg?style=flat-square)](https://996.icu/#/zh_CN)
[![Go](https://github.com/WeiJiadong/gocache/workflows/Go/badge.svg?branch=master)](https://github.com/WeiJiadong/gocache/actions)
[![Go Report Card](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat)](https://goreportcard.com/report/github.com/WeiJiadong/gocache)
[![GoDoc](https://godoc.org/github.com/WeiJiadong/gocache?status.svg)](https://pkg.go.dev/github.com/WeiJiadong/gocache@v1.1.7)
[![Latest](https://img.shields.io/badge/latest-v1.1.2-blue.svg)](https://github.com/WeiJiadong/gocache/tree/v1.1.2)
[![codecov](https://codecov.io/gh/WeiJiadong/gocache/branch/master/graph/badge.svg?token=6RG0W91RF2)](https://codecov.io/gh/WeiJiadong/gocache)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
### 支持的功能
1.并发安全；  
2.LRU淘汰策略；  
3.数据过期，过期策略为懒更新;  
4.value支持interface;  
5.支持返回过期数据和对应error；  
6.支持singelfight方式更新下游;</br>
7.支持一键式读取并更新缓存;</br>
8.支持缓存状态信息打印。

### 使用示例
```go
func main() {
    cache := NewGoCacheBuilder(WithExpire(time.Second), WithKeyCnt(3), WithName("test_cache"))
    v, err := cache.Get(context.TODO(), "1", func() (interface{}, error) {
        return "1", nil
    })
}
```
