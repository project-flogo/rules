package ruleaction

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"

	"github.com/project-flogo/core/data/metadata"

	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/app/resource"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/support/log"
	rulecache "github.com/project-flogo/rules/cache"
	"github.com/project-flogo/rules/common"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/config"
	"github.com/project-flogo/rules/ruleapi"
)

const (
	sRuleSession   = "rulesession"
	sTupleDescFile = "tupleDescriptorFile"
	ivValues       = "values"
)

var actionMetadata = action.ToMetadata(&Settings{})

var manager *config.ResourceManager

//var resManager *config.ResManager

type Settings struct {
	RuleSessionURI string                  `json:"ruleSessionURI"`
	TupleDescFile  string                  `json:"tupleDescriptorFile"`
	Tds            []model.TupleDescriptor `json:"tds"`
}

func init() {
	action.Register(&RuleAction{}, &ActionFactory{})
}

type ActionFactory struct {
}

func (f *ActionFactory) Initialize(ctx action.InitContext) error {

	if manager != nil {
		return nil
	}

	manager = config.NewResourceManager()
	resource.RegisterLoader(config.RESTYPE_RULESESSION, manager)

	return nil
}

// New implements action.Factory.New
func (f *ActionFactory) New(cfg *action.Config) (action.Action, error) {

	settings := &Settings{}

	jsonSettings, err := json.Marshal(cfg.Settings)

	er := json.Unmarshal(jsonSettings, settings)
	if er != nil {
		return nil, err
	}

	rsCfg, err := manager.GetRuleActionDescriptor(settings.RuleSessionURI)
	if err != nil {
		return nil, err
	}

	if rsCfg == nil {
		return nil, fmt.Errorf("unable to resolve rulesession: %s", settings.RuleSessionURI)
	}

	if settings.TupleDescFile != "" {
		//Load the tuple descriptor file (relative to GOPATH)
		tupleDescAbsFileNm := common.GetAbsPathForResource(settings.TupleDescFile)
		tupleDescriptor := common.FileToString(tupleDescAbsFileNm)

		log.RootLogger().Info("Loaded tuple descriptor: \n%s\n", tupleDescriptor)

		//First register the tuple descriptors
		err := model.RegisterTupleDescriptors(tupleDescriptor)
		if err != nil {
			return nil, fmt.Errorf("failed to register tuple descriptors : %s", err.Error())
		}
	} else if settings.Tds != nil {
		err = model.RegisterTupleDescriptorsFromTds(settings.Tds)
		if err != nil {
			return nil, fmt.Errorf("failed to register tuple descriptors : %s", err.Error())
		}
	}

	ruleAction := &RuleAction{}
	ruleSessionDescriptor, err := manager.GetRuleSessionDescriptor(settings.RuleSessionURI)
	if err != nil {
		return nil, fmt.Errorf("failed to get RuleSessionDescriptor for %s\n%s", settings.RuleSessionURI, err.Error())
	}
	ruleCollectionJSON, err := json.Marshal(ruleSessionDescriptor)

	if err != nil {
		return nil, fmt.Errorf("failed to marshall RuleSessionDescriptor : %s", err.Error())
	}
	ruleAction.rs, err = ruleapi.GetOrCreateRuleSessionFromConfig(settings.RuleSessionURI, string(ruleCollectionJSON))

	if err != nil {
		return nil, fmt.Errorf("failed to create rulesession for %s\n %s", settings.RuleSessionURI, err.Error())
	}

	ruleAction.ioMetadata = rsCfg.IOMetadata

	//start the rule session here, calls the startup rule function
	err = ruleAction.rs.Start(nil)

	//Initialize CacheManager
	var rcm *rulecache.RedisCacheManager = &rulecache.RedisCacheManager{}

	if rsCfg.CacheConfig != nil {
		rcm.Init(*rsCfg.CacheConfig)

		//Load tuples from cache
		tds := model.GetAllTupleDescriptors()
		for _, td := range tds {
			if model.OMModeMap[td.PersistMode] == model.ReadOnlyCache {
				err = rcm.LoadTuples(context.TODO(), &td, ruleAction.rs)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return ruleAction, err
}

// RuleAction wraps RuleSession
type RuleAction struct {
	rs         model.RuleSession
	ioMetadata *metadata.IOMetadata
}

func (a *RuleAction) Metadata() *action.Metadata {
	return actionMetadata
}

func (a *RuleAction) IOMetadata() *metadata.IOMetadata {
	return a.ioMetadata
}

// Run implements action.Action.Run
func (a *RuleAction) Run(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {

	defer func() {
		if r := recover(); r != nil {
			log.RootLogger().Warnf("Unhandled Error executing rule action \n")

			// todo: useful for debugging
			log.RootLogger().Debugf("StackTrace: %s", debug.Stack())
		}

	}()

	tupleTypeData, exists := inputs["tupletype"]

	if !exists {
		err := fmt.Errorf("No tuple name recieved")
		log.RootLogger().Debugf(err.Error())
		return nil, err
	}

	str, _ := tupleTypeData.(string)
	tupleType := model.TupleType(str)
	valAttr, exists := inputs[ivValues]

	if !exists {
		err := fmt.Errorf("No values recieved")
		log.RootLogger().Debugf(err.Error())
		return nil, err
	}

	val, _ := valAttr.(string)
	valuesMap := make(map[string]interface{})

	//metadata section allows receiving 'values' string as json format. (i.e. 'name=Bob' as '{"Name":"Bob"}')
	err := json.Unmarshal([]byte(val), &valuesMap)

	if err != nil {
		err := fmt.Errorf("values for [%s] are malformed:\n %v", string(tupleType), val)
		log.RootLogger().Warnf(err.Error())
		return nil, err
	}

	td := model.GetTupleDescriptor(tupleType)
	if td == nil {
		err := fmt.Errorf("Tuple descriptor for type [%s] not found", string(tupleType))
		log.RootLogger().Warnf(err.Error())
		return nil, err
	}

	for _, keyProp := range td.GetKeyProps() {
		_, found := valuesMap[keyProp]
		if !found {
			//set unique ids to string key properties, if not present in the payload
			if td.GetProperty(keyProp).PropType == data.TypeString {
				uid, err := common.GetUniqueId()
				if err == nil {
					valuesMap[keyProp] = uid
				} else {
					err := fmt.Errorf("Failed to generate a unique id, discarding event [%s]", string(tupleType))
					log.RootLogger().Warnf(err.Error())
					return nil, err
				}
			}
		}
	}

	tuple, _ := model.NewTuple(tupleType, valuesMap)
	err = a.rs.Assert(ctx, tuple)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
