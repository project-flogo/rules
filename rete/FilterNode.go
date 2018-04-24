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

// void assertObjects(Handle[] handles, Object[] objects, boolean right) {
// 	if(m_wm.isLoadingObjects() && loadStopHere) return;  //if it is in loadObject mode, return
// 	if (m_condition == null){
// 		propagateObjects(handles, objects);
// 		return;
// 	}
// 	else {
// 		boolean evalSuccess = false;
// 		if(m_wm.m_reteListener != null) {
// 			m_wm.m_reteListener.filterConditionStart(this);
// 		}
// 		try {
// 			if(m_convertIndex == null) {
// 				Object[] _objects = objects;
// 				boolean _eval = m_condition.eval(_objects);
// 				if(_eval) {
// 					evalSuccess = true;
// 					propagateObjects(handles, objects);
// 				}
// 			}
// 			else {
// 				convert(handles, objects);
// 				Handle [] convertedHandles = new Handle[m_convertIndex.length];
// 				Object [] convertedObjects = new Object[m_convertIndex.length];

// 				for(int i=0; i < m_convertIndex.length; i++) {
// 					convertedObjects[i] = objects[m_convertIndex[i]];
// 					convertedHandles[i] = handles[m_convertIndex[i]];
// 				}

// 				Object[] _objects = convertedObjects;
// 				boolean _eval = m_condition.eval(_objects);
// 				if(_eval) {
// 					evalSuccess = true;
// 					propagateObjects(handles, objects);
// 				}
// 			}
// 		}
// 		catch(ForceConditionFailureException fcfe) {
// 			m_logger.log(Level.DEBUG,"Forced condition failure in " + getWorkingMemory().getName() + " : " + m_condition.getRule().getName() + " : " + m_condition + " with " + Format.objsToStr(objects));
// 		}
// 		catch(RuntimeException ex) {
// 			String errMsg =ResourceManager.formatString("rule.condition.exception", m_condition.getRule().getName(), m_condition, Format.objsToStr(objects));
// 			m_logger.log(Level.ERROR,errMsg, ex);
// 		}
// 		if(m_wm.m_reteListener != null)
// 			m_wm.m_reteListener.filterConditionEnd(evalSuccess);
// 	}
// }
