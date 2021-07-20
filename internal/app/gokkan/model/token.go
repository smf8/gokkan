package model

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// TokenBlacklistRepo is the store for blacklisted tokens (i.e. signed out users).
type TokenBlacklistRepo interface {
	Save(jti string) error
	Check(jti string) bool
}

// CacheTokenBlacklistRepo is the in memory implementation of TokenBlacklistRepo.
type CacheTokenBlacklistRepo struct {
	cache *cache.Cache
}

// NewCacheTokenBlacklistRepo creates a new cache repo.
func NewCacheTokenBlacklistRepo(expiration time.Duration) *CacheTokenBlacklistRepo {
	return &CacheTokenBlacklistRepo{
		cache: cache.New(expiration, expiration),
	}
}

// Save saves the token as blocked.
func (c *CacheTokenBlacklistRepo) Save(jti string) error {
	return c.cache.Add(jti, nil, cache.DefaultExpiration)
}

// Check checks blacklist repo to see if the token is blocked.
func (c *CacheTokenBlacklistRepo) Check(jti string) bool {
	_, found := c.cache.Get(jti)

	return found
}
