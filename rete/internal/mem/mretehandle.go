package mem

import (
	"sync/atomic"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/internal/types"
)

//Holds a tuple reference and related state

type reteHandleImpl struct {
	types.NwElemIdImpl
	tuple    model.Tuple
	tupleKey model.TupleKey
	status   types.ReteHandleStatus
	id       int64
}

func newReteHandleImpl(nw types.Network, tuple model.Tuple, status types.ReteHandleStatus, id int64) *reteHandleImpl {
	h1 := reteHandleImpl{}
	h1.initHandleImpl(nw, tuple, status, id)
	return &h1
}

func (hdl *reteHandleImpl) SetTuple(tuple model.Tuple) {
	hdl.tuple = tuple
	hdl.tupleKey = tuple.GetKey()
}

func (hdl *reteHandleImpl) initHandleImpl(nw types.Network, tuple model.Tuple, status types.ReteHandleStatus, id int64) {
	hdl.SetID(nw)
	hdl.SetTuple(tuple)
	hdl.tupleKey = tuple.GetKey()
	hdl.status = status
	hdl.id = id
}

func (hdl *reteHandleImpl) GetTuple() model.Tuple {
	return hdl.tuple
}

func (hdl *reteHandleImpl) GetTupleKey() model.TupleKey {
	return hdl.tupleKey
}

func (hdl *reteHandleImpl) SetStatus(status types.ReteHandleStatus) {
	hdl.status = status
}

func (hdl *reteHandleImpl) Unlock() {
	atomic.StoreInt64(&hdl.id, -1)
}

func (hdl *reteHandleImpl) GetStatus() types.ReteHandleStatus {
	return hdl.status
}

func (hdl *reteHandleImpl) AddJoinTableRowRef(joinTableRowVar types.JoinTableRow, joinTableVar types.JoinTable) {
	hdl.Nw.GetJtRefService().AddEntry(hdl, joinTableVar.GetName(), joinTableRowVar.GetID())
}
