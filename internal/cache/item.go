// internal/cache/item.go	数据单元
package cache

import "time"

type Item struct {
	Value     string
	ExpiresAt time.Time //零值 time.Time{}
	HasExpiry bool      //零值 false
}
