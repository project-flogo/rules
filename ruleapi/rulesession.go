package ruleapi

import (
	"context"
	"fmt"

	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/rete"
)

type rulesessionImpl struct {
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

func (rs *rulesessionImpl) AddRule(rule model.Rule) (int, bool) {
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

func (rs *rulesessionImpl) printNetwork() {
	fmt.Println(rs.reteNetwork.String())
}
