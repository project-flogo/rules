package redis

import (
	"fmt"
	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/common"
	"github.com/project-flogo/rules/rete/internal/types"
)

type idGenServiceImpl struct {
	types.NwServiceImpl
	config common.Config
	key    string // key used to access idgen
	fld    string // redis field in key
	redisutils.RedisHdl
}

func NewIdGenImpl(nw types.Network, config common.Config) types.IdGen {
	r := idGenServiceImpl{
		NwServiceImpl: types.NwServiceImpl{
			Nw: nw,
		},
		config:   config,
		RedisHdl: redisutils.NewRedisHdl(config.IDGens.Redis),
	}
	return &r
}

func (ri *idGenServiceImpl) Init() {
	ri.key = ri.Nw.GetPrefix() + ":idgen"
	j := ri.GetMaxID()
	fmt.Printf("maxid : [%d]\n ", j)
}

func (ri *idGenServiceImpl) GetMaxID() int {
	return ri.GetAsInt(ri.key)
}

func (ri *idGenServiceImpl) GetNextID() int {
	return ri.IncrBy(ri.key, 1)
}
