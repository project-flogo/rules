package rete

import (
	"context"
	"strconv"
)

//joinNode holds the join tables for unmatched entries
type joinNode interface {
	node
}

type joinNodeImpl struct {
	nodeImpl
	conditionVar condition

	leftIdrs  []identifier
	rightIdrs []identifier

	leftIdrLen  int
	rightIdrLen int
	totalIdrLen int

	joinIndexForLeft  []int
	joinIndexForRight []int

	leftTable  joinTable
	rightTable joinTable
}

func newJoinNode(leftIdrs []identifier, rightIdrs []identifier, conditionVar condition) joinNode {
	jn := joinNodeImpl{}
	jn.initjoinNodeImplVar(leftIdrs, rightIdrs, conditionVar)
	return &jn
}

func (jn *joinNodeImpl) initjoinNodeImplVar(leftIdrs []identifier, rightIdrs []identifier, conditionVar condition) {
	jn.initNodeImpl(nil)
	jn.leftIdrs = leftIdrs
	jn.rightIdrs = rightIdrs
	jn.conditionVar = conditionVar
	jn.leftTable = newJoinTable(leftIdrs)
	jn.rightTable = newJoinTable(rightIdrs)
	jn.setJoinIdentifiers()
}

func (jn *joinNodeImpl) GetLeftIdentifiers() []identifier {
	return jn.leftIdrs
}

func (jn *joinNodeImpl) GetRightIdentifiers() []identifier {
	return jn.rightIdrs
}

func (jn *joinNodeImpl) setJoinIdentifiers() {
	jn.leftIdrLen = len(jn.leftIdrs)
	jn.rightIdrLen = len(jn.rightIdrs)
	jn.totalIdrLen = jn.leftIdrLen + jn.rightIdrLen

	jn.identifiers = make([]identifier, jn.totalIdrLen)

	jn.joinIndexForLeft = make([]int, jn.leftIdrLen)
	jn.joinIndexForRight = make([]int, jn.rightIdrLen)

	for i := 0; i < jn.leftIdrLen; i++ {
		jn.joinIndexForLeft[i] = -1
	}
	for i := 0; i < jn.rightIdrLen; i++ {
		jn.joinIndexForRight[i] = -1
	}
	conditionIdrLen := 0
	if jn.conditionVar != nil {
		conditionIdrLen = len(jn.conditionVar.getIdentifiers())
		for i := 0; i < conditionIdrLen; i++ {
			idx := GetIndex(jn.leftIdrs, jn.conditionVar.getIdentifiers()[i])
			if idx != -1 {
				jn.joinIndexForLeft[idx] = i
				jn.identifiers[i] = jn.leftIdrs[idx]
				continue
			}
			idx = GetIndex(jn.rightIdrs, jn.conditionVar.getIdentifiers()[i])
			if idx != -1 {
				jn.joinIndexForRight[idx] = i
				jn.identifiers[i] = jn.rightIdrs[idx]
				continue
			}
			//TODO ERROR HANDLING!
		}
	}

	outIndex := conditionIdrLen
	for i := 0; i < jn.leftIdrLen; i++ {
		if jn.joinIndexForLeft[i] == -1 {
			jn.joinIndexForLeft[i] = outIndex
			jn.identifiers[outIndex] = jn.leftIdrs[i]
			outIndex++
		}
	}
	for i := 0; i < jn.rightIdrLen; i++ {
		if jn.joinIndexForRight[i] == -1 {
			jn.joinIndexForRight[i] = outIndex
			jn.identifiers[outIndex] = jn.rightIdrs[i]
			outIndex++
		}
	}

	if outIndex != jn.totalIdrLen {
		//TODO ERROR HANDLING!
	}
}

