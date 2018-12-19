package redisutils

import (
	"github.com/gomodule/redigo/redis"
)

var rd RedisHdl

type RedisHdl = *RedisHandle

type RedisHandle struct {
	config string
	pool   *redis.Pool
}

func InitService(config string) {
	if rd == nil {
		rd = &RedisHandle{}
		rd.config = config
		rd.newPool("tcp", ":6379")
	}
}

func GetRedisHdl() RedisHdl {
	return rd
}

func (rh *RedisHandle) newPool(network string, address string) {
	rh.pool = &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(network, address)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func (rh *RedisHandle) GetPool() *redis.Pool {
	return rh.pool
}
