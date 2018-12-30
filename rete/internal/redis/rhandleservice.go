package redis

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/internal/types"
	"github.com/project-flogo/rules/redisutils"
	"fmt"
)

type handleServiceImpl struct {
	//allHandles map[string]types.ReteHandle
	types.NwServiceImpl
}

func NewHandleCollection(config map[string]interface{}) types.HandleService {
	hc := handleServiceImpl{}
	//hc.allHandles = make(map[string]types.ReteHandle)
	return &hc
}

func (hc *handleServiceImpl) Init() {

}

//func (hc *handleServiceImpl) AddHandle(hdl types.ReteHandle) {
//	//hc.allHandles[hdl.GetTupleKey().String()] = hdl
//
//	redisutils.GetRedisHdl().HSetAll(key, m)
//}

func (hc *handleServiceImpl) RemoveHandle(tuple model.Tuple) types.ReteHandle {

	numDeleted := redisutils.GetRedisHdl().Del("h-"+tuple.GetKey().String())
	fmt.Printf("Deleted: [%d] keys\n", numDeleted)

	//TODO: Dummy handle
	h := newReteHandleImpl(hc.GetNw(), tuple)
	return h

}

func (hc *handleServiceImpl) GetHandle(tuple model.Tuple) types.ReteHandle {
	return hc.GetHandleByKey(tuple.GetKey())
}

func (hc *handleServiceImpl) GetHandleByKey(key model.TupleKey) types.ReteHandle {
	rkey := "h-" + key.String()

	m := redisutils.GetRedisHdl().HGetAll(rkey)
	if len (m) == 0 {
		return nil
	} else {
		tuple := hc.Nw.GetTupleStore().GetTupleByStringKey(key.String())
		if tuple == nil {
			//TODO: error handling
			return nil
		}
		h := newReteHandleImpl(hc.GetNw(), tuple)
		return h
	}
}

func (hc *handleServiceImpl) GetOrCreateHandle(nw types.Network, tuple model.Tuple) types.ReteHandle {

	key := "h-" + tuple.GetKey().String()

	m := redisutils.GetRedisHdl().HGetAll(key)
	if len (m) == 0 {
		m := make(map[string]interface{})
		m["k"] = "v"
		redisutils.GetRedisHdl().HSetAll(key, m)
	}

	h := newReteHandleImpl(nw, tuple)
	//hc.allHandles[tuple.GetKey().String()] = h

	return h
}
