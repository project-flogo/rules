package rete

import (
	"container/list"
	"context"

	"github.com/project-flogo/rules/common/model"
)

type conflictRes interface {
	addAgendaItem(rule model.Rule, tupleMap map[model.TupleType]model.Tuple)
	resolveConflict(ctx context.Context)
	deleteAgendaFor(ctx context.Context, tuple model.Tuple, changeProps map[string]bool)
}

type conflictResImpl struct {
	agendaList list.List
}

func newConflictRes() conflictRes {
	cr := conflictResImpl{}
	cr.initCR()
	return &cr
}

func (cr *conflictResImpl) initCR() {
	cr.agendaList = list.List{}
}

func (cr *conflictResImpl) addAgendaItem(rule model.Rule, tupleMap map[model.TupleType]model.Tuple) {
	item := newAgendaItem(rule, tupleMap)
	v := rule.GetPriority()
	found := false
	for e := cr.agendaList.Front(); e != nil; e = e.Next() {
		curr := e.Value.(agendaItem)
		if v <= curr.getRule().GetPriority() {
			cr.agendaList.InsertBefore(item, e)
			found = true
			break
		}
	}
	if !found {
		cr.agendaList.PushBack(item)
	}
}

func (cr *conflictResImpl) resolveConflict(ctx context.Context) {
	var item agendaItem

	front := cr.agendaList.Front()
	for front != nil {
		val := cr.agendaList.Remove(front)
		if val != nil {
			item = val.(agendaItem)
			actionTuples := item.getTuples()
			actionFn := item.getRule().GetActionFn()
			if actionFn != nil {
				reteCtxV := getReteCtx(ctx)
				actionFn(ctx, reteCtxV.getRuleSession(), item.getRule().GetName(), actionTuples, item.getRule().GetContext())
			}
		}

		reteCtxV := getReteCtx(ctx)

		reteCtxV.addRuleModifiedToOpsList()

		reteCtxV.copyRuleModifiedToRtcModified()
		//action scoped, clear it for the next action
		reteCtxV.resetModified()

		if reteCtxV != nil {
			opsFront := reteCtxV.getOpsList().Front()
			for opsFront != nil {
				opsVal := reteCtxV.getOpsList().Remove(opsFront)
				oprn := opsVal.(opsEntry)
				oprn.execute(ctx)
				opsFront = reteCtxV.getOpsList().Front()
			}
		}

		front = cr.agendaList.Front()
	}

	reteCtxV := getReteCtx(ctx)
	reteCtxV.normalize()
	//reteCtxV.printRtcChangeList()

}

func (cr *conflictResImpl) deleteAgendaFor(ctx context.Context, modifiedTuple model.Tuple, changeProps map[string]bool) {

	hdlModified := getOrCreateHandle(ctx, modifiedTuple)

	for e := cr.agendaList.Front(); e != nil; {
		item := e.Value.(agendaItem)
		next := e.Next()
		for _, tuple := range item.getTuples() {
			hdl := getOrCreateHandle(ctx, tuple)
			if hdl == hdlModified { //this agendaitem has the modified tuple, remove the agenda item!
				toRemove := true
				//check if the rule depends on this change prop
				if changeProps != nil {
					if depProps, found := item.getRule().GetDeps()[tuple.GetTupleType()]; found {
						if len(depProps) > 0 {
							for prop := range depProps {
								if _, fnd := changeProps[prop]; fnd {
									toRemove = true
									break
								}
							}
						} else {
							toRemove = false
						}
					} else {
						toRemove = false
					}
				}
				if toRemove {
					cr.agendaList.Remove(e)
					break
				}
			}
		}
		e = next
	}

}
