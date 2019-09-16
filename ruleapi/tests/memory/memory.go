package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/project-flogo/rules/common"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

func createRuleSession() (model.RuleSession, error) {
	rs, _ := ruleapi.GetOrCreateRuleSession("test")

	tupleDescFileAbsPath := common.GetPathForResource("ruleapi/tests/tests.json", "./../tests.json")

	dat, err := ioutil.ReadFile(tupleDescFileAbsPath)
	if err != nil {
		log.Fatal(err)
	}
	err = model.RegisterTupleDescriptors(string(dat))
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func trueCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	return true
}

func emptyAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
}

var (
	adddelete   = flag.Bool("adddelete", false, "add followed by delete mode")
	addsdeletes = flag.Bool("addsdeletes", false, "adds followed by deletes mode")
	ttl         = flag.Bool("ttl", false, "add with ttl mode")
	same        = flag.Bool("same", false, "add same tuple")
)

func main() {
	flag.Parse()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	if *adddelete {
		rs, _ := createRuleSession()
		defer rs.Unregister()
		rule := ruleapi.NewRule("R2")
		rule.AddCondition("R2_c1", []string{"t3.none"}, trueCondition, nil)
		rule.SetAction(emptyAction)
		rule.SetPriority(1)
		rs.AddRule(rule)
		log.Printf("Rule added: [%s]\n", rule.GetName())
		rs.Start(nil)

		i := 0
		for {
			t1, _ := model.NewTupleWithKeyValues("t3", fmt.Sprintf("tuple%d", i))
			err := rs.Assert(context.TODO(), t1)
			if err != nil {
				log.Fatalf("err should be nil: %v", err)
			}
			time.Sleep(10 * time.Millisecond)
			rs.Delete(context.TODO(), t1)
			time.Sleep(10 * time.Millisecond)
			i++
		}
	} else if *addsdeletes {
		rs, _ := createRuleSession()
		defer rs.Unregister()
		rule := ruleapi.NewRule("R2")
		rule.AddCondition("R2_c1", []string{"t3.none"}, trueCondition, nil)
		rule.SetAction(emptyAction)
		rule.SetPriority(1)
		rs.AddRule(rule)
		log.Printf("Rule added: [%s]\n", rule.GetName())
		rs.Start(nil)

		i := 0
		for {
			tuples := make([]model.MutableTuple, 0, 10)
			for j := 0; j < 10; j++ {
				t1, _ := model.NewTupleWithKeyValues("t3", fmt.Sprintf("tuple%d_%d", i, j))
				err := rs.Assert(context.TODO(), t1)
				if err != nil {
					log.Fatalf("err should be nil: %v", err)
				}
				tuples = append(tuples, t1)
				time.Sleep(10 * time.Millisecond)
			}

			for _, tuple := range tuples {
				rs.Delete(context.TODO(), tuple)
				time.Sleep(10 * time.Millisecond)
			}
			i++
		}
	} else if *ttl {
		rs, _ := createRuleSession()
		defer rs.Unregister()
		rule := ruleapi.NewRule("R2")
		rule.AddCondition("R2_c1", []string{"t4.none"}, trueCondition, nil)
		rule.SetAction(emptyAction)
		rule.SetPriority(1)
		rs.AddRule(rule)
		log.Printf("Rule added: [%s]\n", rule.GetName())
		rs.Start(nil)

		i := 0
		for {
			t1, _ := model.NewTupleWithKeyValues("t4", fmt.Sprintf("tuple%d", i))
			err := rs.Assert(context.TODO(), t1)
			if err != nil {
				log.Fatalf("err should be nil: %v", err)
			}
			time.Sleep(time.Second)
			i++
		}
	} else if *same {
		rs, _ := createRuleSession()
		defer rs.Unregister()
		rule := ruleapi.NewRule("R2")
		rule.AddCondition("R2_c1", []string{"t3.none"}, trueCondition, nil)
		rule.SetAction(emptyAction)
		rule.SetPriority(1)
		rs.AddRule(rule)
		log.Printf("Rule added: [%s]\n", rule.GetName())
		rs.Start(nil)

		t1, _ := model.NewTupleWithKeyValues("t3", "t3")
		err := rs.Assert(context.TODO(), t1)
		if err != nil {
			log.Fatalf("err should be nil: %v", err)
		}
		for {
			rs.Assert(context.TODO(), t1)
			time.Sleep(10 * time.Millisecond)
		}
	}
}
