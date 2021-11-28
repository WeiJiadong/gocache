# gocache
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
