package rete

import (
	"context"
	"fmt"
	"math"
	"strconv"

	"github.com/project-flogo/rules/common/model"

	"container/list"
	"github.com/project-flogo/rules/rete/internal/types"
	"sync"
	"time"
	//"github.com/project-flogo/rules/rete/common"
	"github.com/project-flogo/rules/rete/common"
)

type reteNetworkImpl struct {
	//unique name of the network. used for namespacing in storage, etc
	prefix string

	//All rules in the network
	allRules map[string]model.Rule //(Rule)

	//Holds the DataSource name as key, and ClassNodes as value
	allClassNodes map[string]classNode //ClassNode in network

	//Holds the Rule name as key and pointer to a slice of RuleNodes as value
	ruleNameNodesOfRule map[string]*list.List //*list.List of Nodes of rule

	//Holds the Rule name as key and a pointer to a slice of NodeLinks as value
	ruleNameClassNodeLinksOfRule map[string]*list.List //*list.List of ClassNodeLink

	//handleService map[string]types.ReteHandle
	handleService types.HandleService

	crudLock   sync.RWMutex
	txnHandler []model.RtcTransactionHandler
	txnContext []interface{}

	//jtService map[int]types.JoinTable
	jtService types.JtService

	jtRefsService types.JtRefsService

	config map[string]string

	factory       *TypeFactory
	idGen         types.IdGen
	tupleStore    model.TupleStore
	joinNodeNames int
}

//NewReteNetwork ... creates a new rete network
func NewReteNetwork(sessionName string, jsonConfig string) types.Network {
	reteNetworkImpl := reteNetworkImpl{}
	reteNetworkImpl.initReteNetwork(sessionName, jsonConfig)
	return &reteNetworkImpl
}

func (nw *reteNetworkImpl) initReteNetwork(sessionName string, config string) error {
	//nw.currentId = 0
	nw.allRules = make(map[string]model.Rule)
	nw.allClassNodes = make(map[string]classNode)
	nw.ruleNameNodesOfRule = make(map[string]*list.List)
	nw.ruleNameClassNodeLinksOfRule = make(map[string]*list.List)
	nw.txnHandler = []model.RtcTransactionHandler{}

	factory, err := NewFactory(nw, config)
	if err != nil {
		return err
	}
	//nw.factory = factory

	//if factory.parsedJson != nil {
	//	reteCfg := factory.parsedJson["rs"].(map[string]interface{})
	//	nw.prefix = reteCfg["prefix"].(string)
	//}
	nw.prefix = sessionName
	nw.idGen = factory.getIdGen()
	nw.jtService = factory.getJoinTableCollection()
	nw.handleService = factory.getHandleCollection()
	nw.jtRefsService = factory.getJoinTableRefs()
	nw.initNwServices()
	return nil

}

func (nw *reteNetworkImpl) initNwServices() {
	nw.idGen.Init()
	nw.jtService.Init()
	nw.handleService.Init()
	nw.jtRefsService.Init()
}

