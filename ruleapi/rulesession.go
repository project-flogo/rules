package ruleapi

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/rete"
)

var (
	sessionMap sync.Map
)

type rulesessionImpl struct {
	name        string
	reteNetwork rete.Network

	timers map[interface{}]*time.Timer
	startupFn model.StartupRSFunction
	started bool
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

func (rs *rulesessionImpl) initRuleSession(name string) {
	rs.reteNetwork = rete.NewReteNetwork()
	rs.name = name
	rs.timers = make(map[interface{}]*time.Timer)
	rs.started = false
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

func (rs *rulesessionImpl) Assert(ctx context.Context, tuple model.Tuple) (err error) {
	if !rs.started {
		return fmt.Errorf("Cannot assert tuple. Rulesession [%s] not started", rs.name)
	}
	if ctx == nil {
		ctx = context.Context(context.Background())
	}
	rs.reteNetwork.Assert(ctx, rs, tuple, nil)
	return nil
}

func (rs *rulesessionImpl) Retract(ctx context.Context, tuple model.Tuple) {
	rs.reteNetwork.Retract(ctx, tuple, nil)
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

func (rs *rulesessionImpl) SetStartupFunction (startupFn model.StartupRSFunction) {
	rs.startupFn = startupFn
}

func (rs *rulesessionImpl) GetStartupFunction() (startupFn model.StartupRSFunction) {
	return rs.startupFn
}

func (rs *rulesessionImpl) Start(startupCtx map[string]interface{}) error {

	if !rs.started {
		if rs.startupFn != nil {
			err := rs.startupFn (context.TODO(), rs, startupCtx)
			if err != nil {
				return err
			}
		}
		rs.started = true
	} else {
		return fmt.Errorf("Rulesession [%s] already started", rs.name)
	}
	return nil
}