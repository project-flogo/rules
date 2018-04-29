package rete

import (
	"strconv"

	"github.com/TIBCOSoftware/bego/common/model"
)

//filter node holds the filter condition
type filterNode interface {
	node
}

type filterNodeImpl struct {
	nodeImpl
	conditionVar condition
	convert      []int
}

//NewFilterNode ... C'tor
func newFilterNode(identifiers []Identifier, conditionVar condition) filterNode {
	filterNodeImplVar := filterNodeImpl{}
	filterNodeImplVar.initFilterNodeImpl(identifiers, conditionVar)
	return &filterNodeImplVar
}

func (filterNodeImplVar *filterNodeImpl) initFilterNodeImpl(identifiers []Identifier, conditionVar condition) {
	filterNodeImplVar.nodeImpl.initNodeImpl(identifiers)
	filterNodeImplVar.conditionVar = conditionVar
	filterNodeImplVar.setConvert()
}

func (filterNodeImplVar *filterNodeImpl) setConvert() {

	if filterNodeImplVar.conditionVar == nil {
		return
	}
	conIdrs := filterNodeImplVar.conditionVar.getIdentifiers()

	if conIdrs != nil && len(conIdrs) == 0 {
		for i, condIdr := range conIdrs {
			idx := GetIndex(filterNodeImplVar.identifiers, condIdr)
			if idx != -1 {
				filterNodeImplVar.convert[i] = idx
			} else {
				//TODO ERROR HANDLING
			}
		}
	}

}

func (filterNodeImplVar *filterNodeImpl) String() string {
	cond := ""
	for _, idr := range filterNodeImplVar.conditionVar.getIdentifiers() {
		cond += idr.String() + " "
	}

	linkTo := ""
	switch filterNodeImplVar.nodeLinkVar.getChild().(type) {
	case *joinNodeImpl:
		if filterNodeImplVar.nodeLinkVar.isRightNode() {
			linkTo += "j" + strconv.Itoa(filterNodeImplVar.nodeLinkVar.getChild().getID()) + "R"
		} else {
			linkTo += "j" + strconv.Itoa(filterNodeImplVar.nodeLinkVar.getChild().getID()) + "L"
		}
	case *filterNodeImpl:
		linkTo += "f" + strconv.Itoa(filterNodeImplVar.nodeLinkVar.getChild().getID())
	case *ruleNodeImpl:
		linkTo += "r" + strconv.Itoa(filterNodeImplVar.nodeLinkVar.getChild().getID())
	}

	return "\t[FilterNode id(" + strconv.Itoa(filterNodeImplVar.nodeImpl.id) + ") link(" + linkTo + "):\n" +
		"\t\tIdentifier            = " + IdentifiersToString(filterNodeImplVar.identifiers) + " ;\n" +
		"\t\tCondition Identifiers = " + cond + ";\n" +
		"\t\tCondition             = " + filterNodeImplVar.conditionVar.String() + "]"
}

func (filterNodeImplVar *filterNodeImpl) assertObjects(handles []reteHandle, isRight bool) {
	if filterNodeImplVar.conditionVar == nil {
		filterNodeImplVar.nodeLinkVar.propagateObjects(handles)
	} else {
		//TODO: rete listeners...
		var tuples []model.StreamTuple
		// tupleMap := map[model.StreamSource]model.StreamTuple{}
		if filterNodeImplVar.convert == nil {
			tuples = copyIntoTupleArray(handles)
		} else {
			tuples = make([]model.StreamTuple, len(filterNodeImplVar.convert))
			for i := 0; i < len(filterNodeImplVar.convert); i++ {
				tuples[i] = handles[filterNodeImplVar.convert[i]].getTuple()
				// tupleMap[tuples[i].GetStreamDataSource()] = tuples[i]
			}
		}
		tupleMap := convertToTupleMap(tuples)
		cv := filterNodeImplVar.conditionVar
		toPropagate := cv.getEvaluator()(cv.getName(), cv.getRule().GetName(), tupleMap)
		if toPropagate {
			filterNodeImplVar.nodeLinkVar.propagateObjects(handles)
		}
	}
}
