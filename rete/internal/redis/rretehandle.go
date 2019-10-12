package redis

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/internal/types"
)

//Holds a tuple reference and related state

type reteHandleImpl struct {
	types.NwElemIdImpl
	tuple    model.Tuple
	tupleKey model.TupleKey
	key      string
	status   types.ReteHandleStatus
	id       int64
	//jtRefs   types.JtRefsService
}

func newReteHandleImpl(nw types.Network, tuple model.Tuple, key string, status types.ReteHandleStatus, id int64) types.ReteHandle {
	h1 := reteHandleImpl{}
	h1.initHandleImpl(nw, tuple, key, status)
	return &h1
}

func (hdl *reteHandleImpl) SetTuple(tuple model.Tuple) {
	hdl.tuple = tuple
	if tuple != nil {
		hdl.tupleKey = tuple.GetKey()
	}
}

func (hdl *reteHandleImpl) initHandleImpl(nw types.Network, tuple model.Tuple, key string, status types.ReteHandleStatus) {
	hdl.SetID(nw)
	hdl.SetTuple(tuple)
	hdl.key = key
	hdl.status = status
}

func (hdl *reteHandleImpl) GetTuple() model.Tuple {
	return hdl.tuple
}

func (hdl *reteHandleImpl) GetTupleKey() model.TupleKey {
	return hdl.tupleKey
}

func (hdl *reteHandleImpl) SetStatus(status types.ReteHandleStatus) {
	if hdl.key == "" {
		return
	}
	redisutils.GetRedisHdl().HSet(hdl.key, "status", status)
}

func (hdl *reteHandleImpl) Unlock() {
	redisutils.GetRedisHdl().HDel(hdl.key, "id")
}

func (hdl *reteHandleImpl) GetStatus() types.ReteHandleStatus {
	return hdl.status
}

func (hdl *reteHandleImpl) AddJoinTableRowRef(joinTableRowVar types.JoinTableRow, joinTableVar types.JoinTable) {
	hdl.Nw.GetJtRefService().AddEntry(hdl, joinTableVar.GetName(), joinTableRowVar.GetID())
}

func (hdl *reteHandleImpl) GetRefTableIterator() types.JointableIterator {
	refTblIterator := hdl.Nw.GetJtRefService().GetRowIterator(nil, hdl)
	return refTblIterator
}
