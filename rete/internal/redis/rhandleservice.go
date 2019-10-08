package redis

import (
	"context"
	"strconv"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/internal/types"
)

type handleServiceImpl struct {
	//allHandles map[string]types.ReteHandle
	types.NwServiceImpl
	prefix string
	config map[string]interface{}
}

func NewHandleCollection(nw types.Network, config map[string]interface{}) types.HandleService {
	hc := handleServiceImpl{}
	hc.Nw = nw
	hc.config = config
	//hc.allHandles = make(map[string]types.ReteHandle)
	return &hc
}

func (hc *handleServiceImpl) Init() {
	hc.prefix = hc.Nw.GetPrefix() + ":h:"
	reteCfg := hc.config["rete"].(map[string]interface{})
	jtRef := reteCfg["jt-ref"].(string)
	jts := hc.config["jts"].(map[string]interface{})
	redisCfg := jts[jtRef].(map[string]interface{})
	redisutils.InitService(redisCfg)
}

func (hc *handleServiceImpl) RemoveHandle(tuple model.Tuple) types.ReteHandle {
	rkey := hc.prefix + tuple.GetKey().String()
	redisutils.GetRedisHdl().Del(rkey)
	//TODO: Dummy handle
	h := newReteHandleImpl(hc.GetNw(), tuple, rkey, types.ReteHandleStatusUnknown)
	return h

}

func (hc *handleServiceImpl) GetHandle(ctx context.Context, tuple model.Tuple) types.ReteHandle {
	return hc.GetHandleByKey(ctx, tuple.GetKey())
}

func (hc *handleServiceImpl) GetHandleByKey(ctx context.Context, key model.TupleKey) types.ReteHandle {
	rkey := hc.prefix + key.String()

	m := redisutils.GetRedisHdl().HGetAll(rkey)
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

	h := newReteHandleImpl(hc.GetNw(), tuple, rkey, status)
	return h
}

func (hc *handleServiceImpl) GetOrCreateHandle(nw types.Network, tuple model.Tuple) (types.ReteHandle, bool) {
	key, status := hc.prefix+tuple.GetKey().String(), types.ReteHandleStatusCreating
	exists, _ := redisutils.GetRedisHdl().HSetNX(key, "status", status)
	if exists {
		m := redisutils.GetRedisHdl().HGetAll(key)
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

	h := newReteHandleImpl(nw, tuple, key, status)
	return h, exists
}
