package ruleapi

import (
	"github.com/TIBCOSoftware/bego/common/model"
	"strings"
)

type ruleImpl struct {
	name        string
	identifiers []model.TupleType
	conditions  []model.Condition
	actionFn    model.ActionFunction
	priority    int
	deps        map[model.TupleType]map[string]bool
}

//NewRule ... Create a new rule
func NewRule(name string) model.MutableRule {
	rule := ruleImpl{}
	rule.initRuleImpl(name)
	return &rule
}

func (rule *ruleImpl) initRuleImpl(name string) {
	rule.name = name
	rule.identifiers = []model.TupleType{}
	rule.conditions = []model.Condition{}
	rule.deps = make(map[model.TupleType]map[string]bool)
}

func (rule *ruleImpl) GetName() string {
	return rule.name
}

func (rule *ruleImpl) GetActionFn() model.ActionFunction {
	return rule.actionFn
}

func (rule *ruleImpl) GetConditions() []model.Condition {
	return rule.conditions
}

func (rule *ruleImpl) SetActionFn(actionFn model.ActionFunction) {
	rule.actionFn = actionFn
}

func (rule *ruleImpl) AddCondition(conditionName string, idrs []model.TupleType, cfn model.ConditionEvaluator, ctx model.ConditionContext) {
	rule.addCond(conditionName, idrs, cfn, ctx, false)
}

func (rule *ruleImpl) addCond(conditionName string, idrs []model.TupleType, cfn model.ConditionEvaluator, ctx model.ConditionContext, setIdr bool) {
	condition := newCondition(conditionName, rule, idrs, cfn, ctx)
	rule.conditions = append(rule.conditions, condition)

	for _, cidr := range idrs {
		if len(rule.identifiers) == 0 {
			rule.identifiers = append(rule.identifiers, cidr)
		} else {
			for _, ridr := range rule.identifiers {
				if cidr != ridr {
					rule.identifiers = append(rule.identifiers, cidr)
					break
				}
			}
		}
	}

}

func (rule *ruleImpl) GetPriority() int {
	return rule.priority
}

func (rule *ruleImpl) SetPriority(priority int) {
	rule.priority = priority
}

func (rule *ruleImpl) String() string {
	str := ""
	str += "[Rule: (" + ") " + rule.name + "\n"
	//str += "[Rule: (" + strconv.Itoa(rule.id) + ") " + rule.name + "\n"

	str += "\t[Conditions:\n"
	for _, cond := range rule.conditions {
		str += "\t\t" + cond.String() + "\n"
	}
	// idrs := ""
	// for i := 0; i < len(rule.identifiers); i++ {
	// 	idrs += rule.identifiers[i].String() + ", "
	// }
	str += "\t[Idrs:" + model.IdentifiersToString(rule.identifiers) + "]\n"
	return str
	// return str + idrs + "]\n"
}

func (rule *ruleImpl) GetIdentifiers() []model.TupleType {
	return rule.identifiers
}

func (rule *ruleImpl) SetAction(actionFn model.ActionFunction) {
	rule.actionFn = actionFn
}

////IdentifiersToString Take a slice of Identifiers and return a string representation
//func IdentifiersToString(identifiers []model.TupleType) string {
//	str := ""
//	for _, idr := range identifiers {
//		str += string(idr) + ", "
//	}
//	return str
//}

func (rule *ruleImpl) AddConditionWithDependency(conditionName string, idrs []string, cFn model.ConditionEvaluator, ctx model.ConditionContext) {
	typeDepMap := map[model.TupleType]bool{}
	//cwd := model.ConditionAndDep{"n1", []string{"p1", "p2", "p3"}}
	for _, idr := range idrs {
		aliasProp := strings.Split(idr, ".")

		alias := model.TupleType(aliasProp[0])
		typeDepMap[alias] = true
		prop := aliasProp[1]

		propMap, found := rule.deps[alias]
		if !found {
			propMap = map[string]bool{}
			rule.deps[alias] = propMap
		}
		propMap[prop] = true
	}
	typeDeps := []model.TupleType{}

	for key, _ := range typeDepMap {
		typeDeps = append(typeDeps, key)
	}

	rule.addCond(conditionName, typeDeps, cFn, ctx, true)

}

func (rule *ruleImpl) GetDeps() map[model.TupleType]map[string]bool {
	return rule.deps
}
