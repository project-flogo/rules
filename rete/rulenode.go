package rete

import (
	"context"
	"strconv"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/internal/types"
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

func newRuleNode(rule model.Rule) ruleNode {
	rn := ruleNodeImpl{}
	rn.identifiers = rule.GetIdentifiers()
	rn.rule = rule
	return &rn
}

func (rn *ruleNodeImpl) String() string {
	return "\t[RuleNode id(" + strconv.Itoa(rn.GetID()) + "): \n" +
		"\t\tIdentifier           = " + model.IdentifiersToString(rn.identifiers) + " ;\n" +
		"\t\tRule                 = " + rn.rule.GetName() + "]\n"
}

func (rn *ruleNodeImpl) assertObjects(ctx context.Context, handles []types.ReteHandle, isRight bool) {

	tupleMap := copyIntoTupleMap(handles)

	cr := getReteCtx(ctx).getConflictResolver()

	cr.addAgendaItem(rn.getRule(), tupleMap)

}

func (rn *ruleNodeImpl) getRule() model.Rule {
	return rn.rule
}