//String Stringer.String interface
func (jn *joinNodeImpl) String() string {

	joinIdsForLeftStr := ""
	for i := range jn.joinIndexForLeft {
		joinIdsForLeftStr += strconv.Itoa(i) + ", "
	}

	joinIdsForRightStr := ""
	for i := range jn.joinIndexForRight {
		joinIdsForRightStr += strconv.Itoa(i) + ", "
	}

	linkTo := ""
	switch jn.nodeLinkVar.getChild().(type) {
	case *joinNodeImpl:
		if jn.nodeLinkVar.isRightNode() {
			linkTo += strconv.Itoa(jn.nodeLinkVar.getChild().getID()) + "R"
		} else {
			linkTo += strconv.Itoa(jn.nodeLinkVar.getChild().getID()) + "L"
		}
	default:
		linkTo += strconv.Itoa(jn.nodeLinkVar.getChild().getID())
	}

	joinConditionStr := "nil"
	joinConditionIdrsStr := "nil"
	if jn.conditionVar != nil {
		joinConditionStr = jn.conditionVar.String()
		joinConditionIdrsStr = IdentifiersToString(jn.conditionVar.getIdentifiers())
	}
	return "\t[JoinNode(" + jn.nodeImpl.String() + ") link(" + linkTo + ")\n" +
		"\t\tLeft identifier      = " + IdentifiersToString(jn.leftIdrs) + ";\n" +
		"\t\tRight identifier     = " + IdentifiersToString(jn.rightIdrs) + ";\n" +
		"\t\tOut identifier       = " + IdentifiersToString(jn.identifiers) + ";\n" +
		"\t\tCondition identifier = " + joinConditionIdrsStr + ";\n" +
		"\t\tJoin Left Index      = " + joinIdsForLeftStr + ";\n" +
		"\t\tJoin Right Index     = " + joinIdsForRightStr + ";\n" +
		"\t\tCondition            = " + joinConditionStr + "]\n"
}

func (jn *joinNodeImpl) assertObjects(ctx context.Context, handles []reteHandle, isRight bool) {
	//TODO:
	joinedHandles := make([]reteHandle, jn.totalIdrLen)
	if isRight {
		jn.assertFromRight(ctx, handles, joinedHandles)
	} else {
		jn.assertFromLeft(ctx, handles, joinedHandles)
	}
}

func (jn *joinNodeImpl) assertFromRight(ctx context.Context, handles []reteHandle, joinedHandles []reteHandle) {

	//TODO: other stuff. right now focus on tuple table
	jn.joinRightObjects(handles, joinedHandles)
	tupleTableRow := newJoinTableRow(handles)
	jn.rightTable.addRow(tupleTableRow)
	//TODO: rete listeners etc.
	for tupleTableRowLeft := range jn.leftTable.getMap() {
		success := jn.joinLeftObjects(tupleTableRowLeft.getHandles(), joinedHandles)
		if !success {
			//TODO: handle it
			continue
		}
		toPropagate := false
		if jn.conditionVar == nil {
			toPropagate = true
		} else {
			tupleMap := copyIntoTupleMap(joinedHandles)
			cv := jn.conditionVar
			toPropagate = cv.getEvaluator()(cv.getName(), cv.getRule().GetName(), tupleMap)
		}
		if toPropagate {
			jn.nodeLinkVar.propagateObjects(ctx, joinedHandles)
		}
	}
}

func (jn *joinNodeImpl) joinLeftObjects(leftHandles []reteHandle, joinedHandles []reteHandle) bool {
	for i := 0; i < jn.leftIdrLen; i++ {
		handle := leftHandles[i]
		if handle.getTuple() == nil {
			return false
		}
		joinedHandles[jn.joinIndexForLeft[i]] = handle
	}
	return true
}

func (jn *joinNodeImpl) joinRightObjects(rightHandles []reteHandle, joinedHandles []reteHandle) bool {
	for i := 0; i < jn.rightIdrLen; i++ {
		handle := rightHandles[i]
		if handle.getTuple() == nil {
			return false
		}
		joinedHandles[jn.joinIndexForRight[i]] = handle
	}
	return true
}

func (jn *joinNodeImpl) assertFromLeft(ctx context.Context, handles []reteHandle, joinedHandles []reteHandle) {
	jn.joinLeftObjects(handles, joinedHandles)
	//TODO: other stuff. right now focus on tuple table
	tupleTableRow := newJoinTableRow(handles)
	jn.leftTable.addRow(tupleTableRow)
	//TODO: rete listeners etc.
	for tupleTableRowRight := range jn.rightTable.getMap() {
		success := jn.joinRightObjects(tupleTableRowRight.getHandles(), joinedHandles)
		if !success {
			//TODO: handle it
			continue
		}
		toPropagate := false
		if jn.conditionVar == nil {
			toPropagate = true
		} else {
			tupleMap := copyIntoTupleMap(joinedHandles)
			cv := jn.conditionVar
			toPropagate = cv.getEvaluator()(cv.getName(), cv.getRule().GetName(), tupleMap)
		}
		if toPropagate {
			jn.nodeLinkVar.propagateObjects(ctx, joinedHandles)
		}
	}
}
