package rete

import (
	"github.com/project-flogo/rules/rete/common"
	"github.com/project-flogo/rules/rete/internal/mem"
	"github.com/project-flogo/rules/rete/internal/redis"
	"github.com/project-flogo/rules/rete/internal/types"
)

type TypeFactory struct {
	nw     *reteNetworkImpl
	config common.Config
}

func NewFactory(nw *reteNetworkImpl, config common.Config) (*TypeFactory, error) {
	tf := TypeFactory{
		nw:     nw,
		config: config,
	}

	return &tf, nil
}

func (f *TypeFactory) getJoinTableRefs() types.JtRefsService {
	switch f.config.Rete.Jt {
	case common.ServiceTypeMem:
		return mem.NewJoinTableRefsInHdlImpl(f.nw, f.config)
	case common.ServiceTypeRedis:
		return redis.NewJoinTableRefsInHdlImpl(f.nw, f.config)
	default:
		panic("invalid service type")
	}
}

func (f *TypeFactory) getJoinTableCollection() types.JtService {
	switch f.config.Rete.Jt {
	case common.ServiceTypeMem:
		return mem.NewJoinTableCollection(f.nw, f.config)
	case common.ServiceTypeRedis:
		return redis.NewJoinTableCollection(f.nw, f.config)
	default:
		panic("invalid service type")
	}
}

func (f *TypeFactory) getHandleCollection() types.HandleService {
	switch f.config.Rete.JtRef {
	case common.ServiceTypeMem:
		return mem.NewHandleCollection(f.nw, f.config)
	case common.ServiceTypeRedis:
		return redis.NewHandleCollection(f.nw, f.config)
	default:
		panic("invalid service type")
	}
}

func (f *TypeFactory) getIdGen() types.IdGen {
	switch f.config.Rete.IDGenRef {
	case common.ServiceTypeMem:
		return mem.NewIdGenImpl(f.nw, f.config)
	case common.ServiceTypeRedis:
		return redis.NewIdGenImpl(f.nw, f.config)
	default:
		panic("invalid service type")
	}
}

func (f *TypeFactory) getLockService() types.LockService {
	switch f.config.Rete.IDGenRef {
	case common.ServiceTypeMem:
		return mem.NewLockServiceImpl(f.nw, f.config)
	case common.ServiceTypeRedis:
		return redis.NewLockServiceImpl(f.nw, f.config)
	default:
		panic("invalid service type")
	}
}
