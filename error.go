package gocache

import "fmt"

var (
	KeyNotFoundErr  = fmt.Errorf("key not found")
	KeyIsExpiredErr = fmt.Errorf("key is expired")
)
