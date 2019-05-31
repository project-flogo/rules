package tests

/* solar check events
if the monthly bill is greater than 200 and the house doesn't have a solar panels installed,
the house with the matching parcel id should be a candidate for solar panel installation promotion.

For now, house tuples are to be pre-loaded from data.txt to cache before running go test by
cat data.txt | redis-cli --pipe
The command above loads tuples with a tuple key as a hash key. Then updates the index named after
tuple type with the hash key.

<house data tuples>
house:parcel:0001 parcel 0001 is_solar true
house:parcel:0002 parcel 0002 is_solar false

solar events are asserted everytime a monthly electiricity bill is generated.

<solar event tuples>
solar:parcel:0001 parcel 0001 bill 300
solar:parcel:0002 parcel 0002 bill 250

*/

import (
	"context"
	"fmt"
	"testing"

	rulecache "github.com/project-flogo/rules/cache"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/config"

	"github.com/project-flogo/rules/ruleapi"
)

func Test_Cache_1(t *testing.T) {

	rs, _ := createRuleSession()

	rule := ruleapi.NewRule("Cache_Test")
	err := rule.AddCondition("TC_1", []string{"house", "solar"}, checkSolarEligibleCondition, nil)
	if err != nil {
		t.Logf("%s", err)
		t.FailNow()
	}
	ruleActionCtx := make(map[string]string)
	rule.SetContext(ruleActionCtx)
	rule.SetAction(solarEligibleAction)
	rule.SetPriority(1)
	err = rs.AddRule(rule)
	if err != nil {
		t.Logf("%s", err)
		t.FailNow()
	}
	t.Logf("Rule added: [%s]\n", rule.GetName())

	err = rs.Start(nil)
	if err != nil {
		t.Logf("%s", err)
		t.FailNow()
	}

	//PreCase: Assert all house tuples from Cache. ToDo: Check if this can be done before rulesession creation.
	{
		var rcm rulecache.CacheManager = &rulecache.RedisCacheManager{}
		cacheConfig := config.CacheConfig{"rediscache", "redis", "localhost:6379", "", 0}
		rcm.Init(cacheConfig)

		//Load tuples from cache
		tds := model.GetAllTupleDescriptors()
		for _, td := range tds {
			if model.OMModeMap[td.PersistMode] == model.ReadOnlyCache {
				err = rcm.LoadTuples(context.TODO(), &td, rs)
				if err != nil {
					t.Logf("%s", err)
					t.FailNow()
				}
			}
		}
	}

	// Case1: Assert an ineligible solar. It should not fire solarEligibleAction.
	{
		ctx := context.WithValue(context.TODO(), "key", t)
		values := make(map[string]interface{})
		values["parcel"] = "0001"
		values["bill"] = 300
		tuple, _ := model.NewTuple("solar", values)
		err := rs.Assert(ctx, tuple)
		if err != nil {
			t.Logf("%s", err)
			t.FailNow()
		}
	}

	// Case2: Assert an eligible solar. solarEligibleAction should be fired.
	{
		ctx := context.WithValue(context.TODO(), "key", t)
		values := make(map[string]interface{})
		values["parcel"] = "0002"
		values["bill"] = 250
		tuple, _ := model.NewTuple("solar", values)
		err := rs.Assert(ctx, tuple)
		if err != nil {
			t.Logf("%s", err)
			t.FailNow()
		}
	}

	// Case3: Assert the ineligible solar again. It should not fire solarEligibleAction.
	{
		ctx := context.WithValue(context.TODO(), "key", t)
		values := make(map[string]interface{})
		values["parcel"] = "0001"
		values["bill"] = 300
		tuple, _ := model.NewTuple("solar", values)
		err := rs.Assert(ctx, tuple)
		if err != nil {
			t.Logf("%s", err)
			t.FailNow()
		}
	}

	// Case2: Assert an eligible solar. solarEligibleAction should be fired.
	{
		ctx := context.WithValue(context.TODO(), "key", t)
		values := make(map[string]interface{})
		values["parcel"] = "0002"
		values["bill"] = 250
		tuple, _ := model.NewTuple("solar", values)
		err := rs.Assert(ctx, tuple)
		if err != nil {
			t.Logf("%s", err)
			t.FailNow()
		}
	}

	rs.Unregister()

}

func solarEligibleAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t := ctx.Value("key").(*testing.T)
	tHouse := tuples["house"]
	tSolar := tuples["solar"]
	t.Logf("Eligible for a solar promotion! [%s], [%s]\n", tHouse.GetKey().String(), tSolar.GetKey().String())
}

func checkSolarEligibleCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	tHouse := tuples["house"]
	tSolar := tuples["solar"]
	if tHouse == nil || tSolar == nil {
		fmt.Println("Should not get nil tuples here in JoinCondition! This is an error")
		return false
	}
	parcelHouse, _ := tHouse.GetString("parcel")
	parcelSolar, _ := tSolar.GetString("parcel")

	isSolarHouse, _ := tHouse.GetBool("is_solar")
	billSolar, _ := tSolar.GetDouble("bill")

	return (parcelHouse == parcelSolar) &&
		(isSolarHouse == false) &&
		(billSolar > 200)
}
