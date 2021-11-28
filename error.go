package gocache

import "fmt"

var (
	// ErrKeyNotFound key不存在错误
	ErrKeyNotFound = fmt.Errorf("key not found")
	// ErrKeyIsExpired key过期错误
	ErrKeyIsExpired = fmt.Errorf("key is expired")
)
