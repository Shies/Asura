package dao

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"Asura/conf"
	_ "Asura/src/logger"
	rpc "Asura/app/rpc/client"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/garyburd/redigo/redis"
)

// noinspection ALL
type Dao struct {
	c 	     *conf.Config
	db       *sql.DB
	redis 	 *redis.Pool
	rslave   *redis.Pool
	kredis   *redis.Pool
	mc 		 *memcache.Client
	grpc	 *rpc.Conn
	err 	 error
}

func New(c *conf.Config) *Dao {
	// cache.NewCache(c)
	dao := &Dao{
		c: c,
		// db: sql.NewMySQL(c.Mysql),
		// mc: cache.NewMemClient(c.Memcache),
		// redis: cache.NewRedisClient(c.Redis.Node1.Master),
		// rslave: cache.NewRedisClient(c.Redis.Node2.Master),
		// kredis: cache.NewKafkaRedis(c.KafkaConsumer.Test.Redis),
		grpc: rpc.NewClient(c.RpcClient),
		// http:   net.NewClient(c.HttpClient),
	}
	if dao.err != nil {
		panic(dao.err)
	}

	return dao
}

func Now(t time.Time) int64 {
	return t.Unix()
}

// JSON2map json to map.
func (d *Dao) JSON2map(msg json.RawMessage) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(msg), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (d *Dao) Ping(c context.Context) (err error) {
	return d.grpc.Call("welcome", rpc.Welcome)
}

func (d *Dao) Close() {
	d.db.Close()
	return
}