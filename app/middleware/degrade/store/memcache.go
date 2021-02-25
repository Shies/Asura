package store

import (
	"context"

	"Asura/conf"
	"Asura/src/logger"

	"github.com/bradfitz/gomemcache/memcache"
)

// Memcache represents the cache with memcached persistence
type Memcache struct {
	conf   *conf.Memcache
	client *memcache.Client
}

// NewMemcache new a memcache store.
func NewMemcache(c *conf.Memcache) *Memcache {
	if c == nil {
		panic("cache config is nil")
	}

	// pick memcache cluster the frist node
	return &Memcache{
		client: memcache.New(c.Host+":"+c.Port),
		conf: c,
	}
}

// Set save the result to memcache store.
func (ms *Memcache) Set(ctx context.Context, key string, value []byte, expire int32) (err error) {
	err = ms.client.Set(&memcache.Item{Key: key, Value: value, Expiration: expire})
	if err != nil {
		logger.Error("memcache set error(%v)", err)
		return
	}

	return
}

// Get get result from mc by locaiton+params.
func (ms *Memcache) Get(ctx context.Context, key string) ([]byte, error) {
	item, err := ms.client.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			item = nil
			err = nil
		} else {
			logger.Error("memcache get error(%v)", err)
			return []byte(""), err
		}
	}
	if item == nil {
		item = &memcache.Item{Key: key, Value: []byte(""), Expiration: 36400}
	}

	return item.Value, nil
}
