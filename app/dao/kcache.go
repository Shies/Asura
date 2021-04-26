package dao

import (
	"fmt"

	log "Asura/src/logger"

	"github.com/garyburd/redigo/redis"
)

const (
	_cacheKey = "cacheKey_%s"
)

func (d *Dao) CacheKey(key string) string {
	return fmt.Sprintf(_cacheKey, key)
}

func (d *Dao) PipeLine(values map[string]string) (reply interface{}, err error) {
	var (
		conn = d.kredis.Get()
	)
	defer conn.Close()
	for key, val := range values {
		cacheKey := d.CacheKey(key)
		_ = conn.Send("SET", cacheKey, val)
	}

	_ = conn.Flush()
	if reply, err = conn.Receive(); err != nil {
		log.Error("Note: the redis pipeline recv fail(%v)", err)
		return
	}
	return
}

func (d *Dao) GetValues(key string) (reply string, err error) {
	var (
		conn = d.kredis.Get()
		_key = d.CacheKey(key)
	)

	defer conn.Close()
	isExists, _ := redis.Bool(conn.Do("EXISTS", _key))
	if (isExists) {
		reply, err = redis.String(conn.Do("GET", _key))
		if err != nil {
			log.Error("redis getValues error(%v)", err)
			return
		}
	}
	return
}