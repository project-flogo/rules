package rete

import "github.com/TIBCOSoftware/bego/common/model"

//condition ... is a rete condtion
type condition interface {
	getName() string
	getIdentifiers() []identifier
	// Eval([]model.StreamTuple) bool

	//Stringer.String interface
	String() string
	getEvaluator() model.ConditionEvaluator
	getRule() Rule
}

type conditionImpl struct {
	name        string
	rule        Rule
	identifiers []identifier
	cfn         model.ConditionEvaluator
}

func newCondition(name string, rule Rule, identifiers []string, cfn model.ConditionEvaluator) condition {
	c := conditionImpl{}
	c.initConditionImpl(name, rule, identifiers, cfn)
	return &c
}

func (cnd *conditionImpl) initConditionImpl(name string, rule Rule, identifiers []string, cfn model.ConditionEvaluator) {
	cnd.name = name
	cnd.rule = rule
	for i := 0; i < len(identifiers); i++ {
		idName := identifiers[i]
		idr := newIdentifier(idName)
		cnd.identifiers = append(cnd.identifiers, idr)
	}
	cnd.cfn = cfn
}

func (cnd *conditionImpl) getIdentifiers() []identifier {
	return cnd.identifiers
}

func (cnd *conditionImpl) getEvaluator() model.ConditionEvaluator {
	return cnd.cfn
}

func (cnd *conditionImpl) String() string {
	return "[Condition: name:" + cnd.name + ", idrs:" + IdentifiersToString(cnd.identifiers) + "]"
}

func (cnd *conditionImpl) getName() string {
	return cnd.name
}

func (cnd *conditionImpl) getRule() Rule {
	return cnd.rule
}
