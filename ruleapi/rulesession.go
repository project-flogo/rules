package ruleapi

import (
	"github.com/TIBCOSoftware/bego/rete"
	"github.com/TIBCOSoftware/bego/common/model"
	"fmt"
)



import (
	"context"
	"sync"
)
var (
	sessionMap sync.Map
)
func init() {

}
type rulesessionImpl struct {
	name string
	reteNetwork rete.Network
}

func GetOrCreateRuleSession(name string) model.RuleSession {

	rs := rulesessionImpl{}
	rs.initRuleSession()
	rs1, _ :=  sessionMap.LoadOrStore(name, &rs)
	return rs1.(*rulesessionImpl)
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
	rs.reteNetwork.Retract(ctx, tuple)
}

func (rs *rulesessionImpl) printNetwork() {
	fmt.Println(rs.reteNetwork.String())
}

func (rs *rulesessionImpl) GetName() string {
	return rs.name
}

func (rs *rulesessionImpl) Unregister() {
	sessionMap.Delete(rs.name)
}
