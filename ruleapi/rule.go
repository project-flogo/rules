package ruleapi

import (
	"github.com/TIBCOSoftware/bego/common/model"
	"strings"
)

type ruleImpl struct {
	name        string
	identifiers []model.TupleTypeAlias
	conditions  []model.Condition
	actionFn    model.ActionFunction
	priority    int
	deps 		map[model.TupleTypeAlias]map[string]bool
}

//NewRule ... Create a new rule
func NewRule(name string) model.MutableRule {
	rule := ruleImpl{}
	rule.initRuleImpl(name)
	return &rule
}

func (rule *ruleImpl) initRuleImpl(name string) {
	rule.name = name
	rule.identifiers =  []model.TupleTypeAlias{}
	rule.conditions = []model.Condition{}
	rule.deps = make (map[model.TupleTypeAlias]map[string]bool)
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

func (rule *ruleImpl) AddCondition(conditionName string, idrs []model.TupleTypeAlias, cfn model.ConditionEvaluator) {
	rule.addCond(conditionName, idrs, cfn, false)
}

func (rule *ruleImpl) addCond(conditionName string, idrs []model.TupleTypeAlias, cfn model.ConditionEvaluator, setIdr bool) {
	condition := newCondition(conditionName, rule, idrs, cfn)
	rule.conditions = append(rule.conditions, condition)

	for _, cidr := range idrs {
		if len (rule.identifiers) == 0 {
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

//func (rule *ruleImpl) String() string {
//	str := ""
//	str += "[Rule: (" + rule.name + "\n"
//	return str
//}

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

func (rule *ruleImpl) GetIdentifiers() []model.TupleTypeAlias {
	return rule.identifiers
}

func (rule *ruleImpl) SetAction(actionFn model.ActionFunction) {
	rule.actionFn = actionFn
}

////IdentifiersToString Take a slice of Identifiers and return a string representation
//func IdentifiersToString(identifiers []model.TupleTypeAlias) string {
//	str := ""
//	for _, idr := range identifiers {
//		str += string(idr) + ", "
//	}
//	return str
//}


func (rule *ruleImpl) AddConditionWithDependency(conditionName string, idrs []string, cFn model.ConditionEvaluator) {
	typeDepMap := map[model.TupleTypeAlias]bool{}
	//cwd := model.ConditionAndDep{"n1", []string{"p1", "p2", "p3"}}
	for _, idr := range idrs {
		aliasProp := strings.Split(idr, ".")

		alias := model.TupleTypeAlias(aliasProp[0])
		typeDepMap[alias] = true
		prop := aliasProp[1]

		propMap, found := rule.deps [alias]
		if !found {
			propMap = map[string]bool{}
			rule.deps[alias] = propMap
		}
		propMap[prop] = true
	}
	typeDeps := []model.TupleTypeAlias{}

	for key, _ := range typeDepMap {
		typeDeps = append(typeDeps, key)
	}

	rule.addCond(conditionName, typeDeps, cFn, true)


}

func (rule *ruleImpl) GetDeps() map[model.TupleTypeAlias]map[string]bool {
	return rule.deps
}