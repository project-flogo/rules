package redis

import (
	"context"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/common"
	"github.com/project-flogo/rules/rete/internal/types"
)

type handleServiceImpl struct {
	types.NwServiceImpl
	prefix string
	config common.Config
	rand.Source
	sync.Mutex
	redisutils.RedisHdl
}

func NewHandleCollection(nw types.Network, config common.Config) types.HandleService {
	hc := handleServiceImpl{
		NwServiceImpl: types.NwServiceImpl{
			Nw: nw,
		},
		config:   config,
		Source:   rand.NewSource(time.Now().UnixNano()),
		RedisHdl: redisutils.NewRedisHdl(config.Jts.Redis),
	}
	return &hc
}

func (hc *handleServiceImpl) Int63() int64 {
	hc.Lock()
	defer hc.Unlock()
	return hc.Source.Int63()
}

func (hc *handleServiceImpl) Init() {
	hc.prefix = hc.Nw.GetPrefix() + ":h:"
}

func (hc *handleServiceImpl) RemoveHandle(tuple model.Tuple) types.ReteHandle {
	rkey := hc.prefix + tuple.GetKey().String()
	hc.Del(rkey)
	//TODO: Dummy handle
	h := newReteHandleImpl(hc.GetNw(), hc.RedisHdl, tuple, rkey, types.ReteHandleStatusUnknown, -1)
	return h

}

func (hc *handleServiceImpl) GetHandle(ctx context.Context, tuple model.Tuple) types.ReteHandle {
	return hc.GetHandleByKey(ctx, tuple.GetKey())
}

func (hc *handleServiceImpl) GetHandleByKey(ctx context.Context, key model.TupleKey) types.ReteHandle {
	rkey := hc.prefix + key.String()

	m := hc.HGetAll(rkey)
	if len(m) == 0 {
		return nil
	}
	status := types.ReteHandleStatusUnknown
	if value, ok := m["status"]; ok {
		if value, ok := value.(string); ok {
			number, err := strconv.Atoi(value)
			if err != nil {
				panic(err)
			}
			status = types.ReteHandleStatus(number)
		} else {
			panic("status not string")
		}
	} else {
		panic("missing status")
	}
	id := int64(-1)
	if value, ok := m["id"]; ok {
		if value, ok := value.(string); ok {
			number, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				panic(err)
			}
			id = number
		}
	}

	var tuple model.Tuple
	if ctx != nil {
		if value := ctx.Value(model.RetecontextKeyType{}); value != nil {
			if value, ok := value.(types.ReteCtx); ok {
				if modified := value.GetRtcModified(); modified != nil {
					if value := modified[key.String()]; value != nil {
						tuple = value.GetTuple()
					}
				}
				if tuple == nil {
					if added := value.GetRtcAdded(); added != nil {
						tuple = added[key.String()]
					}
				}
			}
		}
	}
	if tuple == nil {
		tuple = hc.Nw.GetTupleStore().GetTupleByKey(key)
	}
	if tuple == nil {
		//TODO: error handling
		return nil
	}

	h := newReteHandleImpl(hc.GetNw(), hc.RedisHdl, tuple, rkey, status, id)
	return h
}

func (hc *handleServiceImpl) GetOrCreateLockedHandle(nw types.Network, tuple model.Tuple) (types.ReteHandle, bool) {
	id := hc.Int63()
	key, status := hc.prefix+tuple.GetKey().String(), types.ReteHandleStatusCreating

	exists, _ := hc.HSetNX(key, "id", id)
	if exists {
		return nil, true
	}

	exists, _ = hc.HSetNX(key, "status", status)
	if exists {
		m := hc.HGetAll(key)
		if len(m) > 0 {
			if value, ok := m["status"]; ok {
				if value, ok := value.(string); ok {
					number, err := strconv.Atoi(value)
					if err != nil {
						panic(err)
					}
					status = types.ReteHandleStatus(number)
				} else {
					panic("status not string")
				}
			} else {
				panic("missing status")
			}
		}
	}

	h := newReteHandleImpl(nw, hc.RedisHdl, tuple, key, status, id)
	return h, false
}

func (hc *handleServiceImpl) GetLockedHandle(nw types.Network, tuple model.Tuple) (types.ReteHandle, bool) {
	id := hc.Int63()
	key := hc.prefix + tuple.GetKey().String()

	exists, _ := hc.HSetNX(key, "id", id)
	if exists {
		return nil, true
	}

	m := hc.HGetAll(key)
	if len(m) > 0 {
		if value, ok := m["status"]; ok {
			if value, ok := value.(string); ok {
				number, err := strconv.Atoi(value)
				if err != nil {
					panic(err)
				}
				h := newReteHandleImpl(nw, hc.RedisHdl, tuple, key, types.ReteHandleStatus(number), id)
				return h, false
			}
		}
	}

	return nil, true
}

func (hc *handleServiceImpl) GetHandleWithTuple(nw types.Network, tuple model.Tuple) types.ReteHandle {
	key, status, id := hc.prefix+tuple.GetKey().String(), types.ReteHandleStatusCreating, int64(-1)
	m := hc.HGetAll(key)
	if len(m) > 0 {
		if value, ok := m["status"]; ok {
			if value, ok := value.(string); ok {
				number, err := strconv.Atoi(value)
				if err != nil {
					panic(err)
				}
				status = types.ReteHandleStatus(number)
			} else {
				panic("status not string")
			}
		} else {
			panic("missing status")
		}
		if value, ok := m["id"]; ok {
			if value, ok := value.(string); ok {
				number, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					panic(err)
				}
				id = number
			}
		}
	}

	h := newReteHandleImpl(nw, hc.RedisHdl, tuple, key, status, id)
	return h
}
