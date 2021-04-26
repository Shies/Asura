package dao

import (
	log "Asura/src/logger"

	"github.com/bradfitz/gomemcache/memcache"
)

const (
	ExpireTime = 3*24*3600
	_healthCheck = "hello world"
)

func (d *Dao) MemcacheGet(key string) (item *memcache.Item, err error) {
	item, err = d.mc.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			item = nil
			err = nil
		} else {
			log.Error("memcache get error(%v)", err)
			return
		}
	}
	if item == nil {
		item = &memcache.Item{Key: key, Value: []byte(""), Expiration: ExpireTime}
	}
	return
}

func (d *Dao) MemcacheSet(key, val string) (err error) {
	err = d.mc.Set(&memcache.Item{Key: key, Value: []byte(val), Expiration: ExpireTime})
	if err != nil {
		log.Error("memcache set error(%v)", err)
		return
	}

	return
}