func (nw *reteNetworkImpl) AddRule(rule model.Rule) (err error) {

	nw.crudLock.Lock()
	defer nw.crudLock.Unlock()

	if nw.allRules[rule.GetName()] != nil {
		return fmt.Errorf("Rule already exists.." + rule.GetName())
	}
	conditionSet := list.New()
	conditionSetNoIdr := list.New()
	nodeSet := list.New()

	nodesOfRule := list.New()
	classNodeLinksOfRule := list.New()

	conditions := rule.GetConditions()
	noIdrConditionCnt := 0
	if len(conditions) == 0 {
		identifierVar := pickIdentifier(rule.GetIdentifiers())
		nw.createClassFilterNode(rule, nodesOfRule, classNodeLinksOfRule, identifierVar, nil, nodeSet)
	} else {
		for i := 0; i < len(conditions); i++ {
			if conditions[i].GetIdentifiers() == nil || len(conditions[i].GetIdentifiers()) == 0 {
				conditionSetNoIdr.PushBack(conditions[i])
				noIdrConditionCnt++
			} else if len(conditions[i].GetIdentifiers()) == 1 &&
				!contains(nodeSet, conditions[i].GetIdentifiers()[0]) {
				cond := conditions[i]
				nw.createClassFilterNode(rule, nodesOfRule, classNodeLinksOfRule, cond.GetIdentifiers()[0], cond, nodeSet)
			} else {
				conditionSet.PushBack(conditions[i])
			}
		}
	}
	if len(rule.GetConditions()) != 0 && noIdrConditionCnt == len(rule.GetConditions()) {
		idr := pickIdentifier(rule.GetIdentifiers())
		nw.createClassFilterNode(rule, nodesOfRule, classNodeLinksOfRule, idr, nil, nodeSet)
	}
	nw.buildNetwork(rule, nodesOfRule, classNodeLinksOfRule, conditionSet, nodeSet, conditionSetNoIdr)

	cntxt := make([]interface{}, 2)
	cntxt[0] = nw
	cntxt[1] = nodesOfRule
	for _, classNode := range nw.allClassNodes {
		optimizeNetwork(classNode, cntxt)
	}
	// nw.optimizeNetwork(nodesOfRule)

	nw.setClassNodeAndLinkJoinTables(nodesOfRule, classNodeLinksOfRule)

	//Add the rule to the network
	nw.allRules[rule.GetName()] = rule

	//Add RuleNodes
	nw.ruleNameNodesOfRule[rule.GetName()] = nodesOfRule

	//Add NodeLinks
	nw.ruleNameClassNodeLinksOfRule[rule.GetName()] = classNodeLinksOfRule
	return nil
}

func (nw *reteNetworkImpl) setClassNodeAndLinkJoinTables(nodesOfRule *list.List,
	classNodeLinksOfRule *list.List) {
}

func (nw *reteNetworkImpl) RemoveRule(ruleName string) model.Rule {

	nw.crudLock.Lock()
	defer nw.crudLock.Unlock()

	rule := nw.allRules[ruleName]
	delete(nw.allRules, ruleName)
	if rule == nil {
		//TODO: log a message
		return nil
	}

	classNodeLinksOfRule := nw.ruleNameClassNodeLinksOfRule[ruleName]
	delete(nw.ruleNameClassNodeLinksOfRule, ruleName)
	if classNodeLinksOfRule != nil {
		for e := classNodeLinksOfRule.Front(); e != nil; e = e.Next() {
			removeRuleHelper(e.Value.(classNodeLink))
		}
	}

	nodesOfRuleItem := nw.ruleNameNodesOfRule[ruleName]
	delete(nw.ruleNameNodesOfRule, ruleName)
	if nodesOfRuleItem != nil {
		for e := nodesOfRuleItem.Front(); e != nil; e = e.Next() {
			n := e.Value.(abstractNode)
			switch nodeImpl := n.(type) {
			//Only interested in joinnodes
			//case *filterNodeImpl:
			//case *classNodeImpl:
			//case *ruleNodeImpl:
			case *joinNodeImpl:
				//nw.removeRefsFromReteHandles(nodeImpl.leftTable)
				//nw.removeRefsFromReteHandles(nodeImpl.rightTable)
				nodeImpl.leftTable.RemoveAllRows()
				nodeImpl.rightTable.RemoveAllRows()
			}
		}
	}

	return rule
}

func (nw *reteNetworkImpl) GetRules() []model.Rule {
	nw.crudLock.RLock()
	defer nw.crudLock.RUnlock()

	rules := make([]model.Rule, 0)

	for _, rule := range nw.allRules {
		rules = append(rules, rule)
	}
	return rules
}

func (nw *reteNetworkImpl) removeRefsFromReteHandles(joinTableVar types.JoinTable) {
	if joinTableVar == nil {
		return
	}
	rIterator := joinTableVar.GetRowIterator()
	for rIterator.HasNext() {
		tableRow := rIterator.Next()
		for _, handle := range tableRow.GetHandles() {
			nw.removeJoinTableRowRefs(handle, nil)
		}
	}
}

