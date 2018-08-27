package ruleaction

import (
	"github.com/TIBCOSoftware/bego/common/model"
	"testing"
	"time"
	"fmt"
)

func TestPkgFlowNormalWithDeps(t *testing.T) {

	rs, err := createRuleSessionAndRules()
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}

	loadPkgRulesWithDeps(rs)
	rs.Start(nil)

	pkgEvt, _ := model.NewTupleWithKeyValues("packageevent", "1")
	pkgEvt.SetString(nil, "next", "sfo")
	pkgEvt.SetString(nil, "status", "normal")
	pkgEvt.SetString(nil, "isnew", "true")

	err = rs.Assert(nil, pkgEvt)
	if err != nil {
		fmt.Printf("Error...[%s]\n", err)
		return
	}
	//time.Sleep(time.Second*20)
	scanEv, _ := model.NewTupleWithKeyValues("scanevent", "1")
	scanEv.SetString(nil, "curr", "sfo")
	scanEv.SetString(nil, "next", "ny")
	err = scanEv.SetValue(nil, "eta", 10)
	if err != nil {
		fmt.Printf("[%s]\n", err)
	}

	rs.Assert(nil, scanEv)

	scanEv1, _ := model.NewTupleWithKeyValues("scanevent",  "1")
	scanEv1.SetString(nil, "curr", "ny")
	scanEv1.SetString(nil, "next", "done")
	scanEv.SetString(nil, "next", "ny")

	rs.Assert(nil, scanEv1)

}

func TestPkgFlowNormalWithDepsWithTimeout(t *testing.T) {

	rs, err := createRuleSessionAndRules()
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}

	loadPkgRulesWithDeps(rs)
	rs.Start(nil)

	pkgEvt, _ := model.NewTupleWithKeyValues("packageevent", "1")
	pkgEvt.SetString(nil, "next", "sfo")
	pkgEvt.SetString(nil, "status", "normal")
	pkgEvt.SetString(nil, "isnew", "true")


	rs.Assert(nil, pkgEvt)
	//time.Sleep(time.Second*20)
	scanEv, _ := model.NewTupleWithKeyValues("scanevent", "1")
	scanEv.SetString(nil, "curr", "sfo")
	scanEv.SetString(nil, "next", "ny")
	err = scanEv.SetValue(nil, "eta", 3)
	if err != nil {
		fmt.Printf("[%s]\n", err)
	}

	rs.Assert(nil, scanEv)

	time.Sleep(time.Second * time.Duration(20))

}

func TestPkgFlowNormalWithDepsMapValues(t *testing.T) {

	rs, err := createRuleSessionAndRules()
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}

	loadPkgRulesWithDeps(rs)
	rs.Start(nil)

	pkgEvt, err := model.NewTupleWithKeyValues("packageevent", "1")
	if err != nil {
		fmt.Printf("error: [%s]\n", err)
		return
	}
	pkgEvt.SetString(nil, "next", "sfo")
	pkgEvt.SetString(nil, "status", "normal")
	pkgEvt.SetString(nil, "isnew", "true")
	fmt.Printf("Asserting package with key [%s]\n", pkgEvt.GetKey().String())

	rs.Assert(nil, pkgEvt)

	values := make (map[string]interface{})
	values["packageid"] = "1"
	values["curr"] = "sfo"
	values["next"] = "ny"
	values["eta"] = 5

	scanEv, _ := model.NewTuple("scanevent", values)

	rs.Assert(nil, scanEv)

	values = make (map[string]interface{})
	values["packageid"] = "1"
	values["curr"] = "ny"
	values["next"] = "done"
	values["eta"] = 5

	scanEv2, _ := model.NewTuple("scanevent", values)
	rs.Assert(nil, scanEv2)

	time.Sleep(time.Second * time.Duration(20))
}
