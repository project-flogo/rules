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

func (rs *rulesessionImpl) initRuleSession() {
	rs.reteNetwork = rete.NewReteNetwork()
}

func (rs *rulesessionImpl) AddRule(apiRule Rule) (int, bool) {
	rule := convertAPIRuleToReteRule(apiRule)

	ret := rs.reteNetwork.AddRule(rule)
	if ret == 0 {
		return 0, true
	}
	return ret, false
}

func (rs *rulesessionImpl) DeleteRule(ruleName string) {
	rs.reteNetwork.RemoveRule(ruleName)
}

func (rs *rulesessionImpl) Assert(tuple model.StreamTuple) {
	rs.reteNetwork.Assert(tuple)
}

func (rs *rulesessionImpl) Retract(tuple model.StreamTuple) {
	rs.reteNetwork.Retract(tuple)
}

func (rs *rulesessionImpl) PrintNetwork() {
	fmt.Println(rs.reteNetwork.String())
}
func convertAPIRuleToReteRule(apiRule Rule) rete.Rule {
	reteRule := rete.NewRule(apiRule.GetName())
	for _, c := range apiRule.GetConditions() {
		reteRule.AddCondition(c.GetName(), c.GetStreamSource(), c.GetEvaluator())
	}
	reteRule.SetAction(apiRule.GetActionFn())
	reteRule.SetPriority(apiRule.GetPriority())
	return reteRule
}