func removeRuleHelper(classNodeLinkOfRule classNodeLink) {
	classNodeLinkOfRule.getClassNode().removeClassNodeLink(classNodeLinkOfRule)
}

func optimizeNetwork(classNodeVar classNode, context []interface{}) {
	nodesOfRule := context[1].(*list.List)
	for e := classNodeVar.getClassNodeLinks().Front(); e != nil; e = e.Next() {
		nodeLink := e.Value.(classNodeLink)
		childNode := nodeLink.getChild()
		switch nodeImpl := childNode.(type) {
		case *filterNodeImpl:
			if nodeImpl.conditionVar == nil {
				nodeLink.setChild(nodeImpl.nodeLinkVar.getChild())
				nodeLink.setIsRightChild(nodeImpl.nodeLinkVar.isRightNode())
				removeFromList(nodesOfRule, nodeImpl)
			}
		}
	}
}

func removeFromList(listVar *list.List, val interface{}) {
	for e := listVar.Front(); e != nil; e = e.Next() {
		if e.Value == val {
			listVar.Remove(e)
			break
		}
	}
}

func contains(nodeSet *list.List, identifierVar model.TupleType) bool {
	identifiers := []model.TupleType{identifierVar}
	for e := nodeSet.Front(); e != nil; e = e.Next() {
		n := e.Value.(node)
		if ContainedByFirst(n.getIdentifiers(), identifiers) {
			return true
		}
	}
	return false
}

