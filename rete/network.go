package rete

import (
	"context"
	"fmt"
	"math"

	"github.com/TIBCOSoftware/bego/common/model"

	"github.com/TIBCOSoftware/bego/utils"
)

//Network ... the rete network
type Network interface {
	AddRule(model.Rule) int
	String() string
	RemoveRule(string) model.Rule
	Assert(ctx context.Context, rs model.RuleSession, tuple model.StreamTuple)
	Retract(tuple model.StreamTuple)

	assertInternal(ctx context.Context, tuple model.StreamTuple)
	getOrCreateHandle(tuple model.StreamTuple) reteHandle
}

type reteNetworkImpl struct {
	//All rules in the network
	allRules map[string]model.Rule //(Rule)

	//Holds the DataSource name as key, and ClassNodes as value
	allClassNodes map[string]classNode //ClassNode in network

	//Holds the Rule name as key and pointer to a slice of RuleNodes as value
	ruleNameNodesOfRule map[string]utils.ArrayList //utils.ArrayList of Nodes of rule

	//Holds the Rule name as key and a pointer to a slice of NodeLinks as value
	ruleNameClassNodeLinksOfRule map[string]utils.ArrayList //utils.ArrayList of ClassNodeLink

	allHandles map[model.StreamTuple]reteHandle
}

//NewReteNetwork ... creates a new rete network
func NewReteNetwork() Network {
	reteNetworkImpl := reteNetworkImpl{}
	reteNetworkImpl.initReteNetwork()
	return &reteNetworkImpl
}

func (nw *reteNetworkImpl) initReteNetwork() {
	nw.allRules = make(map[string]model.Rule)
	nw.allClassNodes = make(map[string]classNode)
	nw.ruleNameNodesOfRule = make(map[string]utils.ArrayList)
	nw.ruleNameClassNodeLinksOfRule = make(map[string]utils.ArrayList)
	nw.allHandles = make(map[model.StreamTuple]reteHandle)
}

func (nw *reteNetworkImpl) AddRule(rule model.Rule) int {

	if nw.allRules[rule.GetName()] != nil {
		fmt.Println("Rule already exists.." + rule.GetName())
		return 1
	}
	//TODO: Worry about nonEqJoin warnings later.
	conditionSet := utils.NewArrayList()
	conditionSetNoIdr := utils.NewArrayList()
	nodeSet := utils.NewArrayList()

	nodesOfRule := utils.NewArrayList()
	classNodeLinksOfRule := utils.NewArrayList()

	if len(rule.GetConditions()) == 0 {
		identifierVar := pickIdentifier(rule.GetIdentifiers())
		nw.createClassFilterNode(rule, nodesOfRule, classNodeLinksOfRule, identifierVar, nil, nodeSet)
	} else {
		for i := 0; i < len(rule.GetConditions()); i++ {
			if rule.GetConditions()[i].GetIdentifiers() == nil || len(rule.GetConditions()[i].GetIdentifiers()) == 0 {
				//TODO: condition with no identifiers
				conditionSetNoIdr.Add(rule.GetConditions()[i])
			} else if len(rule.GetConditions()[i].GetIdentifiers()) == 1 &&
				!contains(nodeSet, rule.GetConditions()[i].GetIdentifiers()[0]) {
				cond := rule.GetConditions()[i]
				nw.createClassFilterNode(rule, nodesOfRule, classNodeLinksOfRule, cond.GetIdentifiers()[0], cond, nodeSet)
			} else {
				conditionSet.Add(rule.GetConditions()[i])
			}
		}
	}

	nw.buildNetwork(rule, nodesOfRule, classNodeLinksOfRule, conditionSet, nodeSet, conditionSetNoIdr)

	context := make([]interface{}, 2)
	context[0] = nw
	context[1] = nodesOfRule

	for key, classNode := range nw.allClassNodes {
		optimizeNetwork(key, classNode, context)
	}
	// nw.optimizeNetwork(nodesOfRule)

	nw.setClassNodeAndLinkJoinTables(nodesOfRule, classNodeLinksOfRule)

	//Add the rule to the network
	nw.allRules[rule.GetName()] = rule

	//Add RuleNodes
	nw.ruleNameNodesOfRule[rule.GetName()] = nodesOfRule

	//Add NodeLinks
	nw.ruleNameClassNodeLinksOfRule[rule.GetName()] = classNodeLinksOfRule
	return 0
}

func (nw *reteNetworkImpl) setClassNodeAndLinkJoinTables(nodesOfRule utils.ArrayList,
	classNodeLinksOfRule utils.ArrayList) {
	//TODO: add join table ids to nodes and links
}

