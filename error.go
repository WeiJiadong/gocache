package gocache

import "fmt"

var (
	ErrKeyNotFound  = fmt.Errorf("key not found")
	ErrKeyIsExpired = fmt.Errorf("key is expired")
)
