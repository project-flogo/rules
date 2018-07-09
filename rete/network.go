package rete

import (
	"fmt"
	"math"
	"strconv"

	"github.com/TIBCOSoftware/bego/common/model"

	"github.com/TIBCOSoftware/bego/utils"
)

//Network ... the rete network
type Network interface {
	AddRule(Rule) int
	String() string
	RemoveRule(string) Rule
	Assert(tuple model.StreamTuple)
	Retract(tuple model.StreamTuple)
}

type reteNetworkImpl struct {
	//All rules in the network
	allRules utils.Map //(Rule)

	//Holds the DataSource name as key, and ClassNodes as value
	allClassNodes utils.Map //ClassNode in network

	//Holds the Rule name as key and pointer to a slice of RuleNodes as value
	ruleNameNodesOfRule utils.Map //utils.ArrayList of Nodes of rule

	//Holds the Rule name as key and a pointer to a slice of NodeLinks as value
	ruleNameClassNodeLinksOfRule utils.Map //utils.ArrayList of ClassNodeLink
}

//NewReteNetwork ... creates a new rete network
func NewReteNetwork() Network {
	reteNetworkImpl := reteNetworkImpl{}
	reteNetworkImpl.initReteNetwork()
	return &reteNetworkImpl
}

func (nw *reteNetworkImpl) initReteNetwork() {
	nw.allRules = utils.NewHashMap()
	nw.allClassNodes = utils.NewHashMap()
	nw.ruleNameNodesOfRule = utils.NewHashMap()
	nw.ruleNameClassNodeLinksOfRule = utils.NewHashMap()
}

func (nw *reteNetworkImpl) AddRule(rule Rule) int {

	if nw.allRules.Get(rule.GetName()) != nil {
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
			if rule.GetConditions()[i].getIdentifiers() == nil || len(rule.GetConditions()[i].getIdentifiers()) == 0 {
				//TODO: condition with no identifiers
				conditionSetNoIdr.Add(rule.GetConditions()[i])
			} else if len(rule.GetConditions()[i].getIdentifiers()) == 1 &&
				!contains(nodeSet, rule.GetConditions()[i].getIdentifiers()[0]) {
				cond := rule.GetConditions()[i]
				nw.createClassFilterNode(rule, nodesOfRule, classNodeLinksOfRule, cond.getIdentifiers()[0], cond, nodeSet)
			} else {
				conditionSet.Add(rule.GetConditions()[i])
			}
		}
	}

	nw.buildNetwork(rule, nodesOfRule, classNodeLinksOfRule, conditionSet, nodeSet, conditionSetNoIdr)

	context := make([]interface{}, 2)
	context[0] = nw
	context[1] = nodesOfRule

	nw.allClassNodes.ForEach(optimizeNetwork, context)
	// nw.optimizeNetwork(nodesOfRule)

	nw.setClassNodeAndLinkJoinTables(nodesOfRule, classNodeLinksOfRule)

	//Add the rule to the network
	nw.allRules.Put(rule.GetName(), rule)

	//Add RuleNodes
	nw.ruleNameNodesOfRule.Put(rule.GetName(), nodesOfRule)

	//Add NodeLinks
	nw.ruleNameClassNodeLinksOfRule.Put(rule.GetName(), classNodeLinksOfRule)
	return 0
}

func (nw *reteNetworkImpl) setClassNodeAndLinkJoinTables(nodesOfRule utils.ArrayList,
	classNodeLinksOfRule utils.ArrayList) {
	//TODO: add join table ids to nodes and links
}

