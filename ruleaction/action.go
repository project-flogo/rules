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
	//"io/ioutil"
	//"log"
	"io/ioutil"
	"log"
)

const (
	ACTION_REF = "github.com/TIBCOSoftware/bego/ruleaction"
)

type RuleAction struct {
	rs model.RuleSession
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
	Tds []model.TupleDescriptor
}

func (ff *ActionFactory) New(config *action.Config) (action.Action, error) {

	ruleAction := &RuleAction{}

	ruleAction.rs = ruleapi.GetOrCreateRuleSession("flogosession")

	actionData := ActionData{}
	actionData.Tds = []model.TupleDescriptor{}

	err := json.Unmarshal(config.Data, &actionData)
	if err != nil {
		return nil, fmt.Errorf("failed to read rule action data '%s' error '%s'", config.Id, err.Error())
	}
	tdss , _:= json.Marshal(&actionData.Tds)
	fmt.Printf ("**ACTION DATA: [%s]\n**", string(tdss))

	ruleAction.rs.RegisterTupleDescriptors(string(tdss))

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
	//fmt.Printf("Received event from stream source [%s]", h.Name)

	streamSrc := model.TupleTypeAlias(h.Name)

	streamTuple := model.NewStreamTuple(streamSrc) //n1 -> will be replaced by contextual information coming in the data


	queryParams := inputs["queryParams"].Value().(map[string]string)


	//map input data into stream tuples, only string. ignore the rest for now
	for key, value := range queryParams {
		fmt.Printf("[%s]\n", "a")
		if key == "balance" {
			f, _ := strconv.ParseFloat(value, 64)
			streamTuple.SetFloat(ctx, a.rs, key, f)
		} else {
			streamTuple.SetString(ctx, a.rs, key, value)
		}
	}

	a.rs.Assert(ctx, streamTuple)
	fmt.Printf("[%s]\n", "b")
	return nil, nil
}

func loadRules(rs model.RuleSession) {

	rule := ruleapi.NewRule("customer-event")
	rule.AddCondition("customer", []model.TupleTypeAlias{"customerevent"}, truecondition) // check for name "Bob" in n1
	rule.SetAction(customerEvent)
	rule.SetPriority(1)
	rs.AddRule(rule)
	fmt.Printf("Rule added: [%s]\n", rule.GetName())

	rule2 := ruleapi.NewRule("debit-event")
	rule2.AddCondition("debitevent", []model.TupleTypeAlias{"debitevent"}, truecondition)
	rule2.SetAction(debitEvent)
	rule2.SetPriority(2)
	rs.AddRule(rule2)
	fmt.Printf("Rule added: [%s]\n", rule2.GetName())

	rule3 := ruleapi.NewRule("customer-debit")
	rule3.AddCondition("customerdebit", []model.TupleTypeAlias{"debitevent", "customerevent"}, customerdebitjoincondition)
	rule3.SetAction(debitAction)
	rule3.SetPriority(3)
	rs.AddRule(rule3)
	fmt.Printf("Rule added: [%s]\n", rule3.GetName())

	rule4 := ruleapi.NewRule("check-balance")
	rule4.AddCondition("customerdebit", []model.TupleTypeAlias{"customerevent"}, checkBalance) // check for name "Bob" in n1
	rule4.SetAction(balanceAlert)
	rule4.SetPriority(-1)
	rs.AddRule(rule4)
	fmt.Printf("Rule added: [%s]\n", rule4.GetName())


	rule5 := ruleapi.NewRule("timeout-rule")
	rule5.AddCondition("packagetimeout", []model.TupleTypeAlias{"packagetimeout"}, truecondition) // check for name "Bob" in n1
	rule5.SetAction(packageTimeout)
	rule5.SetPriority(-1)
	rs.AddRule(rule5)
	fmt.Printf("Rule added: [%s]\n", rule5.GetName())



}


func truecondition (ruleName string, condName string, tuples map[model.TupleTypeAlias]model.StreamTuple) bool {
	return true
}
func customerEvent(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {
	streamTuple := tuples["customerevent"]
	if streamTuple == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
	} else {
		name := streamTuple.GetString("name")
		fmt.Printf("Received a customer event with customer name [%s]\n", name)
	}
}
func debitEvent(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {
	streamTuple := tuples["debitevent"]
	if streamTuple == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
	} else {
		name := streamTuple.GetString("name")
		amount := streamTuple.GetString("debit")

		fmt.Printf("Received a debit event for customer [%s], amount [%s]\n", name, amount)
	}
}

func customerdebitjoincondition(ruleName string, condName string, tuples map[model.TupleTypeAlias]model.StreamTuple) bool {

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

func debitAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {
	//fmt.Printf("Rule fired: [%s]\n", ruleName)
	customerTuple := tuples["customerevent"].(model.MutableStreamTuple)
	debitTuple := tuples["debitevent"]
	dbt := debitTuple.GetString("debit")
	debitAmt, _ := strconv.ParseFloat(dbt, 64)
	currBal := customerTuple.GetFloat("balance")
	if (customerTuple.GetString("status") == "active") {
		customerTuple.SetFloat(ctx, rs,"balance", customerTuple.GetFloat("balance")-debitAmt)
	}
	fmt.Printf("Customer [%s], Balance [%f], Debit [%f], NewBalance [%f]\n", customerTuple.GetString("name"), currBal, debitAmt, customerTuple.GetFloat("balance"))
}

func checkBalance(ruleName string, condName string, tuples map[model.TupleTypeAlias]model.StreamTuple) bool {
	customerTuple := tuples["customerevent"]
	if customerTuple == nil {
		fmt.Println("Should not get a nil tuple here! This is an error")
		return false
	}
	balance := customerTuple.GetFloat("balance")
	issuspended := customerTuple.GetString("status")

	return balance <= 0 && issuspended == "active"
}


func balanceAlert(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {
	//fmt.Printf("Rule fired: [%s]\n", ruleName)
	customerTuple := tuples["customerevent"].(model.MutableStreamTuple)
	fmt.Printf("**** Account Suspended *** Customer balance is 0 or negative ! [%s], Balance [%f]\n", customerTuple.GetString("name"), customerTuple.GetFloat("balance"))
	customerTuple.SetString(ctx,rs,"status", "suspended")
}

func packageTimeout(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {
	packageTimeout := tuples["packagetimeout"]
	if packageTimeout == nil {
		fmt.Println("Should not get a nil tuple here! This is an error")
		return
	}
	fmt.Printf("Received a package timeout event for package-id [%s]\n", packageTimeout.GetString("packageid"))
}

func createRuleSessionAndRules() model.RuleSession {
	rs := ruleapi.GetOrCreateRuleSession("asession")
	dat, err := ioutil.ReadFile("/home/bala/go/src/github.com/TIBCOSoftware/bego/common/model/tupledescriptor.json")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Tuple descriptors: [%s]\n", string(dat))
	rs.RegisterTupleDescriptors(string(dat))
	loadRules(rs)
	return rs
}