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
		log.RootLogger().Debugf("No tuple name recieved")
		//no input, should we return an error?
		return nil, nil
	}

	str, _ := tupleTypeData.(string)
	tupleType := model.TupleType(str)
	valAttr, exists := inputs[ivValues]

	if !exists {
		log.RootLogger().Debugf("No values recieved")
		//no input, should we return an error?
		return nil, nil
	}

	val, _ := valAttr.(string)
	valuesMap := make(map[string]interface{})

	//metadata section allows receiving 'values' string as json format. (i.e. 'name=Bob' as '{"Name":"Box"}')
	err := json.Unmarshal([]byte(val), &valuesMap)

	if err != nil {
		log.RootLogger().Warnf("values for [%s] are malformed:\n %v\n", string(tupleType), val)
		return nil, nil
	}

	td := model.GetTupleDescriptor(tupleType)
	if td == nil {
		log.RootLogger().Warnf("Tuple descriptor for type [%s] not found\n", string(tupleType))
		return nil, nil
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
					log.RootLogger().Warnf("Failed to generate a unique id, discarding event [%s]\n", string(tupleType))
					return nil, nil
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
