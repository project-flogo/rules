package rete

import (
	"strconv"
)

//ruleNode the leaf node of the rule network for a Rule
type ruleNode interface {
	node
	getRule() Rule
}

type ruleNodeImpl struct {
	nodeImpl
	rule Rule
}

func newRuleNode(rule Rule) ruleNode {
	rn := ruleNodeImpl{}
	rn.identifiers = rule.GetIdentifiers()
	rn.rule = rule
	return &rn
}

func (ruleNodeImplVar *ruleNodeImpl) String() string {
	return "\t[RuleNode id(" + strconv.Itoa(ruleNodeImplVar.id) + "): \n" +
		"\t\tIdentifier           = " + IdentifiersToString(ruleNodeImplVar.identifiers) + " ;\n" +
		"\t\tRule                 = " + ruleNodeImplVar.rule.GetName() + "]\n"
}

func (ruleNodeImplVar *ruleNodeImpl) assertObjects(handles []reteHandle, isRight bool) {
	// fmt.Println("Rule " + ruleNodeImplVar.getRule().GetName() + " fired, total tuples:" + strconv.Itoa(len(handles)))
	// tuples := copyIntoTupleArray(handles)
	// ruleNodeImplVar.getRule.performAction(tuples)
	tupleMap := copyIntoTupleMap(handles)
	actionFn := ruleNodeImplVar.getRule().GetActionFn()
	if actionFn != nil {
		actionFn(ruleNodeImplVar.getRule().GetName(), tupleMap)
	}

}
func (ruleNodeImplVar *ruleNodeImpl) getRule() Rule {
	return ruleNodeImplVar.rule
}
