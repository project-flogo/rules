package ruleapi

import (
	"fmt"

	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/rete"
)

type RuleSession interface {
	AddRule(rule Rule) (int, bool)
	DeleteRule(ruleName string)

	Assert(tuple model.StreamTuple)
	Retract(tuple model.StreamTuple)
	PrintNetwork()
}

type rulesessionImpl struct {
	allRules    map[string]Rule
	reteNetwork rete.Network
}

func NewRuleSession() RuleSession {
	rs := rulesessionImpl{}
	rs.initRuleSession()
	return &rs
}

func (rulesessionImplVar *rulesessionImpl) initRuleSession() {
	rulesessionImplVar.reteNetwork = rete.NewReteNetwork()
}

func (rulesessionImplVar *rulesessionImpl) AddRule(apiRule Rule) (int, bool) {
	rule := convertAPIRuleToReteRule(apiRule)

	ret := rulesessionImplVar.reteNetwork.AddRule(rule)
	if ret == 0 {
		return 0, true
	}
	return ret, false
}

func (rulesessionImplVar *rulesessionImpl) DeleteRule(ruleName string) {
	rulesessionImplVar.reteNetwork.RemoveRule(ruleName)
}

func (rulesessionImplVar *rulesessionImpl) Assert(tuple model.StreamTuple) {
	rulesessionImplVar.reteNetwork.Assert(tuple)
}

func (rulesessionImplVar *rulesessionImpl) Retract(tuple model.StreamTuple) {
	rulesessionImplVar.reteNetwork.Retract(tuple)
}

func (rulesessionImplVar *rulesessionImpl) PrintNetwork() {
	fmt.Println(rulesessionImplVar.reteNetwork.String())
}
func convertAPIRuleToReteRule(apiRule Rule) rete.Rule {
	reteRule := rete.NewRule(apiRule.GetName())
	for _, c := range apiRule.GetConditions() {
		reteRule.AddCondition(c.GetName(), c.GetStreamSource(), c.GetEvaluator())
	}
	reteRule.SetAction(apiRule.GetActionFn())
	return reteRule
}
