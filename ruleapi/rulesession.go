package ruleapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	utils "github.com/project-flogo/rules/common"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/config"
	"github.com/project-flogo/rules/rete"
	"github.com/project-flogo/rules/rete/common"
	"github.com/project-flogo/rules/ruleapi/internal/store/mem"
	"github.com/project-flogo/rules/ruleapi/internal/store/redis"
)

var (
	sessionMap sync.Map
)

type rulesessionImpl struct {
	sync.RWMutex

	name        string
	reteNetwork common.Network

	timers      map[interface{}]*time.Timer
	startupFn   model.StartupRSFunction
	started     bool
	storeConfig string
	tupleStore  model.TupleStore
	config      common.Config
}

func ClearSessions() {
	sessionMap.Range(func(key, value interface{}) bool {
		sessionMap.Delete(key)
		return true
	})
}

func GetOrCreateRuleSession(name, config string) (model.RuleSession, error) {
	if name == "" {
		return nil, errors.New("RuleSession name cannot be empty")
	}
	rs := rulesessionImpl{}
	rs.loadStoreConfig(config)
	rs.initRuleSession(name)

	rs1, _ := sessionMap.LoadOrStore(name, &rs)
	return rs1.(*rulesessionImpl), nil
}

// GetOrCreateRuleSessionFromConfig returns rule session from created from config
func GetOrCreateRuleSessionFromConfig(name, store, jsonConfig string) (model.RuleSession, error) {
	rs, err := GetOrCreateRuleSession(name, store)

	if err != nil {
		return nil, err
	}

	ruleSessionDescriptor := config.RuleSessionDescriptor{}
	err = json.Unmarshal([]byte(jsonConfig), &ruleSessionDescriptor)
	if err != nil {
		return nil, err
	}

	// inflate action services
	aServices := make(map[string]model.ActionService)
	for _, s := range ruleSessionDescriptor.Services {
		aService, err := NewActionService(s)
		if err != nil {
			return nil, err
		}
		aServices[s.Name] = aService
	}

	for _, ruleCfg := range ruleSessionDescriptor.Rules {
		rule := NewRule(ruleCfg.Name)
		rule.SetContext("This is a test of context")
		// set action service to rule, if exist
		if ruleCfg.ActionService != nil {
			aService, found := aServices[ruleCfg.ActionService.Service]
			if !found {
				return nil, fmt.Errorf("rule action service[%s] not found", ruleCfg.ActionService.Service)
			}
			aService.SetInput(ruleCfg.ActionService.Input)
			rule.SetActionService(aService)
		}
		rule.SetPriority(ruleCfg.Priority)

		for _, condCfg := range ruleCfg.Conditions {
			if condCfg.Expression == "" {
				rule.AddCondition(condCfg.Name, condCfg.Identifiers, condCfg.Evaluator, nil)
			} else {
				rule.AddExprCondition(condCfg.Name, condCfg.Expression, nil)
			}
		}
		//now add explicit rule identifiers if any
		if ruleCfg.Identifiers != nil {
			idrs := []model.TupleType{}
			for _, idr := range ruleCfg.Identifiers {
				idrs = append(idrs, model.TupleType(idr))
			}
			rule.AddIdrsToRule(idrs)
		}

		rs.AddRule(rule)
	}

	rs.SetStartupFunction(config.GetStartupRSFunction(name))

	return rs, nil
}

const defaultConfig = `{
	"mode": "consistency",
  "rs": {
    "prefix": "x",
    "store-ref": "mem"
  },
  "rete": {
    "jt-ref": "mem",
    "idgen-ref": "mem",
    "jt":"mem"
  },
  "stores": {
    "mem": {
    },
    "redis": {
      "network": "tcp",
      "address": ":6379"
    }
  },
  "idgens": {
    "mem": {
    },
    "redis": {
      "network": "tcp",
      "address": ":6379"
    }
  },
  "jts": {
    "mem": {
    },
    "redis": {
      "network": "tcp",
      "address": ":6379"
    }
  }
}`

