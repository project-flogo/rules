package ruleaction

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"

	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/ruleapi"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

const (
	ACTION_REF = "github.com/TIBCOSoftware/bego/ruleaction"
)

type RuleAction struct {
	rs ruleapi.RuleSession
}
type ActionFactory struct {
}

//todo fix this
var metadata = &action.Metadata{ID: ACTION_REF, Async: false}

func init() {
	action.RegisterFactory(ACTION_REF, &ActionFactory{})
}

func (ff *ActionFactory) Init() error {
	return nil
}

type ActionData struct {
	Ref string `json:"ref"`
}

func (ff *ActionFactory) New(config *action.Config) (action.Action, error) {

	ruleAction := &RuleAction{}
	ruleAction.rs = ruleapi.NewRuleSession()

	var actionData ActionData
	err := json.Unmarshal(config.Data, &actionData)
	if err != nil {
		return nil, fmt.Errorf("failed to read rule action data '%s' error '%s'", config.Id, err.Error())
	}

	loadRules(ruleAction.rs)

	return ruleAction, nil
}

// Metadata get the Action's metadata
func (a *RuleAction) Metadata() *action.Metadata {
	return metadata
}

// IOMetadata get the Action's IO metadata
func (a *RuleAction) IOMetadata() *data.IOMetadata {
	return nil
}

func (a *RuleAction) Run(ctx context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {

	defer func() {
		if r := recover(); r != nil {
			logger.Warnf("Unhandled Error executing rule action \n")

			// todo: useful for debugging
			logger.Debugf("StackTrace: %s", debug.Stack())
		}

	}()

	streamTuple := model.NewStreamTuple("n1") //n1 -> will be replaced by contextual information coming in the data


	queryParams := inputs["queryParams"].Value().(map[string]string)



	//map input data into stream tuples, only string. ignore the rest for now
	for key, value := range queryParams {
		streamTuple.SetString(key, value)
	}

	a.rs.Assert(streamTuple)
	return nil, nil
}

func loadRules(rs ruleapi.RuleSession) {
	rule := ruleapi.NewRule("n1.name == Bob")
	fmt.Printf("Rule added: [%s]\n", rule.GetName())
	rule.AddCondition("c1", []model.StreamSource{"n1"}, checkForBob) // check for name "Bob" in n1
	//in effect, fire the rule when name field in both tuples in "Bob"
	rule.SetActionFn(myActionFn)
	rs.AddRule(rule)
}

func checkForBob(ruleName string, condName string, tuples map[model.StreamSource]model.StreamTuple) bool {
	//This conditions filters on name="Bob"
	streamTuple := tuples["n1"]
	if streamTuple == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	name := streamTuple.GetString("name")
	return name == "Bob"
}

func myActionFn(ruleName string, tuples map[model.StreamSource]model.StreamTuple) {
	fmt.Printf("Rule fired: [%s]\n", ruleName)
	streamTuple1 := tuples["n1"]
	if streamTuple1 == nil {
		fmt.Println("Should not get nil tuples here in Action! This is an error")
	}
	name1 := streamTuple1.GetString("name")
	fmt.Printf("n1.name = [%s]", name1)
}
