package rete

import (
	"context"
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

func (rn *ruleNodeImpl) String() string {
	return "\t[RuleNode id(" + strconv.Itoa(rn.id) + "): \n" +
		"\t\tIdentifier           = " + IdentifiersToString(rn.identifiers) + " ;\n" +
		"\t\tRule                 = " + rn.rule.GetName() + "]\n"
}

func (rn *ruleNodeImpl) assertObjects(ctx context.Context, handles []reteHandle, isRight bool) {
	// fmt.Println("Rule " + rn.getRule().GetName() + " fired, total tuples:" + strconv.Itoa(len(handles)))
	// tuples := copyIntoTupleArray(handles)
	// rn.getRule.performAction(tuples)
	tupleMap := copyIntoTupleMap(handles)

	cr := getReteCtx(ctx).getConflictResolver()

	cr.addAgendaItem(rn.getRule(), tupleMap)

	// actionFn := rn.getRule().GetActionFn()
	// if actionFn != nil {
	// 	actionFn(rn.getRule().GetName(), tupleMap)
	// }
}
func (rn *ruleNodeImpl) getRule() Rule {
	return rn.rule
}