func (nw *reteNetworkImpl) RemoveRule(ruleName string) model.Rule {

	rule := nw.allRules[ruleName]
	delete(nw.allRules, ruleName)
	if rule == nil {
		//TODO: log a message
		return nil
	}

	classNodeLinksOfRule := nw.ruleNameClassNodeLinksOfRule[ruleName].(utils.ArrayList)
	delete(nw.ruleNameClassNodeLinksOfRule, ruleName)
	if classNodeLinksOfRule != nil {
		classNodeLinksOfRule.ForEach(removeRuleHelper, nil)
	}

	nodesOfRuleItem := nw.ruleNameNodesOfRule[ruleName]
	delete(nw.ruleNameNodesOfRule, ruleName)
	if nodesOfRuleItem != nil {
		nodesOfRule := nodesOfRuleItem.(utils.ArrayList)
		for i := 0; i < nodesOfRule.Len(); i++ {
			node := nodesOfRule.Get(i).(abstractNode)
			switch nodeImpl := node.(type) {
			//Only interested in joinnodes
			//case *filterNodeImpl:
			//case *classNodeImpl:
			//case *ruleNodeImpl:
			case *joinNodeImpl:
				removeRefsFromReteHandles(nodeImpl.leftTable)
				removeRefsFromReteHandles(nodeImpl.rightTable)
			}
		}
	}

	return rule
}

func removeRefsFromReteHandles(joinTableVar joinTable) {
	if joinTableVar == nil {
		return
	}
	for tableRow := range joinTableVar.getMap() {
		for _, handle := range tableRow.getHandles() {
			handle.removeJoinTable(joinTableVar)
		}
	}
}

func removeRuleHelper(entry interface{}, context []interface{}) {
	classNodeLinkOfRule := entry.(classNodeLink)
	classNodeLinkOfRule.getClassNode().removeClassNodeLink(classNodeLinkOfRule)
}

func optimizeNetwork(key string, classNodeVar classNode, context []interface{}) {
	nodesOfRule := context[1].(utils.ArrayList)
	for j := 0; j < classNodeVar.getClassNodeLinks().Len(); j++ {
		nodeLink := classNodeVar.getClassNodeLinks().Get(j).(classNodeLink)
		childNode := nodeLink.getChild()
		switch nodeImpl := childNode.(type) {
		case *filterNodeImpl:
			if nodeImpl.conditionVar == nil {
				nodeLink.setChild(nodeImpl.nodeLinkVar.getChild())
				nodeLink.setIsRightChild(nodeImpl.nodeLinkVar.isRightNode())
				nodesOfRule.Remove(nodeImpl)
			}
		}
	}
}

func contains(nodeSet utils.ArrayList, identifierVar model.TupleTypeAlias) bool {
	identifiers := []model.TupleTypeAlias{identifierVar}
	for i := 0; i < nodeSet.Len(); i++ {
		node := nodeSet.Get(i).(node)
		if ContainedByFirst(node.getIdentifiers(), identifiers) {
			return true
		}
	}
	return false
}

