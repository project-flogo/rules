package rete

import (
	"strconv"
)

//joinNode holds the join tables for unmatched entries
type joinNode interface {
	node
}

type joinNodeImpl struct {
	nodeImpl
	conditionVar condition

	leftIdrs  []Identifier
	rightIdrs []Identifier

	leftIdrLen  int
	rightIdrLen int
	totalIdrLen int

	joinIndexForLeft  []int
	joinIndexForRight []int

	leftTable  joinTable
	rightTable joinTable
}

func newJoinNode(leftIdrs []Identifier, rightIdrs []Identifier, conditionVar condition) joinNode {
	joinNodeImplVar := joinNodeImpl{}
	joinNodeImplVar.initjoinNodeImplVar(leftIdrs, rightIdrs, conditionVar)
	return &joinNodeImplVar
}

func (joinNodeImplVar *joinNodeImpl) initjoinNodeImplVar(leftIdrs []Identifier, rightIdrs []Identifier, conditionVar condition) {
	joinNodeImplVar.initNodeImpl(nil)
	joinNodeImplVar.leftIdrs = leftIdrs
	joinNodeImplVar.rightIdrs = rightIdrs
	joinNodeImplVar.conditionVar = conditionVar
	joinNodeImplVar.leftTable = newJoinTable(leftIdrs)
	joinNodeImplVar.rightTable = newJoinTable(rightIdrs)
	joinNodeImplVar.setJoinIdentifiers()
}

func (joinNodeImplVar *joinNodeImpl) GetLeftIdentifiers() []Identifier {
	return joinNodeImplVar.leftIdrs
}

func (joinNodeImplVar *joinNodeImpl) GetRightIdentifiers() []Identifier {
	return joinNodeImplVar.rightIdrs
}

func (joinNodeImplVar *joinNodeImpl) setJoinIdentifiers() {
	joinNodeImplVar.leftIdrLen = len(joinNodeImplVar.leftIdrs)
	joinNodeImplVar.rightIdrLen = len(joinNodeImplVar.rightIdrs)
	joinNodeImplVar.totalIdrLen = joinNodeImplVar.leftIdrLen + joinNodeImplVar.rightIdrLen

	joinNodeImplVar.identifiers = make([]Identifier, joinNodeImplVar.totalIdrLen)

	joinNodeImplVar.joinIndexForLeft = make([]int, joinNodeImplVar.leftIdrLen)
	joinNodeImplVar.joinIndexForRight = make([]int, joinNodeImplVar.rightIdrLen)

	for i := 0; i < joinNodeImplVar.leftIdrLen; i++ {
		joinNodeImplVar.joinIndexForLeft[i] = -1
	}
	for i := 0; i < joinNodeImplVar.rightIdrLen; i++ {
		joinNodeImplVar.joinIndexForRight[i] = -1
	}
	conditionIdrLen := 0
	if joinNodeImplVar.conditionVar != nil {
		conditionIdrLen = len(joinNodeImplVar.conditionVar.getIdentifiers())
		for i := 0; i < conditionIdrLen; i++ {
			idx := GetIndex(joinNodeImplVar.leftIdrs, joinNodeImplVar.conditionVar.getIdentifiers()[i])
			if idx != -1 {
				joinNodeImplVar.joinIndexForLeft[idx] = i
				joinNodeImplVar.identifiers[i] = joinNodeImplVar.leftIdrs[idx]
				continue
			}
			idx = GetIndex(joinNodeImplVar.rightIdrs, joinNodeImplVar.conditionVar.getIdentifiers()[i])
			if idx != -1 {
				joinNodeImplVar.joinIndexForRight[idx] = i
				joinNodeImplVar.identifiers[i] = joinNodeImplVar.rightIdrs[idx]
				continue
			}
			//TODO ERROR HANDLING!
		}
	}

	outIndex := conditionIdrLen
	for i := 0; i < joinNodeImplVar.leftIdrLen; i++ {
		if joinNodeImplVar.joinIndexForLeft[i] == -1 {
			joinNodeImplVar.joinIndexForLeft[i] = outIndex
			joinNodeImplVar.identifiers[outIndex] = joinNodeImplVar.leftIdrs[i]
			outIndex++
		}
	}
	for i := 0; i < joinNodeImplVar.rightIdrLen; i++ {
		if joinNodeImplVar.joinIndexForRight[i] == -1 {
			joinNodeImplVar.joinIndexForRight[i] = outIndex
			joinNodeImplVar.identifiers[outIndex] = joinNodeImplVar.rightIdrs[i]
			outIndex++
		}
	}

	if outIndex != joinNodeImplVar.totalIdrLen {
		//TODO ERROR HANDLING!
	}
}

