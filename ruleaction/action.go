package ruleaction

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/config"
	"github.com/TIBCOSoftware/flogo-lib/app/resource"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/TIBCOSoftware/bego/common"
)

// Action ref to register the action factory
const (
	ActionRef = "github.com/TIBCOSoftware/bego/ruleaction"
)

var manager *config.ResourceManager

const (
	sRuleSession   = "rulesession"
	sTupleDescFile = "tupleDescriptorFile"

	ivValues = "values"
)

type Settings struct {
	RuleSessionURI string
	TupleDescFile  string
}

// RuleAction wraps RuleSession
type RuleAction struct {
	rs model.RuleSession
}

// ActionFactory wrapper to register with the action
type ActionFactory struct {
}

//todo fix this
var metadata = &action.Metadata{ID: ActionRef, Async: false,
	Settings: map[string]*data.Attribute{sRuleSession: data.NewZeroAttribute(sRuleSession, data.TypeString),
		sTupleDescFile: data.NewZeroAttribute(sTupleDescFile, data.TypeString)},
	Input: map[string]*data.Attribute{ivValues: data.NewZeroAttribute(ivValues, data.TypeObject)}}

//todo fix this
var iometadata = &data.IOMetadata{Input: map[string]*data.Attribute{ivValues: data.NewZeroAttribute(ivValues, data.TypeObject)}}

func init() {
	action.RegisterFactory(ActionRef, &ActionFactory{})
}

// Init implements action.Factory.Init
func (f *ActionFactory) Init() error {

	if manager != nil {
		return nil
	}

	manager = config.NewResourceManager()
	resource.RegisterManager(config.RESTYPE_RULESESSION, manager)

	return nil
}

// ActionData maintains Tuple descriptor details
type ActionData struct {
	Tds json.RawMessage `json:"tds"`
}

// New implements action.Factory.New
func (f *ActionFactory) New(cfg *action.Config) (action.Action, error) {

	settings, err := getSettings(cfg)
	if err != nil {
		return nil, err
	}

	rsCfg, err := manager.GetRuleSessionConfig(settings.RuleSessionURI)

	if err != nil {
		return nil, err
	} else {
		if rsCfg == nil {
			return nil, fmt.Errorf("unable to resolve rulesession: %s", settings.RuleSessionURI)
		}
	}

	if settings.TupleDescFile != "" {
		//Load the tuple descriptor file (relative to GOPATH)
		tupleDescAbsFileNm := getAbsPathForResource(settings.TupleDescFile)
		tupleDescriptor := fileToString(tupleDescAbsFileNm)

		logger.Info("Loaded tuple descriptor: \n%s\n", tupleDescriptor)

		//First register the tuple descriptors
		err := model.RegisterTupleDescriptors(tupleDescriptor)
		if err != nil {
			return nil, fmt.Errorf("failed to register tuple descriptors : %s", err.Error())
		}
	} else {
		actionData := ActionData{}

		err := json.Unmarshal(cfg.Data, &actionData)
		if err != nil {
			return nil, fmt.Errorf("failed to read rule action data '%s' error '%s'", cfg.Id, err.Error())
		}

		err = model.RegisterTupleDescriptors(string(actionData.Tds))
		if err != nil {
			return nil, fmt.Errorf("failed to register tuple descriptors : %s", err.Error())
		}
	}

	ruleAction := &RuleAction{}
	ruleAction.rs, _ = config.GetOrCreateRuleSessionFromConfig(settings.RuleSessionURI, rsCfg)

	//start the rule session here, calls the startup rule function
	err = ruleAction.rs.Start(nil)

	return ruleAction, err
}

// Metadata get the Action's metadata
func (a *RuleAction) Metadata() *action.Metadata {
	return metadata
}

// IOMetadata get the Action's IO metadata
func (a *RuleAction) IOMetadata() *data.IOMetadata {
	return iometadata
}

// Run implements action.Action.Run
func (a *RuleAction) Run(ctx context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {

	defer func() {
		if r := recover(); r != nil {
			logger.Warnf("Unhandled Error executing rule action \n")

			// todo: useful for debugging
			logger.Debugf("StackTrace: %s", debug.Stack())
		}

	}()

	h, _ok := trigger.HandlerFromContext(ctx)
	if !_ok {
		return nil, nil
	}

	tupleType := model.TupleType(h.Name)
	valAttr, exists := inputs[ivValues]
	if !exists {
		logger.Debugf("No values recieved")
		//no input, should we return an error?
		return nil, nil
	}

	values := valAttr.Value().(map[string]interface{})

	td := model.GetTupleDescriptor(tupleType)
	if td == nil {
		logger.Warnf("Tuple descriptor for type [%s] not found\n", string(tupleType))
		return nil, nil
	}

	for _, keyProp := range td.GetKeyProps() {
		_, found := values[keyProp]
		if !found {
			//set unique ids to string key properties, if not present in the payload
			if td.GetProperty(keyProp).PropType == data.TypeString {
				uid, err := common.GetUniqueId()
				if err != nil {
					values[keyProp] = uid
				} else {
					logger.Warnf("Failed to generate a unique id, discarding event [%s]\n", string(tupleType))
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

func getSettings(config *action.Config) (*Settings, error) {

	settings := &Settings{}

	setting, exists := config.Settings[sRuleSession]
	if exists {
		val, err := data.CoerceToString(setting)
		if err != nil {
			return nil, err
		}
		settings.RuleSessionURI = val
	} else {
		return nil, fmt.Errorf("RuleSession not specified")
	}

	setting, exists = config.Settings[sTupleDescFile]
	if exists {
		val, err := data.CoerceToString(setting)
		if err != nil {
			return nil, err
		}
		settings.TupleDescFile = val
	}

	return settings, nil
}

//////////////
// File Utils

func getAbsPathForResource(resourcepath string) string {
	GOPATH := os.Getenv("GOPATH")
	fmt.Printf("path[%s]\n", GOPATH)
	paths := strings.Split(GOPATH, ":")
	for _, path := range paths {
		fmt.Printf("path[%s]\n", path)
		absPath := path + "/" + resourcepath
		_, err := os.Stat(absPath)
		if err == nil {
			return absPath
		}
	}
	return ""
}

func fileToString(fileName string) string {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(dat)
}
