package cache

import (
	"errors"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

type TokenCache struct {
	cache *cache.Cache
	mutex sync.Mutex
}

func NewTokenCache(defaultExpiration time.Duration, cleanupInterval time.Duration) *TokenCache {
	return &TokenCache{
		cache: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (dc *TokenCache) ProcessToken(token string) error {
	dc.mutex.Lock()
	defer dc.mutex.Unlock()

	if _, found := dc.cache.Get(token); found {
		return errors.New("duplicate request: this deduplication token has already been used")
	}

	dc.cache.Set(token, true, cache.DefaultExpiration)
	return nil
}
