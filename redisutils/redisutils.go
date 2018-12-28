package redisutils

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

var rd RedisHdl

type RedisHdl = *RedisHandle

type RedisHandle struct {
	config map[string]interface{}
	pool   *redis.Pool
}

func InitService(config map[string]interface{}) {
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

func (rh *RedisHandle) getPool() *redis.Pool {
	return rh.pool
}

func (rh *RedisHandle) HSetAll(key string, kvs map[string]interface{}) error {

	var args = []interface{}{key}
	for f, v := range kvs {
		args = append(args, f, v)
	}
	c := rd.getPool().Get()
	defer c.Close()
	_, error := c.Do("HMSET", args...)

	return error
}

func (rh *RedisHandle) HGetAll(key string) map[string]string {
	hgetall := make(map[string]string)
	c := rh.getPool()
	defer c.Close()
	vals, error := c.Get().Do("HGETALL", key)
	if error != nil {
		fmt.Printf("error [%v]\n", error)
	} else {
		vals, err2 := redis.Values(vals, error)
		if err2 != nil {
			fmt.Printf("error [%v]\n", err2)
		} else {
			i := 0
			key := ""
			value := ""
			for _, val := range vals {
				ba := val.([]byte)
				s := string(ba)
				//fmt.Printf("Value [%s]\n", s)
				if i%2 == 0 {
					key = s
				} else {
					value = s
					hgetall[key] = value
				}
				i++
			}
		}
	}
	return hgetall
}

func (rh *RedisHandle) HIncrBy(key string, field string, by int) int {
	c := rh.getPool().Get()
	defer c.Close()
	i, err := c.Do("HINCRBY", key, field, 1)

	if err != nil {
		fmt.Printf("error: [%s]", err)
		return -1
	}
	current := int(i.(int64))
	return current
}

func (rh *RedisHandle) HGetAsInt(key string, field string) int {
	c := rh.getPool().Get()
	defer c.Close()
	i, err := c.Do("HGET", key, field)
	j := -1
	if err == nil {
		j, _ = redis.Int(i, err)
	}
	return j
}
