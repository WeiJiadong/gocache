# gocache
[![GoDoc](https://godoc.org/github.com/WeiJiadong/gocache?status.svg)](https://pkg.go.dev/github.com/WeiJiadong/gocache@v1.0.4)
[![Go Report Card](https://goreportcard.com/badge/github.com/WeiJiadong/gocache)](https://goreportcard.com/report/github.com/WeiJiadong/gocache)
[![Latest](https://img.shields.io/badge/latest-v1.0.4-blue.svg)](https://github.com/WeiJiadong/gocache/tree/v1.0.4)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
### 支持的功能
1.并发安全；  
2.LRU淘汰策略；  
3.数据过期，过期策略为懒更新;  
4.key和value支持interface;  
5.支持返回过期数据和对应error；  
6.支持singelfight方式更新下游。

### 使用示例
```go
func main() {
    cache := gocache.New(gocache.WithExpire(500*time.Second), gocache.WithKeyCnt(10))
    cache.Set(1, 1)
    fmt.Println(cache.Get(1))
}
```
