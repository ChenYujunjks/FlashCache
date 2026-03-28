// internal/cache/store.go	存储接口
package cache

import "time"

type Store interface {
	Set(key string, value string, ttl time.Duration) error
	Get(key string) (string, bool)
	Delete(key string) bool
}
