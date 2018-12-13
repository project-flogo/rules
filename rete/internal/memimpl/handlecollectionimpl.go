package memimpl

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/internal/types"
)

type handleCollectionImpl struct {
	allHandles map[string]types.ReteHandle
}

func NewHandleCollection() types.HandleCollection {
	hc := handleCollectionImpl{}
	hc.allHandles = make(map[string]types.ReteHandle)
	return &hc
}

func (hc *handleCollectionImpl) AddHandle(hdl types.ReteHandle) {
	hc.allHandles[hdl.GetTupleKey().String()] = hdl
}

func (hc *handleCollectionImpl) RemoveHandle(tuple model.Tuple) types.ReteHandle {
	rh, found := hc.allHandles[tuple.GetKey().String()]
	if found {
		delete(hc.allHandles, tuple.GetKey().String())
		return rh
	}
	return nil
}

func (hc *handleCollectionImpl) GetHandle(tuple model.Tuple) types.ReteHandle {
	return hc.allHandles[tuple.GetKey().String()]
}

func (hc *handleCollectionImpl) GetHandleByKey(key model.TupleKey) types.ReteHandle {
	return hc.allHandles[key.String()]
}
