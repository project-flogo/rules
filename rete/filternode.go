package rete

import (
	"context"
	"strconv"

	"github.com/project-flogo/rules/common/model"
)

//filter node holds the filter condition
type filterNode interface {
	node
}

type filterNodeImpl struct {
	nodeImpl
	conditionVar model.Condition
	convert      []int
}

func newFilterNode(nw Network, rule model.Rule, identifiers []model.TupleType, conditionVar model.Condition) filterNode {
	fn := filterNodeImpl{}
	fn.initFilterNodeImpl(nw, rule, identifiers, conditionVar)
	return &fn
}

func (fn *filterNodeImpl) initFilterNodeImpl(nw Network, rule model.Rule, identifiers []model.TupleType, conditionVar model.Condition) {
	fn.nodeImpl.initNodeImpl(nw, rule, identifiers)
	fn.conditionVar = conditionVar
	fn.setConvert()
}

func (fn *filterNodeImpl) setConvert() {

	if fn.conditionVar == nil {
		return
	}
	conIdrs := fn.conditionVar.GetIdentifiers()

	if conIdrs != nil && len(conIdrs) == 0 {
		for i, condIdr := range conIdrs {
			idx := GetIndex(fn.identifiers, condIdr)
			if idx != -1 {
				fn.convert[i] = idx
			} else {
				//TODO ERROR HANDLING
			}
		}
	}

}

func (fn *filterNodeImpl) String() string {
	cond := ""
	for _, idr := range fn.conditionVar.GetIdentifiers() {
		cond += string(idr) + " "
	}

	linkTo := ""
	switch fn.nodeLinkVar.getChild().(type) {
	case *joinNodeImpl:
		if fn.nodeLinkVar.isRightNode() {
			linkTo += "j" + strconv.Itoa(fn.nodeLinkVar.getChild().getID()) + "R"
		} else {
			linkTo += "j" + strconv.Itoa(fn.nodeLinkVar.getChild().getID()) + "L"
		}
	case *filterNodeImpl:
		linkTo += "f" + strconv.Itoa(fn.nodeLinkVar.getChild().getID())
	case *ruleNodeImpl:
		linkTo += "r" + strconv.Itoa(fn.nodeLinkVar.getChild().getID())
	}

	return "\t[FilterNode id(" + strconv.Itoa(fn.nodeImpl.getID()) + ") link(" + linkTo + "):\n" +
		"\t\tIdentifier            = " + model.IdentifiersToString(fn.identifiers) + " ;\n" +
		"\t\tCondition Identifiers = " + cond + ";\n" +
		"\t\tCondition             = " + fn.conditionVar.String() + "]"
}

func (fn *filterNodeImpl) assertObjects(ctx context.Context, handles []reteHandle, isRight bool) {
	if fn.conditionVar == nil {
		fn.nodeLinkVar.propagateObjects(ctx, handles)
	} else {
		//TODO: rete listeners...
		var tuples []model.Tuple
		// tupleMap := map[model.TupleType]model.Tuple{}
		if fn.convert == nil {
			tuples = copyIntoTupleArray(handles)
		} else {
			tuples = make([]model.Tuple, len(fn.convert))
			for i := 0; i < len(fn.convert); i++ {
				tuples[i] = handles[fn.convert[i]].getTuple()
				// tupleMap[tuples[i].GetTupleType()] = tuples[i]
			}
		}
		tupleMap := convertToTupleMap(tuples)
		cv := fn.conditionVar
		toPropagate := cv.GetEvaluator()(cv.GetName(), cv.GetRule().GetName(), tupleMap, cv.GetContext())
		if toPropagate {
			fn.nodeLinkVar.propagateObjects(ctx, handles)
		}
	}
}
