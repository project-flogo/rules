package ruleapi

import (
	"github.com/TIBCOSoftware/bego/rete"
	"github.com/TIBCOSoftware/bego/common/model"
	"fmt"
)



import (
	"context"
	"sync"
	"encoding/json"
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
	rs.reteNetwork.Assert(ctx, rs, tuple, nil)
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

func (rs *rulesessionImpl) RegisterTupleDescriptors (jsonRegistry string) {

	tds := make([]model.TupleDescriptor, 0)

	json.Unmarshal([]byte(jsonRegistry),&tds)

	rs.reteNetwork.RegisterTupleDescriptors(tds)
}

func (rs *rulesessionImpl) ValidateUpdate(alias model.TupleTypeAlias, name string, value interface{}) bool {

	td := rs.reteNetwork.GetTupleDescriptor(alias)
	//TODO: type not registered, meaning no validation
	if  td == nil {
		return true
	}

	//TODO: Check property's type and value's type compatibility
	_, _ok := td.GetProperty(name)
	if !_ok {
		return false
	}

	return true
}

