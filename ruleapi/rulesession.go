package ruleapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/config"
	"github.com/project-flogo/rules/rete"
)

var (
	sessionMap sync.Map
)

type rulesessionImpl struct {
	name        string
	reteNetwork rete.Network

	timers    map[interface{}]*time.Timer
	startupFn model.StartupRSFunction
	started   bool
	config    map[string]string
}

func GetOrCreateRuleSession(name string) (model.RuleSession, error) {
	if name == "" {
		return nil, errors.New("RuleSession name cannot be empty")
	}
	rs := rulesessionImpl{}
	rs.initRuleSession(name)
	rs1, _ := sessionMap.LoadOrStore(name, &rs)
	return rs1.(*rulesessionImpl), nil
}

func GetOrCreateRuleSessionFromConfig(name string, jsonConfig string) (model.RuleSession, error) {
	rs, err := GetOrCreateRuleSession(name)

	if err != nil {
		return nil, err
	}

	ruleSessionDescriptor := config.RuleSessionDescriptor{}
	err = json.Unmarshal([]byte(jsonConfig), &ruleSessionDescriptor)
	if err != nil {
		return nil, err
	}

	for _, ruleCfg := range ruleSessionDescriptor.Rules {
		rule := NewRule(ruleCfg.Name)
		rule.SetContext("This is a test of context")
		rule.SetAction(ruleCfg.ActionFunc)
		rule.SetPriority(ruleCfg.Priority)

		for _, condCfg := range ruleCfg.Conditions {
			rule.AddCondition(condCfg.Name, condCfg.Identifiers, condCfg.Evaluator, nil)
		}

		rs.AddRule(rule)
	}

	rs.SetStartupFunction(config.GetStartupRSFunction(name))

	return rs, nil
}

func (rs *rulesessionImpl) initRuleSession(name string) {
	rs.reteNetwork = rete.NewReteNetwork()
	rs.name = name
	rs.timers = make(map[interface{}]*time.Timer)
	rs.started = false
}

func (rs *rulesessionImpl) AddRule(rule model.Rule) (err error) {
	return rs.reteNetwork.AddRule(rule)
}

func (rs *rulesessionImpl) DeleteRule(ruleName string) {
	rs.reteNetwork.RemoveRule(ruleName)
}

func (rs *rulesessionImpl) GetRules() []model.Rule {
	return rs.reteNetwork.GetRules()
}

func (rs *rulesessionImpl) Assert(ctx context.Context, tuple model.Tuple) (err error) {
	if !rs.started {
		return fmt.Errorf("Cannot assert tuple. Rulesession [%s] not started", rs.name)
	}
	assertedTuple := rs.GetAssertedTuple(tuple.GetKey())
	if assertedTuple == tuple {
		return fmt.Errorf("Tuple with key [%s] already asserted", tuple.GetKey().String())
	} else if assertedTuple != nil {
		return fmt.Errorf("Tuple with key [%s] already asserted", tuple.GetKey().String())
	}
	if ctx == nil {
		ctx = context.Context(context.Background())
	}
	rs.reteNetwork.Assert(ctx, rs, tuple, nil, rete.ADD)
	return nil
}

func (rs *rulesessionImpl) Retract(ctx context.Context, tuple model.Tuple) {
	rs.reteNetwork.Retract(ctx, tuple, nil, rete.RETRACT)
}

func (rs *rulesessionImpl) Delete(ctx context.Context, tuple model.Tuple) {
	rs.reteNetwork.Retract(ctx, tuple, nil, rete.DELETE)
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

func (rs *rulesessionImpl) ScheduleAssert(ctx context.Context, delayInMillis uint64, key interface{}, tuple model.Tuple) {

	timer := time.AfterFunc(time.Millisecond*time.Duration(delayInMillis), func() {
		ctxNew := context.TODO()
		delete(rs.timers, key)
		rs.Assert(ctxNew, tuple)
	})

	rs.timers[key] = timer
}

func (rs *rulesessionImpl) CancelScheduledAssert(ctx context.Context, key interface{}) {
	timer, ok := rs.timers[key]
	if ok {
		fmt.Printf("Cancelling timer attached to key [%v]\n", key)
		delete(rs.timers, key)
		timer.Stop()
	}
}

func (rs *rulesessionImpl) SetStartupFunction(startupFn model.StartupRSFunction) {
	rs.startupFn = startupFn
}

func (rs *rulesessionImpl) GetStartupFunction() (startupFn model.StartupRSFunction) {
	return rs.startupFn
}

func (rs *rulesessionImpl) Start(startupCtx map[string]interface{}) error {

	if !rs.started {
		rs.started = true
		if rs.startupFn != nil {
			err := rs.startupFn(context.TODO(), rs, startupCtx)
			if err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("Rulesession [%s] already started", rs.name)
	}
	return nil
}

func (rs *rulesessionImpl) GetAssertedTuple(key model.TupleKey) model.Tuple {
	return rs.reteNetwork.GetAssertedTuple(key)
}

func (rs *rulesessionImpl) RegisterRtcTransactionHandler(txnHandler model.RtcTransactionHandler, txnContext interface{}) {
	rs.reteNetwork.RegisterRtcTransactionHandler(txnHandler, txnContext)
}

func (rs *rulesessionImpl) SetConfig(config map[string]string) {
	if rs.config == nil {
		rs.config = config
	}
	if rs.reteNetwork != nil && rs.reteNetwork.GetConfig() != nil {
		rs.reteNetwork.SetConfig(config)
	}
}
