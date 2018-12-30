package redis

import (
	"fmt"
	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/internal/types"
)

const (
	hget    = "HGET"
	hincrby = "HINCRBY"
)

type idGenServiceImpl struct {
	types.NwServiceImpl

	config map[string]interface{}
	//current int
	rh redisutils.RedisHdl
}

func NewIdGenImpl(config map[string]interface{}) types.IdGen {
	r := idGenServiceImpl{}
	r.config = config
	return &r
}

func (ri *idGenServiceImpl) Init() {
	redisutils.InitService(ri.config)
	ri.rh = redisutils.GetRedisHdl()
	j := ri.GetMaxID()
	fmt.Printf("maxid : [%d]\n ", j)
}

func (ri *idGenServiceImpl) GetMaxID() int {
	return ri.rh.HGetAsInt("IDGEN", "ID")
}

func (ri *idGenServiceImpl) GetNextID() int {
	return ri.rh.HIncrBy("IDGEN", "ID", 1)
}