func (nw *reteNetworkImpl) buildNetwork(rule model.Rule, nodesOfRule utils.ArrayList, classNodeLinksOfRule utils.ArrayList,
	conditionSet utils.ArrayList, nodeSet utils.ArrayList, conditionSetNoIdr utils.ArrayList) {
	if conditionSet.Len() == 0 {
		if nodeSet.Len() == 1 {
			node := nodeSet.Get(0).(node)
			if ContainedByFirst(node.getIdentifiers(), rule.GetIdentifiers()) {
				//TODO: Re evaluate set later..

				lastNode := node
				//check conditions with no identifierVar
				for i := 0; i < conditionSetNoIdr.Len(); i++ {
					conditionVar := conditionSetNoIdr.Get(i).(model.Condition)
					fNode := newFilterNode(node.getIdentifiers(), conditionVar)
					nodesOfRule.Add(fNode)
					newNodeLink(lastNode, fNode, false)
					lastNode = fNode
				}
				//Yoohoo! We have a Rule!!
				ruleNode := newRuleNode(rule)
				newNodeLink(node, ruleNode, false)
				nodesOfRule.Add(ruleNode)
			} else {
				idrs := SecondMinusFirst(node.getIdentifiers(), rule.GetIdentifiers())
				fNode := nw.createClassFilterNode(rule, nodesOfRule, classNodeLinksOfRule, idrs[0], nil, nodeSet)
				nw.createJoinNode(rule, nodesOfRule, node, fNode, nil, conditionSet, nodeSet)
				nw.buildNetwork(rule, nodesOfRule, classNodeLinksOfRule, conditionSet, nodeSet, conditionSetNoIdr)
			}
		} else {
			nodes := findSimilarNodes(nodeSet)
			nw.createJoinNode(rule, nodesOfRule, nodes[0], nodes[1], nil, conditionSet, nodeSet)
			nw.buildNetwork(rule, nodesOfRule, classNodeLinksOfRule, conditionSet, nodeSet, conditionSetNoIdr)
		}
	} else {
		if nw.createFilterNode(rule, nodesOfRule, conditionSet, nodeSet) {
			nw.buildNetwork(rule, nodesOfRule, classNodeLinksOfRule, conditionSet, nodeSet, conditionSetNoIdr)
		} else if nw.createJoinNodeFromExisting(rule, nodesOfRule, conditionSet, nodeSet) {
			nw.buildNetwork(rule, nodesOfRule, classNodeLinksOfRule, conditionSet, nodeSet, conditionSetNoIdr)
		} else if nw.createJoinNodeFromSome(rule, nodesOfRule, classNodeLinksOfRule, conditionSet, nodeSet) {
			nw.buildNetwork(rule, nodesOfRule, classNodeLinksOfRule, conditionSet, nodeSet, conditionSetNoIdr)
		} else {
			conditionVar := nw.findConditionWithLeastIdentifiers(conditionSet)
			nw.createClassFilterNode(rule, nodesOfRule, classNodeLinksOfRule, conditionVar.GetIdentifiers()[0], nil, nodeSet)
			nw.buildNetwork(rule, nodesOfRule, classNodeLinksOfRule, conditionSet, nodeSet, conditionSetNoIdr)
		}
	}
}

func (nw *reteNetworkImpl) createFilterNode(rule model.Rule, nodesOfRule utils.ArrayList, conditionSet utils.ArrayList, nodeSet utils.ArrayList) bool {
	for i := 0; i < conditionSet.Len(); i++ {
		conditionVar := conditionSet.Get(i).(model.Condition)
		for i := 0; i < nodeSet.Len(); i++ {
			node := nodeSet.Get(i).(node)
			if ContainedByFirst(node.getIdentifiers(), conditionVar.GetIdentifiers()) {
				//TODO
				filterNode := newFilterNode(nil, conditionVar)
				newNodeLink(node, filterNode, false)
				nodeSet.Remove(node)
				nodeSet.Add(filterNode)
				nodesOfRule.Add(filterNode)
				return true
			}
		}
	}

	return false
}

func (nw *reteNetworkImpl) createJoinNodeFromExisting(rule model.Rule, nodesOfRule utils.ArrayList, conditionSet utils.ArrayList, nodeSet utils.ArrayList) bool {
	maxCommonIdr := -1
	numOfIdentifiers := 0
	joinThese := make([]node, 2)
	var targetCondition model.Condition
	for i := 0; i < conditionSet.Len(); i++ {
		conditionVar := conditionSet.Get(i).(model.Condition)
		for j := 0; j < nodeSet.Len(); j++ {
			leftNode := nodeSet.Get(j).(node)
			for k := j + 1; k < nodeSet.Len(); k++ {
				rightNode := nodeSet.Get(k).(node)
				if OtherTwoAreContainedByFirst(conditionVar.GetIdentifiers(), leftNode.getIdentifiers(), rightNode.getIdentifiers()) {
					commonIdr := len(IntersectionIdentifiers(leftNode.getIdentifiers(), rightNode.getIdentifiers()))
					if maxCommonIdr < commonIdr {
						maxCommonIdr = commonIdr
						joinThese[0] = leftNode
						joinThese[1] = rightNode
						targetCondition = conditionVar
						numOfIdentifiers = len(UnionIdentifiers(leftNode.getIdentifiers(), rightNode.getIdentifiers()))
					} else if maxCommonIdr == commonIdr {
						numIdrs := len(UnionIdentifiers(leftNode.getIdentifiers(), rightNode.getIdentifiers()))
						if numIdrs < numOfIdentifiers {
							joinThese[0] = leftNode
							joinThese[1] = rightNode
							targetCondition = conditionVar
							numOfIdentifiers = numIdrs
						}
					}
				}
			}
		}
		if maxCommonIdr != -1 {
			nw.createJoinNode(rule, nodesOfRule, joinThese[0], joinThese[1], targetCondition, conditionSet, nodeSet)
			return true
		}
	}

	return false
}

