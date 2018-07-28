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
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"strconv"
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

	h, _ok := trigger.HandlerFromContext(ctx)
	if !_ok {
		return nil, nil
	}
	fmt.Printf("Received event from stream source [%s]", h.Name)

	streamSrc := model.StreamSource(h.Name)

	streamTuple := model.NewStreamTuple(streamSrc) //n1 -> will be replaced by contextual information coming in the data


	queryParams := inputs["queryParams"].Value().(map[string]string)


	//map input data into stream tuples, only string. ignore the rest for now
	for key, value := range queryParams {
		streamTuple.SetString(key, value)
	}

	a.rs.Assert(streamTuple)
	return nil, nil
}

func loadRules(rs ruleapi.RuleSession) {

	rule := ruleapi.NewRule("customer-event")
	rule.AddCondition("customer", []model.StreamSource{"customerevent"}, truecondition) // check for name "Bob" in n1
	rule.SetActionFn(customer)
	rule.SetPriority(1)
	rs.AddRule(rule)
	fmt.Printf("Rule added: [%s]\n", rule.GetName())

	rule2 := ruleapi.NewRule("debit-event")
	rule2.AddCondition("debitevent", []model.StreamSource{"debitevent"}, truecondition)
	rule2.SetActionFn(debit)
	rule2.SetPriority(2)
	rs.AddRule(rule2)
	fmt.Printf("Rule added: [%s]\n", rule2.GetName())

	rule3 := ruleapi.NewRule("customer-debit")
	rule3.AddCondition("customerdebit", []model.StreamSource{"debitevent", "customerevent"}, customerdebitjoincondition)
	rule3.SetActionFn(debitAction)
	rule3.SetPriority(3)
	rs.AddRule(rule3)
	fmt.Printf("Rule added: [%s]\n", rule3.GetName())

	//rule4 := ruleapi.NewRule("check-balance")
	//rule4.AddCondition("customerdebit", []model.StreamSource{"customerevent"}, checkBalance) // check for name "Bob" in n1
	//rule4.SetActionFn(balanceAlert)
	//rule4.SetPriority(4)
	//rs.AddRule(rule4)
	//fmt.Printf("Rule added: [%s]\n", rule4.GetName())


}


func truecondition (ruleName string, condName string, tuples map[model.StreamSource]model.StreamTuple) bool {
	return true
}
func customer(ruleName string, tuples map[model.StreamSource]model.StreamTuple) {
	streamTuple := tuples["customerevent"]
	if streamTuple == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
	} else {
		name := streamTuple.GetString("name")
		fmt.Printf("Received a customer event with customer name [%s]\n", name)
	}
}
func debit(ruleName string, tuples map[model.StreamSource]model.StreamTuple) {
	streamTuple := tuples["debitevent"]
	if streamTuple == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
	} else {
		name := streamTuple.GetString("name")
		amount := streamTuple.GetString("debit")

		fmt.Printf("Received a debit event for customer [%s], amount [%s]\n", name, amount)
	}
}

func customerdebitjoincondition(ruleName string, condName string, tuples map[model.StreamSource]model.StreamTuple) bool {

	customerTuple := tuples["customerevent"]
	debitTuple := tuples["debitevent"]

	if customerTuple == nil || debitTuple == nil {
		fmt.Println("Should not get a nil tuple here! This is an error")
		return false
	}
	custName := customerTuple.GetString("name")
	acctName := debitTuple.GetString("name")

	return custName == acctName

}

func debitAction(ruleName string, tuples map[model.StreamSource]model.StreamTuple) {
	//fmt.Printf("Rule fired: [%s]\n", ruleName)
	customerTuple := tuples["customerevent"].(model.MutableStreamTuple)
	debitTuple := tuples["debitevent"]
	dbt := debitTuple.GetString("debit")
	debitAmt, _ := strconv.ParseFloat(dbt, 64)
	currBal := customerTuple.GetFloat("balance")
	customerTuple.SetFloat("balance", customerTuple.GetFloat("balance") - debitAmt)
	fmt.Printf("Customer [%s], Balance [%f], Debit [%f], NewBalance [%f]\n", customerTuple.GetString("name"), currBal, debitAmt, customerTuple.GetFloat("balance"))
}
