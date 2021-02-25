package dao

import (
	"Asura/src/logger"

	"github.com/garyburd/redigo/redis"
)

const (
	_cacheDegrade = "cache_degrade"
	ExpireTime = 3*24*3600
)

func (d *Dao) DegradeCacheKey() string {
	return _cacheDegrade
}

func (d *Dao) SetDegradeCache(ts int64) (reply bool) {
	var (
		conn = d.redis.Get()
		key = d.DegradeCacheKey()
		err error
	)

	defer conn.Close()
	_, err = conn.Do("SET", key, ts)
	if err != nil {
		logger.Error("redis SetDegradeCache error(%v)", err)
		return
	}

	reply = true
	return
}

func (d *Dao) GetDegradeCache() (result int64) {
	var (
		conn = d.redis.Get()
		key = d.DegradeCacheKey()
		err error
	)

	defer conn.Close()
	isExists, _ := redis.Bool(conn.Do("EXISTS", key))
	if (isExists) {
		result, err = redis.Int64(conn.Do("GET", key))
		if err != nil {
			logger.Error("redis GetDegradeCache error(%v)", err)
			return
		}
	}
	return
}