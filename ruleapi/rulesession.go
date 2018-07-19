package ruleapi

import (
	"context"
	"fmt"

	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/rete"
)

type rulesessionImpl struct {
	allRules    map[string]model.Rule
	reteNetwork rete.Network
}

func NewRuleSession() model.RuleSession {
	rs := rulesessionImpl{}
	rs.initRuleSession()
	return &rs
}

func (rs *rulesessionImpl) initRuleSession() {
	rs.reteNetwork = rete.NewReteNetwork()
}

func (rs *rulesessionImpl) AddRule(apiRule model.Rule) (int, bool) {
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

func (rs *rulesessionImpl) Assert(ctx context.Context, tuple model.StreamTuple) {
	if ctx == nil {
		ctx = context.Context(context.Background())
	}
	rs.reteNetwork.Assert(ctx, rs, tuple)
}

func (rs *rulesessionImpl) Retract(ctx context.Context, tuple model.StreamTuple) {
	rs.reteNetwork.Retract(tuple)
}

func (rs *rulesessionImpl) PrintNetwork() {
	fmt.Println(rs.reteNetwork.String())
}
func convertAPIRuleToReteRule(apiRule model.Rule) rete.Rule {
	reteRule := rete.NewRule(apiRule.GetName())
	for _, c := range apiRule.GetConditions() {
		reteRule.AddCondition(c.GetName(), c.GetStreamSource(), c.GetEvaluator())
	}
	reteRule.SetAction(apiRule.GetActionFn())
	reteRule.SetPriority(apiRule.GetPriority())
	return reteRule
}