func (nw *reteNetworkImpl) RemoveRule(ruleName string) Rule {

	val := nw.allRules.Remove(ruleName)
	if val == nil {
		//TODO: log a message
		return nil
	}
	rule := val.(Rule)

	classNodeLinksOfRule := nw.ruleNameClassNodeLinksOfRule.Remove(ruleName).(utils.ArrayList)
	if classNodeLinksOfRule != nil {
		classNodeLinksOfRule.ForEach(removeRuleHelper, nil)
	}

	nodesOfRuleItem := nw.ruleNameNodesOfRule.Remove(ruleName)
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

	nw.ruleNameNodesOfRule.Remove(ruleName)
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

func optimizeNetwork(key string, val interface{}, context []interface{}) {
	nodesOfRule := context[1].(utils.ArrayList)
	classNode := val.(classNode)
	for j := 0; j < classNode.getClassNodeLinks().Len(); j++ {
		nodeLink := classNode.getClassNodeLinks().Get(j).(classNodeLink)
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

func contains(nodeSet utils.ArrayList, identifierVar identifier) bool {
	identifiers := []identifier{identifierVar}
	for i := 0; i < nodeSet.Len(); i++ {
		node := nodeSet.Get(i).(node)
		if ContainedByFirst(node.getIdentifiers(), identifiers) {
			return true
		}
	}
	return false
}

func (nw *reteNetworkImpl) buildNetwork(rule Rule, nodesOfRule utils.ArrayList, classNodeLinksOfRule utils.ArrayList,
	conditionSet utils.ArrayList, nodeSet utils.ArrayList, conditionSetNoIdr utils.ArrayList) {
	if conditionSet.Len() == 0 {
		if nodeSet.Len() == 1 {
			node := nodeSet.Get(0).(node)
			if ContainedByFirst(node.getIdentifiers(), rule.GetIdentifiers()) {
				//TODO: Re evaluate set later..

				lastNode := node
				//check conditions with no identifierVar
				for i := 0; i < conditionSetNoIdr.Len(); i++ {
					conditionVar := conditionSetNoIdr.Get(i).(condition)
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
			nw.createClassFilterNode(rule, nodesOfRule, classNodeLinksOfRule, conditionVar.getIdentifiers()[0], nil, nodeSet)
			nw.buildNetwork(rule, nodesOfRule, classNodeLinksOfRule, conditionSet, nodeSet, conditionSetNoIdr)
		}
	}
}

func (nw *reteNetworkImpl) createFilterNode(rule Rule, nodesOfRule utils.ArrayList, conditionSet utils.ArrayList, nodeSet utils.ArrayList) bool {
	for i := 0; i < conditionSet.Len(); i++ {
		conditionVar := conditionSet.Get(i).(condition)
		for i := 0; i < nodeSet.Len(); i++ {
			node := nodeSet.Get(i).(node)
			if ContainedByFirst(node.getIdentifiers(), conditionVar.getIdentifiers()) {
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

func (nw *reteNetworkImpl) createJoinNodeFromExisting(rule Rule, nodesOfRule utils.ArrayList, conditionSet utils.ArrayList, nodeSet utils.ArrayList) bool {
	maxCommonIdr := -1
	numOfIdentifiers := 0
	joinThese := make([]node, 2)
	var targetCondition condition
	for i := 0; i < conditionSet.Len(); i++ {
		conditionVar := conditionSet.Get(i).(condition)
		for j := 0; j < nodeSet.Len(); j++ {
			leftNode := nodeSet.Get(j).(node)
			for k := j + 1; k < nodeSet.Len(); k++ {
				rightNode := nodeSet.Get(k).(node)
				if OtherTwoAreContainedByFirst(conditionVar.getIdentifiers(), leftNode.getIdentifiers(), rightNode.getIdentifiers()) {
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

func (nw *reteNetworkImpl) createJoinNodeFromSome(rule Rule, nodesOfRule utils.ArrayList,
	classNodeLinksOfRule utils.ArrayList, conditionSet utils.ArrayList, nodeSet utils.ArrayList) bool {
	leastNeeded := math.MaxUint32
	maxIdentifier := -1
	var targetNode node
	var targetCondition condition
	for i := 0; i < conditionSet.Len(); i++ {
		conditionVar := conditionSet.Get(i).(condition)
		for j := 0; j < nodeSet.Len(); j++ {
			nodeIdentifiers := nodeSet.Get(j).(node).getIdentifiers()
			need := len(SecondMinusFirst(nodeIdentifiers, conditionVar.getIdentifiers()))
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
	nodeIdentifiers := SecondMinusFirst(targetNode.getIdentifiers(), targetCondition.getIdentifiers())
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

func (nw *reteNetworkImpl) createClassFilterNode(rule Rule, nodesOfRule utils.ArrayList, classNodeLinksOfRule utils.ArrayList, identifierVar identifier, conditionVar condition, nodeSet utils.ArrayList) filterNode {
	identifiers := []identifier{identifierVar}
	classNodeVar := getClassNode(nw, identifierVar.getName())
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

func (nw *reteNetworkImpl) createJoinNode(rule Rule, nodesOfRule utils.ArrayList, leftNode node, rightNode node, joinCondition condition, conditionSet utils.ArrayList, nodeSet utils.ArrayList) {

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

func findBestNode(nodeSet utils.ArrayList, matchIdentifiers []identifier, notThis node) node {
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

func (nw *reteNetworkImpl) findConditionWithLeastIdentifiers(conditionSet utils.ArrayList) condition {
	least := math.MaxUint16
	var leastIdentifiers condition
	for i := 0; i < conditionSet.Len(); i++ {
		c := conditionSet.Get(i).(condition)
		lenIdr := len(c.getIdentifiers())
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

func getClassNode(nw *reteNetworkImpl, name string) classNode {
	var classNodeVar classNode
	val := nw.allClassNodes.Get(name)
	if val == nil {
		classNodeVar = newClassNode(name)
		nw.allClassNodes.Put(name, classNodeVar)
	} else {
		classNodeVar = val.(classNode)
	}
	return classNodeVar
}

func (nw *reteNetworkImpl) String() string {

	str := "\n>>> Class View <<<\n"

	for _, val := range nw.allClassNodes.GetMap() {
		classNodeImpl := val.(classNode)
		str += classNodeImpl.String() + "\n"
	}
	str += ">>>> Rule View <<<<\n"

	for _, val := range nw.allRules.GetMap() {
		rule := val.(Rule)
		str += nw.PrintRule(rule)
	}

	return str
}

func pickIdentifier(idrs []identifier) identifier {
	return idrs[0]
}

func (nw *reteNetworkImpl) PrintRule(rule Rule) string {
	str := "[Rule (" + rule.GetName() + ") Id(" + strconv.Itoa(rule.GetID()) + ")]\n"

	nodesOfRule := nw.ruleNameNodesOfRule.Get(rule.GetName()).(utils.ArrayList)

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
	classNodesLinksOfRule := nw.ruleNameClassNodeLinksOfRule.Get(ruleName).(utils.ArrayList)
	links := ""
	for i := 0; i < classNodesLinksOfRule.Len(); i++ {
		classNodeLinkOfRule := classNodesLinksOfRule.Get(i).(classNodeLink)
		if classNodeLinkOfRule.getIdentifier().getName() == classNodeImpl.name {
			links += "\n\t\t" + classNodeLinkOfRule.String()
		}
	}
	return "\t[ClassNode Class(" + classNodeImpl.getName() + ")" + links + "]\n"
}

func (nw *reteNetworkImpl) Assert(tuple model.StreamTuple) {
	cr := newConflictRes()
	dataSource := tuple.GetStreamDataSource()
	listItem := nw.allClassNodes.Get(string(dataSource))
	if listItem != nil {
		classNodeVar := listItem.(classNode)
		classNodeVar.assert(tuple, cr)
	} else {
		fmt.Println("No rule exists for data stream: " + dataSource)
	}

	cr.resolveConflict()

}

func (nw *reteNetworkImpl) Retract(tuple model.StreamTuple) {
	reteHandle := allHandles[tuple]
	if reteHandle == nil {
		//TODO: Nothing to retract!
		return
	}
	reteHandle.removeJoinTableRowRefs()
}