func (rs *rulesessionImpl) loadStoreConfig(name string) {
	if value := os.Getenv("STORECONFIG"); value != "" {
		name = value
	}

	if _, err := os.Stat(name); err != nil {
		rs.storeConfig = defaultConfig
		return
	}

	rs.storeConfig = utils.FileToString(name)
}

func (rs *rulesessionImpl) initRuleSession(name string) error {

	err := json.Unmarshal([]byte(rs.storeConfig), &rs.config)
	if err != nil {
		return err
	}

	rs.name = name
	rs.timers = make(map[interface{}]*time.Timer)
	rs.reteNetwork = rete.NewReteNetwork(rs.name, rs.storeConfig)

	//TODO: Configure it from jconsonfig
	tupleStore := getTupleStore(rs.config)
	if tupleStore != nil {
		tupleStore.Init()
		rs.SetStore(tupleStore)
	}
	rs.started = false
	return nil
}

func getTupleStore(config common.Config) model.TupleStore {
	switch config.Rs.StoreRef {
	case common.ServiceTypeMem:
		return mem.NewStore(config)
	case common.ServiceTypeRedis:
		return redis.NewStore(config)
	default:
		panic("invalid service type")
	}
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
	if ctx == nil {
		ctx = context.Context(context.Background())
	}

	return rs.reteNetwork.Assert(ctx, rs, tuple, nil, common.ADD)
}

func (rs *rulesessionImpl) Retract(ctx context.Context, tuple model.Tuple) error {
	return rs.reteNetwork.Retract(ctx, rs, tuple, nil, common.RETRACT)
}

func (rs *rulesessionImpl) Delete(ctx context.Context, tuple model.Tuple) error {
	return rs.reteNetwork.Retract(ctx, rs, tuple, nil, common.DELETE)
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
		rs.Lock()
		defer rs.Unlock()
		ctxNew := context.TODO()
		delete(rs.timers, key)
		rs.Assert(ctxNew, tuple)
	})

	rs.Lock()
	defer rs.Unlock()
	rs.timers[key] = timer
}

func (rs *rulesessionImpl) CancelScheduledAssert(ctx context.Context, key interface{}) {
	rs.RLock()
	timer, ok := rs.timers[key]
	rs.RUnlock()
	if ok {
		rs.Lock()
		defer rs.Unlock()
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

func (rs *rulesessionImpl) GetAssertedTuple(ctx context.Context, key model.TupleKey) model.Tuple {
	return rs.reteNetwork.GetAssertedTuple(ctx, rs, key)
}

func (rs *rulesessionImpl) RegisterRtcTransactionHandler(txnHandler model.RtcTransactionHandler, txnContext interface{}) {
	rs.reteNetwork.RegisterRtcTransactionHandler(txnHandler, txnContext)
}

func (rs *rulesessionImpl) GetStore() model.TupleStore {
	return rs.tupleStore
}

func (rs *rulesessionImpl) SetStore(store model.TupleStore) error {
	if store == nil {
		return fmt.Errorf("Cannot set nil store")
	}
	if rs.tupleStore != nil {
		return fmt.Errorf("TupleStore already set")
	}
	if rs.started {
		return fmt.Errorf("RuleSession already started")
	}
	rs.tupleStore = store
	rs.reteNetwork.SetTupleStore(store)
	rs.reteNetwork.RegisterRtcTransactionHandler(internalTxnHandler, nil)
	return nil
}

func internalTxnHandler(ctx context.Context, rs model.RuleSession, rtxn model.RtcTxn, handlerCtx interface{}) {
	store := rs.GetStore()
	store.DeleteTuples(rtxn.GetRtcDeleted())
	store.SaveTuples(rtxn.GetRtcAdded())
	store.SaveModifiedTuples(rtxn.GetRtcModified())
}