func (nw *reteNetworkImpl) createJoinNodeFromSome(rule model.Rule, nodesOfRule utils.ArrayList,
	classNodeLinksOfRule utils.ArrayList, conditionSet utils.ArrayList, nodeSet utils.ArrayList) bool {
	leastNeeded := math.MaxUint32
	maxIdentifier := -1
	var targetNode node
	var targetCondition model.Condition
	for i := 0; i < conditionSet.Len(); i++ {
		conditionVar := conditionSet.Get(i).(model.Condition)
		for j := 0; j < nodeSet.Len(); j++ {
			nodeIdentifiers := nodeSet.Get(j).(node).getIdentifiers()
			need := len(SecondMinusFirst(nodeIdentifiers, conditionVar.GetIdentifiers()))
			if need < leastNeeded {
				leastNeeded = need
				maxIdentifier = len(nodeIdentifiers)
				targetNode = nodeSet.Get(j).(node)
				targetCondition = conditionVar
			} else if need == leastNeeded {
				if len(nodeIdentifiers) > maxIdentifier {
					maxIdentifier = len(nodeIdentifiers)
					targetNode = nodeSet.Get(j).(node)
					targetCondition = conditionVar
				}
			}
		}
	}
	if maxIdentifier == -1 {
		return false
	}
	nodeIdentifiers := SecondMinusFirst(targetNode.getIdentifiers(), targetCondition.GetIdentifiers())
	if leastNeeded == 1 {
		filterNode := nw.createClassFilterNode(rule, nodesOfRule, classNodeLinksOfRule, nodeIdentifiers[0], nil, nodeSet)
		nw.createJoinNode(rule, nodesOfRule, targetNode, filterNode, targetCondition, conditionSet, nodeSet)
	} else {
		useThis := findBestNode(nodeSet, nodeIdentifiers, targetNode)
		if useThis == nil {
			nw.createClassFilterNode(rule, nodesOfRule, classNodeLinksOfRule, nodeIdentifiers[0], nil, nodeSet)
		} else {
			nw.createJoinNode(rule, nodesOfRule, targetNode, useThis, nil, conditionSet, nodeSet)
		}
	}

	return true
}

func (nw *reteNetworkImpl) createClassFilterNode(rule model.Rule, nodesOfRule utils.ArrayList, classNodeLinksOfRule utils.ArrayList, identifierVar model.TupleTypeAlias, conditionVar model.Condition, nodeSet utils.ArrayList) filterNode {
	identifiers := []model.TupleTypeAlias{identifierVar}
	classNodeVar := getClassNode(nw, identifierVar)
	filterNodeVar := newFilterNode(identifiers, conditionVar)
	classNodeLink := newClassNodeLink(classNodeVar, filterNodeVar, rule, identifierVar)
	classNodeVar.addClassNodeLink(classNodeLink)
	nodesOfRule.Add(classNodeVar)
	nodesOfRule.Add(filterNodeVar)
	//TODO: Add to RuleLinks
	classNodeLinksOfRule.Add(classNodeLink)
	nodeSet.Add(filterNodeVar)
	return filterNodeVar
}

func (nw *reteNetworkImpl) createJoinNode(rule model.Rule, nodesOfRule utils.ArrayList, leftNode node, rightNode node, joinCondition model.Condition, conditionSet utils.ArrayList, nodeSet utils.ArrayList) {

	//TODO handle equivJoins later..

	joinNode := newJoinNode(leftNode.getIdentifiers(), rightNode.getIdentifiers(), joinCondition)

	newNodeLink(leftNode, joinNode, false)
	newNodeLink(rightNode, joinNode, true)
	nodeSet.Remove(leftNode)
	nodeSet.Remove(rightNode)
	nodeSet.Add(joinNode)
	nodesOfRule.Add(joinNode)
	if joinCondition != nil {
		conditionSet.Remove(joinCondition)
	}
}

func findBestNode(nodeSet utils.ArrayList, matchIdentifiers []model.TupleTypeAlias, notThis node) node {
	var foundNode node
	foundNode = nil
	foundIdr := 0

	for i := 0; i < nodeSet.Len(); i++ {
		node := nodeSet.Get(i).(node)
		if node == notThis {
			continue
		}
		foundMatch := len(IntersectionIdentifiers(node.getIdentifiers(), matchIdentifiers))
		if foundMatch > foundIdr {
			foundIdr = foundMatch
			foundNode = node
		}
	}
	return foundNode
}

