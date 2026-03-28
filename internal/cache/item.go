package cache

import "time"

type Item struct {
	Value     string
	ExpiresAt time.Time
	HasExpiry bool
}
