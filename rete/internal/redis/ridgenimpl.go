package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/internal/types"
)

const (
	hget    = "HGET"
	hincrby = "HINCRBY"
)

type ridImpl struct {
	config string
	//current int
	rh redisutils.RedisHdl
}

func NewIdImpl(config string) types.IdGen {
	r := ridImpl{}
	r.config = config
	return &r
}

func (ri *ridImpl) Init() {
	redisutils.InitService(ri.config)
	j := ri.GetMaxID()
	fmt.Printf("maxid : [%d]\n ", j)
}

func (ri *ridImpl) GetMaxID() int {
	ri.rh = redisutils.GetRedisHdl()
	c := ri.rh.GetPool().Get()
	defer c.Close()

	i, err := c.Do(hget, "IDGEN", "ID")
	if err == nil {
		j, _ := redis.Int(i, err)
		return j
	}
	return -1
}

func (ri *ridImpl) GetNextID() int {
	ri.rh = redisutils.GetRedisHdl()
	c := ri.rh.GetPool().Get()
	defer c.Close()

	i, err := c.Do(hincrby, "IDGEN", "ID", 1)

	if err != nil {
		fmt.Printf("error: [%s]", err)
		return -1
	}
	current := int(i.(int64))
	return current

}