func (nw *reteNetworkImpl) findConditionWithLeastIdentifiers(conditionSet utils.ArrayList) model.Condition {
	least := math.MaxUint16
	var leastIdentifiers model.Condition
	for i := 0; i < conditionSet.Len(); i++ {
		c := conditionSet.Get(i).(model.Condition)
		lenIdr := len(c.GetIdentifiers())
		if lenIdr < least {
			leastIdentifiers = c
			least = lenIdr
		}
	}
	if least == math.MaxUint16 {
		return nil
	}
	return leastIdentifiers
}

func getClassNode(nw *reteNetworkImpl, name model.TupleTypeAlias) classNode {
	classNodeVar := nw.allClassNodes[string(name)]
	if classNodeVar == nil {
		classNodeVar = newClassNode(string(name))
		nw.allClassNodes[string(name)] = classNodeVar
	}
	return classNodeVar
}

func (nw *reteNetworkImpl) String() string {

	str := "\n>>> Class View <<<\n"

	for _, classNodeImpl := range nw.allClassNodes {
		str += classNodeImpl.String() + "\n"
	}
	str += ">>>> Rule View <<<<\n"

	for _, rule := range nw.allRules {
		str += nw.PrintRule(rule)
	}

	return str
}

func pickIdentifier(idrs []model.TupleTypeAlias) model.TupleTypeAlias {
	return idrs[0]
}

func (nw *reteNetworkImpl) PrintRule(rule model.Rule) string {
	//str := "[Rule (" + rule.GetName() + ") Id(" + strconv.Itoa(rule.GetID()) + ")]\n"
	str := "[Rule (" + rule.GetName() + ") Id()]\n"

	nodesOfRule := nw.ruleNameNodesOfRule[rule.GetName()]

	for i := 0; i < nodesOfRule.Len(); i++ {
		node := nodesOfRule.Get(i).(abstractNode)
		switch nodeImpl := node.(type) {
		case *filterNodeImpl:
			str += nodeImpl.String()
		case *joinNodeImpl:
			str += nodeImpl.String()
		case *classNodeImpl:
			str += nw.printClassNode(rule.GetName(), nodeImpl)
		case *ruleNodeImpl:
			str += nodeImpl.String()
		}
		str += "\n"
	}
	return str
}

func (nw *reteNetworkImpl) printClassNode(ruleName string, classNodeImpl *classNodeImpl) string {
	classNodesLinksOfRule := nw.ruleNameClassNodeLinksOfRule[ruleName].(utils.ArrayList)
	links := ""
	for i := 0; i < classNodesLinksOfRule.Len(); i++ {
		classNodeLinkOfRule := classNodesLinksOfRule.Get(i).(classNodeLink)
		if string(classNodeLinkOfRule.GetIdentifier()) == classNodeImpl.name {
			links += "\n\t\t" + classNodeLinkOfRule.String()
		}
	}
	return "\t[ClassNode Class(" + classNodeImpl.getName() + ")" + links + "]\n"
}

func (nw *reteNetworkImpl) Assert(ctx context.Context, rs model.RuleSession, tuple model.StreamTuple) {

	if ctx == nil {
		ctx = context.Background()
	}

	reteCtxVar, isRecursive, newCtx := getOrSetReteCtx(ctx, nw, rs)

	if !isRecursive {
		nw.assertInternal(newCtx, tuple)
	} else {
		reteCtxVar.getOpsList().PushBack(newAssertEntry(tuple))
	}

	reteCtxVar.getConflictResolver().resolveConflict(newCtx)
}

func (nw *reteNetworkImpl) Retract(tuple model.StreamTuple) {
	reteHandle := nw.allHandles[tuple]
	if reteHandle != nil {
		reteHandle.removeJoinTableRowRefs()
	}
}

func (nw *reteNetworkImpl) assertInternal(ctx context.Context, tuple model.StreamTuple) {
	dataSource := tuple.GetTypeAlias()
	listItem := nw.allClassNodes[string(dataSource)]
	if listItem != nil {
		classNodeVar := listItem.(classNode)
		classNodeVar.assert(ctx, tuple)
	} else {
		fmt.Println("No rule exists for data stream: " + dataSource)
	}
}

func (nw *reteNetworkImpl) getOrCreateHandle(tuple model.StreamTuple) reteHandle {
	h := nw.allHandles[tuple]
	if h == nil {
		h1 := handleImpl{}
		h1.initHandleImpl()
		h1.setTuple(tuple)
		h = &h1
		nw.allHandles[tuple] = h
	}
	return h
}
