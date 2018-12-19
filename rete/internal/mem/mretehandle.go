package mem

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/internal/types"
)

//Holds a tuple reference and related state

type reteHandleImpl struct {
	types.NwElemIdImpl
	tuple    model.Tuple
	tupleKey model.TupleKey
	jtRefs   types.JoinTableRefsInHdl
}

func newReteHandleImpl(nw types.Network, tuple model.Tuple) types.ReteHandle {
	h1 := reteHandleImpl{}
	h1.initHandleImpl(nw, tuple)
	return &h1
}

func (hdl *reteHandleImpl) SetTuple(tuple model.Tuple) {
	hdl.tuple = tuple
	hdl.tupleKey = tuple.GetKey()
}

func (hdl *reteHandleImpl) initHandleImpl(nw types.Network, tuple model.Tuple) {
	hdl.SetID(nw)
	hdl.SetTuple(tuple)
	hdl.tupleKey = tuple.GetKey()
	hdl.jtRefs = NewJoinTableRefsInHdlImpl()
}

func (hdl *reteHandleImpl) GetTuple() model.Tuple {
	return hdl.tuple
}

func (hdl *reteHandleImpl) GetTupleKey() model.TupleKey {
	return hdl.tupleKey
}

func (hdl *reteHandleImpl) AddJoinTableRowRef(joinTableRowVar types.JoinTableRow, joinTableVar types.JoinTable) {
	hdl.jtRefs.AddEntry(joinTableVar.GetID(), joinTableRowVar.GetID())
}

//Used when a rule is deleted. See Network.RemoveRule
func (hdl *reteHandleImpl) RemoveJoinTable(joinTableID int) {
	hdl.jtRefs.RemoveEntry(joinTableID)
}

func (hdl *reteHandleImpl) GetRefTableIterator() types.HdlTblIterator {
	refTblIteator := hdl.jtRefs.GetIterator()
	return refTblIteator
}
