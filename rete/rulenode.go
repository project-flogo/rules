package rete

import (
	"context"
	"strconv"

	"github.com/project-flogo/rules/common/model"
)

//ruleNode the leaf node of the rule network for a Rule
type ruleNode interface {
	node
	getRule() model.Rule
}

type ruleNodeImpl struct {
	nodeImpl
	rule model.Rule
}

func newRuleNode(nw Network, rule model.Rule) ruleNode {
	rn := ruleNodeImpl{}
	rn.nodeImpl.initNodeImpl(nw, rule, rule.GetIdentifiers())
	rn.identifiers = rule.GetIdentifiers()
	rn.rule = rule
	return &rn
}

func (rn *ruleNodeImpl) String() string {
	return "\t[RuleNode id(" + strconv.Itoa(rn.id) + "): \n" +
		"\t\tIdentifier           = " + model.IdentifiersToString(rn.identifiers) + " ;\n" +
		"\t\tRule                 = " + rn.rule.GetName() + "]\n"
}

func (rn *ruleNodeImpl) assertObjects(ctx context.Context, handles []reteHandle, isRight bool) {

	tupleMap := copyIntoTupleMap(handles)

	cr := getReteCtx(ctx).getConflictResolver()

	cr.addAgendaItem(rn.getRule(), tupleMap)

}

func (rn *ruleNodeImpl) getRule() model.Rule {
	return rn.rule
}
