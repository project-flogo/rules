package rete

import (
	"context"

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

func newReteHandleImpl(nw *reteNetworkImpl, tuple model.Tuple) types.ReteHandle {
	h1 := reteHandleImpl{}
	h1.initHandleImpl(nw, tuple)
	return &h1
}

func (hdl *reteHandleImpl) SetTuple(tuple model.Tuple) {
	hdl.tuple = tuple
}

func (hdl *reteHandleImpl) initHandleImpl(nw *reteNetworkImpl, tuple model.Tuple) {
	hdl.SetID(nw)
	hdl.SetTuple(tuple)
	hdl.jtRefs = nw.getFactory().getJoinTableRefs()
}

func (hdl *reteHandleImpl) GetTuple() model.Tuple {
	return hdl.tuple
}

func (hdl *reteHandleImpl) GetTupleKey() model.TupleKey {
	return hdl.tupleKey
}

func getOrCreateHandle(ctx context.Context, tuple model.Tuple) types.ReteHandle {
	reteCtxVar := getReteCtx(ctx)
	return reteCtxVar.getNetwork().getOrCreateHandle(ctx, tuple)
}

func (hdl *reteHandleImpl) AddJoinTableRowRef(joinTableRowVar types.JoinTableRow, joinTableVar types.JoinTable) {
	hdl.jtRefs.AddEntry(joinTableVar.GetID(), joinTableRowVar.GetID())
}

func (hdl *reteHandleImpl) RemoveJoinTableRowRefs(changedProps map[string]bool) {

	tuple := hdl.tuple
	alias := tuple.GetTupleType()

	hdlTblIter := hdl.GetRefTableIterator()

	for hdlTblIter.HasNext() {
		joinTableID, rowIDs := hdlTblIter.Next()
		joinTable := hdl.Nw.GetJoinTable(joinTableID)
		toDelete := false
		if changedProps != nil {
			rule := joinTable.GetRule()
			depProps, found := rule.GetDeps()[alias]
			if found { // rule depends on this type
				for changedProp := range changedProps {
					_, foundProp := depProps[changedProp]
					if foundProp {
						toDelete = true
						break
					}
				}
			}
		} else {
			toDelete = true
		}

		if !toDelete {
			continue
		}
		//this can happen if some other handle removed a row as a result of retraction
		if rowIDs == nil {
			continue
		}
		////Remove rows from corresponding join tables
		for e := rowIDs.Front(); e != nil; e = e.Next() {
			rowID := e.Value.(int)
			row := joinTable.RemoveRow(rowID)
			if row != nil {
				//Remove other refs recursively.
				for _, otherHdl := range row.GetHandles() {
					//if otherHdl != nil {
					otherHdl.RemoveJoinTableRowRefs(nil)
					//}
				}
			}
		}

		//Remove the reference to the table itself
		hdl.jtRefs.RemoveEntry(joinTableID)
	}
}

//Used when a rule is deleted. See Network.RemoveRule
func (hdl *reteHandleImpl) RemoveJoinTable(joinTableID int) {
	hdl.jtRefs.RemoveEntry(joinTableID)
}

func (hdl *reteHandleImpl) GetRefTableIterator() types.HdlTblIterator {
	refTblIteator := hdl.jtRefs.GetIterator()
	return refTblIteator

}
