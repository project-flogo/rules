package rete

import (
	"encoding/json"

	"github.com/project-flogo/rules/rete/common"
	"github.com/project-flogo/rules/rete/internal/mem"
	"github.com/project-flogo/rules/rete/internal/redis"
	"github.com/project-flogo/rules/rete/internal/types"
)

type TypeFactory struct {
	nw     *reteNetworkImpl
	config string
	parsed common.Config
}

func NewFactory(nw *reteNetworkImpl, config string) (*TypeFactory, error) {
	tf := TypeFactory{}
	tf.config = config
	err := json.Unmarshal([]byte(config), &tf.parsed)
	if err != nil {
		return nil, err
	}
	tf.nw = nw

	return &tf, nil
}

func (f *TypeFactory) getJoinTableRefs() types.JtRefsService {
	switch f.parsed.Rete.Jt {
	case common.ServiceTypeMem:
		return mem.NewJoinTableRefsInHdlImpl(f.nw, f.parsed)
	case common.ServiceTypeRedis:
		return redis.NewJoinTableRefsInHdlImpl(f.nw, f.parsed)
	default:
		panic("invalid service type")
	}
}

func (f *TypeFactory) getJoinTableCollection() types.JtService {
	switch f.parsed.Rete.Jt {
	case common.ServiceTypeMem:
		return mem.NewJoinTableCollection(f.nw, f.parsed)
	case common.ServiceTypeRedis:
		return redis.NewJoinTableCollection(f.nw, f.parsed)
	default:
		panic("invalid service type")
	}
}

func (f *TypeFactory) getHandleCollection() types.HandleService {
	switch f.parsed.Rete.JtRef {
	case common.ServiceTypeMem:
		return mem.NewHandleCollection(f.nw, f.parsed)
	case common.ServiceTypeRedis:
		return redis.NewHandleCollection(f.nw, f.parsed)
	default:
		panic("invalid service type")
	}
}

func (f *TypeFactory) getIdGen() types.IdGen {
	switch f.parsed.Rete.IDGenRef {
	case common.ServiceTypeMem:
		return mem.NewIdGenImpl(f.nw, f.parsed)
	case common.ServiceTypeRedis:
		return redis.NewIdGenImpl(f.nw, f.parsed)
	default:
		panic("invalid service type")
	}
}
