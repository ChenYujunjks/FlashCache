/*
Store 是“数据存储层”

# InMemoryStore 是 Store 的一种实现

CacheService 只是持有一个 Store，
并在调用 Store 前

- 参数校验
- 业务规则
- 错误处理
- TTL转换
*/
package service

import (
	"errors"
	"strings"
	"time"

	"github.com/ChenYujunjks/FlashCache/internal/cache"
)

var (
	ErrEmptyKey   = errors.New("key cannot be empty")
	ErrEmptyValue = errors.New("value cannot be empty")
	ErrNotFound   = errors.New("key not found")
	ErrInvalidTTL = errors.New("ttl must be greater than or equal to 0")
)

type CacheService struct { // service 说 我要用一个能 Set/Get/Delete 的东西
	store cache.Store
}

func NewCacheService(store cache.Store) *CacheService {
	return &CacheService{
		store: store,
	}
}

func (s *CacheService) Set(key string, value string, ttlSeconds int) error {
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)

	if key == "" {
		return ErrEmptyKey
	}
	if value == "" {
		return ErrEmptyValue
	}
	if ttlSeconds < 0 {
		return ErrInvalidTTL
	}

	var ttl time.Duration
	if ttlSeconds > 0 {
		ttl = time.Duration(ttlSeconds) * time.Second
	}

	return s.store.Set(key, value, ttl)
}

func (s *CacheService) Get(key string) (string, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return "", ErrEmptyKey
	}

	value, ok := s.store.Get(key)
	if !ok {
		return "", ErrNotFound
	}

	return value, nil
}

func (s *CacheService) Delete(key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return ErrEmptyKey
	}

	ok := s.store.Delete(key)
	if !ok {
		return ErrNotFound
	}

	return nil
}
