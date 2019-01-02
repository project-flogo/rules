package redis

import (
	"fmt"
	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/internal/types"
)

type idGenServiceImpl struct {
	types.NwServiceImpl

	config map[string]interface{}
	//current int
	rh redisutils.RedisHdl

	//key used to access idgen
	key string

	//redis field in key
	fld string
}

func NewIdGenImpl(nw types.Network, config map[string]interface{}) types.IdGen {
	r := idGenServiceImpl{}
	r.Nw = nw
	r.config = config
	return &r
}

func (ri *idGenServiceImpl) Init() {
	ri.key = ri.Nw.GetPrefix() + ":idgen"
	redisutils.InitService(ri.config)
	ri.rh = redisutils.GetRedisHdl()
	j := ri.GetMaxID()
	fmt.Printf("maxid : [%d]\n ", j)
}

func (ri *idGenServiceImpl) GetMaxID() int {
	return ri.rh.GetAsInt(ri.key)
}

func (ri *idGenServiceImpl) GetNextID() int {
	return ri.rh.IncrBy(ri.key, 1)
}
