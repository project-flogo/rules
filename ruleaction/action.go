package ruleaction

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/project-flogo/core/data/metadata"
	"runtime/debug"

	"github.com/project-flogo/core/app/resource"
	"github.com/project-flogo/core/action"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/trigger"
	"github.com/project-flogo/core/support/log"
	"github.com/project-flogo/rules/common"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/config"
	"github.com/project-flogo/rules/ruleapi"
)
const (
	sRuleSession   = "rulesession"
	sTupleDescFile = "tupleDescriptorFile"
	ivValues = "values"
)
var actionMetadata = action.ToMetadata(&Settings{})
var manager *config.ResourceManager

type Settings struct {
	RuleSessionURI string `md:"ruleSessionURI,required"`
	TupleDescFile  string `md:"tupleDescriptorFile"`
	Tds json.RawMessage `md:"tds"`
}
type Input struct {
	Data string `md:"data"`
}

//// ActionData maintains Tuple descriptor details
//type ActionData struct {
//	Tds json.RawMessage `json:"tds"`
//}

func init() {
	action.Register(&RuleAction{}, &ActionFactory{})
}

type ActionFactory struct {
	//resManager *resource.Manager
}

func (f *ActionFactory) Initialize(ctx action.InitContext) error {

	//f.resManager = ctx.ResourceManager()
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
	err := metadata.MapToStruct(cfg.Settings, settings, true)
	if err != nil {
		return nil, err
	}

	rsCfg, err := manager.GetRuleSessionDescriptor(settings.RuleSessionURI)

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
		//actionData := ActionData{}

		//err := json.Unmarshal(cfg.Data, &actionData)
		//if err != nil {
		//	return nil, fmt.Errorf("failed to read rule action data '%s' error '%s'", cfg.Id, err.Error())
		//}

		err = model.RegisterTupleDescriptors(string(settings.Tds))
		if err != nil {
			return nil, fmt.Errorf("failed to register tuple descriptors : %s", err.Error())
		}
	}

	ruleAction := &RuleAction{}
	ruleCollectionJSON, err := json.Marshal(rsCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall RuleSessionDescriptor : %s", err.Error())
	}
	ruleAction.rs, _ = ruleapi.GetOrCreateRuleSessionFromConfig(settings.RuleSessionURI, string(ruleCollectionJSON))

	//start the rule session here, calls the startup rule function
	err = ruleAction.rs.Start(nil)

	return ruleAction, err
}

// RuleAction wraps RuleSession
type RuleAction struct {
	rs model.RuleSession
}

func (a *RuleAction) Metadata() *action.Metadata {
	return actionMetadata
}

func (a *RuleAction) IOMetadata() *metadata.IOMetadata {
	return actionMetadata.IOMetadata
}
// Run implements action.Action.Run
func (a *RuleAction) Run(ctx context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {

	defer func() {
		if r := recover(); r != nil {
			log.RootLogger().Warnf("Unhandled Error executing rule action \n")

			// todo: useful for debugging
			log.RootLogger().Debugf("StackTrace: %s", debug.Stack())
		}

	}()

	h, _ok := trigger.HandlerFromContext(ctx)
	if !_ok {
		return nil, nil
	}

	tupleType := model.TupleType(h.Name)
	valAttr, exists := inputs[ivValues]
	if !exists {
		log.RootLogger().Debugf("No values recieved")
		//no input, should we return an error?
		return nil, nil
	}

	values := valAttr.Value().(map[string]interface{})

	td := model.GetTupleDescriptor(tupleType)
	if td == nil {
		log.RootLogger().Warnf("Tuple descriptor for type [%s] not found\n", string(tupleType))
		return nil, nil
	}

	for _, keyProp := range td.GetKeyProps() {
		_, found := values[keyProp]
		if !found {
			//set unique ids to string key properties, if not present in the payload
			if td.GetProperty(keyProp).PropType == data.TypeString {
				uid, err := common.GetUniqueId()
				if err == nil {
					values[keyProp] = uid
				} else {
					log.RootLogger().Warnf("Failed to generate a unique id, discarding event [%s]\n", string(tupleType))
					return nil, nil
				}
			}
		}
	}

	tuple, _ := model.NewTuple(tupleType, values)
	a.rs.Assert(ctx, tuple)
	// does this return anything?

	//fmt.Printf("[%s]\n", "b")
	return nil, nil
}
//
//func getSettings(config *action.Config) (*Settings, error) {
//
//	settings := &Settings{}
//
//	setting, exists := config.Settings[sRuleSession]
//	if exists {
//		val, err := coerce.ToString(setting)
//		if err != nil {
//			return nil, err
//		}
//		settings.RuleSessionURI = val
//	} else {
//		return nil, fmt.Errorf("RuleSession not specified")
//	}
//
//	setting, exists = config.Settings[sTupleDescFile]
//	if exists {
//		val, err := coerce.ToString(setting)
//		if err != nil {
//			return nil, err
//		}
//		settings.TupleDescFile = val
//	}
//
//	return settings, nil
//}
