package dtable

import (
	"context"
	"testing"

	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support/test"
	"github.com/project-flogo/rules/common/model"
	"github.com/stretchr/testify/assert"

	"github.com/project-flogo/core/data/mapper"
)

const tupleDescriptor = `[
	{
	   "name":"applicant",
	   "properties":[
		  {
			 "name":"name",
			 "pk-index":0,
			 "type":"string"
		  },
		  {
			 "name":"gender",
			 "type":"string"
		  },
		  {
			 "name":"age",
			 "type":"int"
		  },
		  {
			 "name":"address",
			 "type":"string"
		  },
		  {
			 "name":"hasDL",
			 "type":"bool"
		  },
		  {
			 "name":"ssn",
			 "type":"long"
		  },
		  {
			 "name":"income",
			 "type":"double"
		  },
		  {
			 "name":"maritalStatus",
			 "type":"string"
		  },
		  {
			 "name":"creditScore",
			 "type":"int"
		  },
		  {
			 "name":"status",
			 "type":"string"
		  },
		  {
			 "name":"eligible",
			 "type":"bool"
		  },
		  {
			 "name":"creditLimit",
			 "type":"double"
		  }
	   ]
	},
	{
	   "name":"processapplication",
	   "ttl":0,
	   "properties":[
		  {
			 "name":"ssn",
			 "pk-index":0,
			 "type":"long"
		  },
		  {
			 "name":"start",
			 "type":"bool"
		  }
	   ]
	}
]`

var testApplicants = []map[string]interface{}{
	{
		"name":          "JohnDoe",
		"gender":        "Male",
		"age":           20,
		"address":       "BoltonUK",
		"hasDL":         true,
		"ssn":           "1231231234",
		"income":        45000,
		"maritalStatus": "single",
		"creditScore":   500,
	},
	{
		"name":          "JaneDoe",
		"gender":        "Female",
		"age":           38,
		"address":       "BoltonUK",
		"hasDL":         false,
		"ssn":           "2424354532",
		"income":        32000,
		"maritalStatus": "single",
		"creditScore":   650,
	},
	{
		"name":          "PrakashY",
		"gender":        "Male",
		"age":           30,
		"address":       "RedwoodShore",
		"hasDL":         true,
		"ssn":           "2345342132",
		"income":        150000,
		"maritalStatus": "married",
		"creditScore":   750,
	},
	{
		"name":          "SandraW",
		"gender":        "Female",
		"age":           26,
		"address":       "RedwoodShore",
		"hasDL":         true,
		"ssn":           "3213214321",
		"income":        50000,
		"maritalStatus": "single",
		"creditScore":   625,
	},
}

func TestCellCompileExpr(t *testing.T) {
	err := model.RegisterTupleDescriptors(string(tupleDescriptor))
	assert.Nil(t, err)

	// test cases input
	tupleType := "applicant"
	propName := "name"
	testcases := make(map[string]string)
	testcases["foo"] = "$.applicant.name == foo"
	testcases["123foo"] = "$.applicant.name == 123foo"
	testcases["foo123"] = "$.applicant.name == foo123"
	testcases["123"] = "$.applicant.name == 123"
	testcases["123.123"] = "$.applicant.name == 123.123"
	testcases[".123"] = "$.applicant.name == .123"
	testcases["!foo"] = "!($.applicant.name == foo)"
	testcases["==foo"] = "$.applicant.name == foo"
	testcases["!=foo"] = "$.applicant.name != foo"
	testcases[">foo"] = "$.applicant.name > foo"
	testcases["!>foo"] = "!($.applicant.name > foo)"
	testcases[">=foo"] = "$.applicant.name >= foo"
	testcases["<foo"] = "$.applicant.name < foo"
	testcases["<=foo"] = "$.applicant.name <= foo"
	testcases["< foo"] = "$.applicant.name < foo" // space test
	testcases["foo&&bar"] = "($.applicant.name == foo && $.applicant.name == bar)"
	testcases[">=foo&&bar"] = "($.applicant.name >= foo && $.applicant.name == bar)"
	testcases["car&&jeep&&bus"] = "(($.applicant.name == car && $.applicant.name == jeep) && $.applicant.name == bus)"
	testcases["foo||bar"] = "($.applicant.name == foo || $.applicant.name == bar)"
	testcases["car||jeep||bus"] = "(($.applicant.name == car || $.applicant.name == jeep) || $.applicant.name == bus)"
	testcases["car&&jeep||bus"] = "(($.applicant.name == car && $.applicant.name == jeep) || $.applicant.name == bus)"
	testcases["!(car&&(jeep||bus))"] = "!(($.applicant.name == car && ($.applicant.name == jeep || $.applicant.name == bus)))"
	testcases["car||jeep&&bus"] = "($.applicant.name == car || ($.applicant.name == jeep && $.applicant.name == bus))"

	// prepare cell
	tupleDesc := model.GetTupleDescriptor(model.TupleType(tupleType))
	assert.NotNil(t, tupleDesc)
	propDesc := tupleDesc.GetProperty(propName)
	assert.NotNil(t, propDesc)
	cell := &genCell{
		metaCell: &metaCell{
			colType:   ctCondition,
			tupleDesc: tupleDesc,
			propDesc:  propDesc,
		},
	}

	// run test cases
	for k, v := range testcases {
		t.Log(k, v)
		cell.rawValue = k
		cell.compileExpr()
		assert.Equal(t, v, cell.cdExpr)
	}
}

func TestNew(t *testing.T) {
	err := model.RegisterTupleDescriptors(string(tupleDescriptor))

	settings := &Settings{
		DTableFile: "test_dtable.csv",
	}
	mf := mapper.NewFactory(resolve.GetBasicResolver())
	initCtx := test.NewActivityInitContext(settings, mf)
	act, err := New(initCtx)
	assert.Nil(t, err)
	assert.NotNil(t, act)
}

func TestEval(t *testing.T) {
	err := model.RegisterTupleDescriptors(string(tupleDescriptor))

	settings := &Settings{
		DTableFile: "test_dtable.xlsx",
	}
	mf := mapper.NewFactory(resolve.GetBasicResolver())
	initCtx := test.NewActivityInitContext(settings, mf)
	act, err := New(initCtx)
	assert.Nil(t, err)
	assert.NotNil(t, act)

	tuples := make(map[model.TupleType]model.Tuple)
	tuple, err := model.NewTuple(model.TupleType("applicant"), testApplicants[0])
	assert.Nil(t, err)
	assert.NotNil(t, tuple)
	tuples[tuple.GetTupleType()] = tuple

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("ctx", context.Background())
	tc.SetInput("tuples", tuples)
	act.Eval(tc)

	expectedStatus, err := tuple.GetString("status")
	assert.Nil(t, err)
	assert.Equal(t, "VISA-Granted", expectedStatus)
	expectedEligible, err := tuple.GetBool("eligible")
	assert.Nil(t, err)
	assert.Equal(t, true, expectedEligible)
	expectedCreditLimit, err := tuple.GetDouble("creditLimit")
	assert.Nil(t, err)
	assert.Equal(t, 2500.0, expectedCreditLimit)
}