//String Stringer.String interface
func (joinNodeImplVar *joinNodeImpl) String() string {

	joinIdsForLeftStr := ""
	for i := range joinNodeImplVar.joinIndexForLeft {
		joinIdsForLeftStr += strconv.Itoa(i) + ", "
	}

	joinIdsForRightStr := ""
	for i := range joinNodeImplVar.joinIndexForRight {
		joinIdsForRightStr += strconv.Itoa(i) + ", "
	}

	linkTo := ""
	switch joinNodeImplVar.nodeLinkVar.getChild().(type) {
	case *joinNodeImpl:
		if joinNodeImplVar.nodeLinkVar.isRightNode() {
			linkTo += strconv.Itoa(joinNodeImplVar.nodeLinkVar.getChild().getID()) + "R"
		} else {
			linkTo += strconv.Itoa(joinNodeImplVar.nodeLinkVar.getChild().getID()) + "L"
		}
	default:
		linkTo += strconv.Itoa(joinNodeImplVar.nodeLinkVar.getChild().getID())
	}

	joinConditionStr := "nil"
	joinConditionIdrsStr := "nil"
	if joinNodeImplVar.conditionVar != nil {
		joinConditionStr = joinNodeImplVar.conditionVar.String()
		joinConditionIdrsStr = IdentifiersToString(joinNodeImplVar.conditionVar.getIdentifiers())
	}
	return "\t[JoinNode(" + joinNodeImplVar.nodeImpl.String() + ") link(" + linkTo + ")\n" +
		"\t\tLeft Identifier      = " + IdentifiersToString(joinNodeImplVar.leftIdrs) + ";\n" +
		"\t\tRight Identifier     = " + IdentifiersToString(joinNodeImplVar.rightIdrs) + ";\n" +
		"\t\tOut Identifier       = " + IdentifiersToString(joinNodeImplVar.identifiers) + ";\n" +
		"\t\tCondition Identifier = " + joinConditionIdrsStr + ";\n" +
		"\t\tJoin Left Index      = " + joinIdsForLeftStr + ";\n" +
		"\t\tJoin Right Index     = " + joinIdsForRightStr + ";\n" +
		"\t\tCondition            = " + joinConditionStr + "]\n"
}

func (joinNodeImplVar *joinNodeImpl) assertObjects(handles []reteHandle, isRight bool) {
	//TODO:
	joinedHandles := make([]reteHandle, joinNodeImplVar.totalIdrLen)
	if isRight {
		joinNodeImplVar.assertFromRight(handles, joinedHandles)
	} else {
		joinNodeImplVar.assertFromLeft(handles, joinedHandles)
	}
}

func (joinNodeImplVar *joinNodeImpl) assertFromRight(handles []reteHandle, joinedHandles []reteHandle) {
	//TODO: other stuff. right now focus on tuple table
	joinNodeImplVar.joinRightObjects(handles, joinedHandles)
	tupleTableRow := newJoinTableRow(handles)
	joinNodeImplVar.rightTable.addRow(tupleTableRow)
	//TODO: rete listeners etc.
	for tupleTableRowLeft := range joinNodeImplVar.leftTable.getMap() {
		success := joinNodeImplVar.joinLeftObjects(tupleTableRowLeft.getHandles(), joinedHandles)
		if !success {
			//TODO: handle it
			continue
		}
		toPropagate := false
		if joinNodeImplVar.conditionVar == nil {
			toPropagate = true
		} else {
			tupleMap := copyIntoTupleMap(joinedHandles)
			cv := joinNodeImplVar.conditionVar
			toPropagate = cv.getEvaluator()(cv.getName(), cv.getRule().GetName(), tupleMap)
		}
		if toPropagate {
			joinNodeImplVar.nodeLinkVar.propagateObjects(joinedHandles)
		}
	}
}

func (joinNodeImplVar *joinNodeImpl) joinLeftObjects(leftHandles []reteHandle, joinedHandles []reteHandle) bool {
	for i := 0; i < joinNodeImplVar.leftIdrLen; i++ {
		handle := leftHandles[i]
		if handle.getTuple() == nil {
			return false
		}
		joinedHandles[joinNodeImplVar.joinIndexForLeft[i]] = handle
	}
	return true
}

func (joinNodeImplVar *joinNodeImpl) joinRightObjects(rightHandles []reteHandle, joinedHandles []reteHandle) bool {
	for i := 0; i < joinNodeImplVar.rightIdrLen; i++ {
		handle := rightHandles[i]
		if handle.getTuple() == nil {
			return false
		}
		joinedHandles[joinNodeImplVar.joinIndexForRight[i]] = handle
	}
	return true
}

func (joinNodeImplVar *joinNodeImpl) assertFromLeft(handles []reteHandle, joinedHandles []reteHandle) {
	joinNodeImplVar.joinLeftObjects(handles, joinedHandles)
	//TODO: other stuff. right now focus on tuple table
	tupleTableRow := newJoinTableRow(handles)
	joinNodeImplVar.leftTable.addRow(tupleTableRow)
	//TODO: rete listeners etc.
	for tupleTableRowRight := range joinNodeImplVar.rightTable.getMap() {
		success := joinNodeImplVar.joinRightObjects(tupleTableRowRight.getHandles(), joinedHandles)
		if !success {
			//TODO: handle it
			continue
		}
		toPropagate := false
		if joinNodeImplVar.conditionVar == nil {
			toPropagate = true
		} else {
			tupleMap := copyIntoTupleMap(joinedHandles)
			cv := joinNodeImplVar.conditionVar
			toPropagate = cv.getEvaluator()(cv.getName(), cv.getRule().GetName(), tupleMap)
		}
		if toPropagate {
			joinNodeImplVar.nodeLinkVar.propagateObjects(joinedHandles)
		}
	}
}