func (nw *reteNetworkImpl) buildNetwork(rule model.Rule, nodesOfRule *list.List, classNodeLinksOfRule *list.List,
	conditionSet *list.List, nodeSet *list.List, conditionSetNoIdr *list.List) {
	if conditionSet.Len() == 0 {
		if nodeSet.Len() == 1 {
			n := nodeSet.Front().Value.(node)
			if ContainedByFirst(n.getIdentifiers(), rule.GetIdentifiers()) {
				//TODO: Re evaluate set later..

				lastNode := n
				//check conditions with no identifierVar
				for e := conditionSetNoIdr.Front(); e != nil; e = e.Next() {
					conditionVar := e.Value.(model.Condition)
					fNode := newFilterNode(nw, rule, n.getIdentifiers(), conditionVar)
					nodesOfRule.PushBack(fNode)
					newNodeLink(nw, lastNode, fNode, false)
					lastNode = fNode
				}
				//Yoohoo! We have a Rule!!
				ruleNode := newRuleNode(rule)
				newNodeLink(nw, lastNode, ruleNode, false)
				nodesOfRule.PushBack(ruleNode)
			} else {
				idrs := SecondMinusFirst(n.getIdentifiers(), rule.GetIdentifiers())
				fNode := nw.createClassFilterNode(rule, nodesOfRule, classNodeLinksOfRule, idrs[0], nil, nodeSet)
				nw.createJoinNode(rule, nodesOfRule, n, fNode, nil, conditionSet, nodeSet)
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

func (nw *reteNetworkImpl) createFilterNode(rule model.Rule, nodesOfRule *list.List, conditionSet *list.List, nodeSet *list.List) bool {
	for e := conditionSet.Front(); e != nil; e = e.Next() {
		conditionVar := e.Value.(model.Condition)
		for f := nodeSet.Front(); f != nil; f = f.Next() {
			n := f.Value.(node)
			if ContainedByFirst(n.getIdentifiers(), conditionVar.GetIdentifiers()) {
				//TODO
				filterNode := newFilterNode(nw, rule, conditionVar.GetIdentifiers(), conditionVar)
				newNodeLink(nw, n, filterNode, false)
				removeFromList(nodeSet, n)
				nodeSet.PushBack(filterNode)
				nodesOfRule.PushBack(filterNode)
				conditionSet.Remove(e)
				return true
			}
		}
	}

	return false
}

func (nw *reteNetworkImpl) createJoinNodeFromExisting(rule model.Rule, nodesOfRule *list.List, conditionSet *list.List, nodeSet *list.List) bool {
	maxCommonIdr := -1
	numOfIdentifiers := 0
	joinThese := make([]node, 2)
	var targetCondition model.Condition
	for e := conditionSet.Front(); e != nil; e = e.Next() {
		conditionVar := e.Value.(model.Condition)
		for j := nodeSet.Front(); j != nil; j = j.Next() {
			leftNode := j.Value.(node)
			for k := j.Next(); k != nil; k = k.Next() {
				rightNode := k.Value.(node)
				if UnionOfOtherTwoContainsAllFromFirst(conditionVar.GetIdentifiers(), leftNode.getIdentifiers(), rightNode.getIdentifiers()) {
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

func (nw *reteNetworkImpl) createJoinNodeFromSome(rule model.Rule, nodesOfRule *list.List,
	classNodeLinksOfRule *list.List, conditionSet *list.List, nodeSet *list.List) bool {
	leastNeeded := math.MaxInt32
	maxIdentifier := -1
	var targetNode node
	var targetCondition model.Condition
	for e := conditionSet.Front(); e != nil; e = e.Next() {
		conditionVar := e.Value.(model.Condition)
		for j := nodeSet.Front(); j != nil; j = j.Next() {
			nodeIdentifiers := j.Value.(node).getIdentifiers()
			need := len(SecondMinusFirst(nodeIdentifiers, conditionVar.GetIdentifiers()))
			if need < leastNeeded {
				leastNeeded = need
				maxIdentifier = len(nodeIdentifiers)
				targetNode = j.Value.(node)
				targetCondition = conditionVar
			} else if need == leastNeeded {
				if len(nodeIdentifiers) > maxIdentifier {
					maxIdentifier = len(nodeIdentifiers)
					targetNode = j.Value.(node)
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

func (nw *reteNetworkImpl) createClassFilterNode(rule model.Rule, nodesOfRule *list.List, classNodeLinksOfRule *list.List, identifierVar model.TupleType, conditionVar model.Condition, nodeSet *list.List) filterNode {
	identifiers := []model.TupleType{identifierVar}
	classNodeVar := getClassNode(nw, identifierVar)
	filterNodeVar := newFilterNode(nw, rule, identifiers, conditionVar)
	classNodeLink := newClassNodeLink(nw, classNodeVar, filterNodeVar, rule, identifierVar)
	classNodeVar.addClassNodeLink(classNodeLink)
	nodesOfRule.PushBack(classNodeVar)
	nodesOfRule.PushBack(filterNodeVar)
	//TODO: Add to RuleLinks
	classNodeLinksOfRule.PushBack(classNodeLink)
	nodeSet.PushBack(filterNodeVar)
	return filterNodeVar
}

func (nw *reteNetworkImpl) createJoinNode(rule model.Rule, nodesOfRule *list.List, leftNode node, rightNode node, joinCondition model.Condition, conditionSet *list.List, nodeSet *list.List) {
	//TODO handle equivJoins later..

	joinNode := newJoinNode(nw, rule, leftNode.getIdentifiers(), rightNode.getIdentifiers(), joinCondition)
	newNodeLink(nw, leftNode, joinNode, false)
	newNodeLink(nw, rightNode, joinNode, true)
	removeFromList(nodeSet, leftNode)
	removeFromList(nodeSet, rightNode)
	nodeSet.PushBack(joinNode)
	nodesOfRule.PushBack(joinNode)
	if joinCondition != nil {
		removeFromList(conditionSet, joinCondition)
	}
}

func findBestNode(nodeSet *list.List, matchIdentifiers []model.TupleType, notThis node) node {
	var foundNode node
	foundNode = nil
	foundIdr := 0

	for e := nodeSet.Front(); e != nil; e = e.Next() {
		n := e.Value.(node)
		if n == notThis {
			continue
		}
		foundMatch := len(IntersectionIdentifiers(n.getIdentifiers(), matchIdentifiers))
		if foundMatch > foundIdr {
			foundIdr = foundMatch
			foundNode = n
		}
	}
	return foundNode
}

func (nw *reteNetworkImpl) findConditionWithLeastIdentifiers(conditionSet *list.List) model.Condition {
	least := math.MaxUint16
	var leastIdentifiers model.Condition
	for e := conditionSet.Front(); e != nil; e = e.Next() {
		c := e.Value.(model.Condition)
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

func getClassNode(nw *reteNetworkImpl, name model.TupleType) classNode {
	classNodeVar := nw.allClassNodes[string(name)]
	if classNodeVar == nil {
		classNodeVar = newClassNode(string(name))
		nw.allClassNodes[string(name)] = classNodeVar
	}
	return classNodeVar
}

func (nw *reteNetworkImpl) String() string {
	nw.crudLock.RLock()
	defer nw.crudLock.RUnlock()

	str := "\n>>> Class View <<<\n"
	for _, classNodeImpl := range nw.allClassNodes {
		str += classNodeImpl.String() + "\n"
	}

	str += ">>>> Rule View <<<<\n"
	for _, rule := range nw.allRules {
		str += "[Rule (" + rule.GetName() + ") Id()]\n"
		nodesOfRule := nw.ruleNameNodesOfRule[rule.GetName()]
		for e := nodesOfRule.Front(); e != nil; e = e.Next() {
			n := e.Value.(abstractNode)
			switch nodeImpl := n.(type) {
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
	}

	return str
}

func pickIdentifier(idrs []model.TupleType) model.TupleType {
	return idrs[0]
}

func (nw *reteNetworkImpl) PrintRule(rule model.Rule) string {
	nw.crudLock.RLock()
	defer nw.crudLock.RUnlock()

	//str := "[Rule (" + rule.GetName() + ") Id(" + strconv.Itoa(rule.GetID()) + ")]\n"
	str := "[Rule (" + rule.GetName() + ") Id()]\n"

	nodesOfRule := nw.ruleNameNodesOfRule[rule.GetName()]

	for e := nodesOfRule.Front(); e != nil; e = e.Next() {
		n := e.Value.(abstractNode)
		switch nodeImpl := n.(type) {
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
	classNodesLinksOfRule := nw.ruleNameClassNodeLinksOfRule[ruleName]
	links := ""
	for e := classNodesLinksOfRule.Front(); e != nil; e = e.Next() {
		classNodeLinkOfRule := e.Value.(classNodeLink)
		if string(classNodeLinkOfRule.GetIdentifier()) == classNodeImpl.name {
			links += "\n\t\t" + classNodeLinkOfRule.String()
		}
	}
	return "\t[ClassNode Class(" + classNodeImpl.getName() + ")" + links + "]\n"
}

func (nw *reteNetworkImpl) Assert(ctx context.Context, rs model.RuleSession, tuple model.Tuple, changedProps map[string]bool, mode common.RtcOprn) error {

	if ctx == nil {
		ctx = context.Background()
	}

	reteCtxVar, isRecursive, newCtx := getOrSetReteCtx(ctx, nw, rs)

	if !isRecursive {
		nw.crudLock.Lock()
		defer nw.crudLock.Unlock()

		err := nw.assertInternal(newCtx, tuple, changedProps, mode)
		if err != nil {
			return err
		}

		reteCtxVar.getConflictResolver().resolveConflict(newCtx)
		//if Timeout is 0, remove it from rete
		td := model.GetTupleDescriptor(tuple.GetTupleType())
		if td != nil {
			if td.TTLInSeconds == 0 { //remove immediately.
				nw.removeTupleFromRete(tuple)
			} else if td.TTLInSeconds > 0 { // TTL for the tuple type, after that, remove it from RETE
				time.AfterFunc(time.Second*time.Duration(td.TTLInSeconds), func() {
					nw.crudLock.Lock()
					defer nw.crudLock.Unlock()
					nw.removeTupleFromRete(tuple)
				})
			} //else, its -ve and means, never expire
		}
		if nw.txnHandler != nil {
			rtcTxn := newRtcTxn(reteCtxVar.getRtcAdded(), reteCtxVar.getRtcModified(), reteCtxVar.getRtcDeleted())
			for i, txnHandler := range nw.txnHandler {
				txnHandler(newCtx, rs, rtcTxn, nw.txnContext[i])
			}
		}
		return nil
	}

	reteCtxVar.getOpsList().PushBack(newAssertEntry(tuple, changedProps, mode))
	return nil
}

func (nw *reteNetworkImpl) removeTupleFromRete(tuple model.Tuple) {
	reteHandle := nw.handleService.RemoveHandle(tuple)
	if reteHandle != nil {
		nw.removeJoinTableRowRefs(reteHandle, nil)
	}
}

func (nw *reteNetworkImpl) Retract(ctx context.Context, rs model.RuleSession, tuple model.Tuple, changedProps map[string]bool, mode common.RtcOprn) error {

	if ctx == nil {
		ctx = context.Background()
	}
	reteCtxVar, isRecursive, ctx := getOrSetReteCtx(ctx, nw, rs)
	if !isRecursive {
		nw.crudLock.Lock()
		defer nw.crudLock.Unlock()
		err := nw.retractInternal(ctx, tuple, changedProps, mode)
		if err != nil {
			return err
		}
		if nw.txnHandler != nil && mode == common.DELETE {
			rtcTxn := newRtcTxn(reteCtxVar.getRtcAdded(), reteCtxVar.getRtcModified(), reteCtxVar.getRtcDeleted())
			for i, txnHandler := range nw.txnHandler {
				txnHandler(ctx, rs, rtcTxn, nw.txnContext[i])
			}
		}
	} else {
		reteCtxVar.getOpsList().PushBack(newDeleteEntry(tuple, mode, changedProps))
	}

	return nil
}

func (nw *reteNetworkImpl) retractInternal(ctx context.Context, tuple model.Tuple, changedProps map[string]bool, mode common.RtcOprn) error {
	handle := nw.handleService.GetHandle(tuple)
	if handle == nil {
		return fmt.Errorf("Tuple with key [%s] doesn't exist", tuple.GetKey().String())
	} else if handle.GetStatus() != "created" {
		return fmt.Errorf("Tuple with key [%s] is being asserted", tuple.GetKey().String())
	}

	if ctx == nil {
		ctx = context.Background()
	}
	rCtx, _, _ := getOrSetReteCtx(ctx, nw, nil)

	nw.removeJoinTableRowRefs(handle, changedProps)
	// add it to the delete list
	if mode == common.DELETE {
		rCtx.addToRtcDeleted(tuple)
	}
	nw.handleService.RemoveHandle(tuple)

	return nil
}

func (nw *reteNetworkImpl) GetAssertedTuple(key model.TupleKey) model.Tuple {
	reteHandle := nw.handleService.GetHandleByKey(key)
	if reteHandle != nil {
		return reteHandle.GetTuple()
	}
	return nil
}

func (nw *reteNetworkImpl) assertInternal(ctx context.Context, tuple model.Tuple, changedProps map[string]bool, mode common.RtcOprn) error {
	if mode == common.ADD {
		handle, exists := nw.getOrCreateHandle(ctx, tuple)
		if exists {
			return fmt.Errorf("Tuple with key [%s] already asserted", tuple.GetKey().String())
		}
		defer handle.SetStatus("created")
	}

	tupleType := tuple.GetTupleType()
	listItem := nw.allClassNodes[string(tupleType)]
	if listItem != nil {
		classNodeVar := listItem.(classNode)
		classNodeVar.assert(ctx, tuple, changedProps)
	}
	td := model.GetTupleDescriptor(tuple.GetTupleType())
	if td != nil {
		if td.TTLInSeconds != 0 && mode == common.ADD {
			rCtx := getReteCtx(ctx)
			if rCtx != nil {
				rCtx.addToRtcAdded(tuple)
			}
		}
	}

	return nil
}

func (nw *reteNetworkImpl) getOrCreateHandle(ctx context.Context, tuple model.Tuple) (types.ReteHandle, bool) {
	return nw.handleService.GetOrCreateHandle(nw, tuple)
}

func (nw *reteNetworkImpl) getHandle(tuple model.Tuple) types.ReteHandle {
	h := nw.handleService.GetHandleByKey(tuple.GetKey())
	return h
}

func (nw *reteNetworkImpl) RegisterRtcTransactionHandler(txnHandler model.RtcTransactionHandler, txnContext interface{}) {
	nw.txnHandler = append(nw.txnHandler, txnHandler)
	nw.txnContext = append(nw.txnContext, txnContext)
}

func (nw *reteNetworkImpl) GetConfigValue(key string) string {
	val, _ := nw.config[key]
	return val
}

func (nw *reteNetworkImpl) GetConfig() map[string]string {
	return nw.config
}

func (nw *reteNetworkImpl) getFactory() *TypeFactory {
	return nw.factory
}

func (nw *reteNetworkImpl) SetTupleStore(tupleStore model.TupleStore) {
	nw.tupleStore = tupleStore
}
func (nw *reteNetworkImpl) GetTupleStore() model.TupleStore {
	return nw.tupleStore
}

func getOrCreateHandle(ctx context.Context, tuple model.Tuple) (types.ReteHandle, bool) {
	reteCtxVar := getReteCtx(ctx)
	return reteCtxVar.getNetwork().getOrCreateHandle(ctx, tuple)
}

func (nw *reteNetworkImpl) removeJoinTableRowRefs(hdl types.ReteHandle, changedProps map[string]bool) {
	tuple := hdl.GetTuple()
	alias := tuple.GetTupleType()

	hdlTblIter := nw.jtRefsService.GetTableIterator(hdl)

	for hdlTblIter.HasNext() {
		joinTable := hdlTblIter.Next()
		if joinTable == nil {
			continue
		}
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

		nw.removeJtRowsForTable(hdl, joinTable)
	}
	//TODO: Remove the current entry from underneath
	hdlTblIter.Remove()
}

func (nw *reteNetworkImpl) removeJtRowsForTable(hdl types.ReteHandle, joinTable types.JoinTable) {
	rowIter := nw.jtRefsService.GetRowIterator(hdl, joinTable.GetName())
	////Remove rows from corresponding join tables
	for rowIter.HasNext() {
		row := rowIter.Next()
		//first, from jTable, remove row
		joinTable.RemoveRow(row.GetID())
		for _, otherHdl := range row.GetHandles() {
			//if otherHdl.GetTupleKey().String() != hdl.GetTupleKey().String() {
			nw.jtRefsService.RemoveRowEntry(otherHdl, joinTable.GetName(), row.GetID())
			//}
		}
		//TODO: Delete the rowRef itself
		rowIter.Remove()
	}
}

func (nw *reteNetworkImpl) getJoinNodeName() string {
	name := strconv.Itoa(nw.joinNodeNames)
	nw.joinNodeNames++
	return name
}

func (nw *reteNetworkImpl) GetIdGenService() types.IdGen {
	return nw.idGen
}

func (nw *reteNetworkImpl) GetJtService() types.JtService {
	return nw.jtService
}

func (nw *reteNetworkImpl) GetJtRefService() types.JtRefsService {
	return nw.jtRefsService
}

func (nw *reteNetworkImpl) GetHandleService() types.HandleService {
	return nw.handleService
}

func (nw *reteNetworkImpl) GetPrefix() string {
	return nw.prefix
}
