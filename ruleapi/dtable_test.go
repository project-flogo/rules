package ruleapi

import (
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/stretchr/testify/assert"
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
			ColType:   ctCondition,
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

func TestDtableXLSX(t *testing.T) {
	err := model.RegisterTupleDescriptors(string(tupleDescriptor))
	assert.Nil(t, err)

	dtable, err := LoadDecisionTableFromFile("test_dtable.xlsx")
	assert.Nil(t, err)

	err = dtable.Compile()
	assert.Nil(t, err)

	vls := make(map[string]interface{})

	vls["name"] = "test"
	vls["age"] = 20
	vls["hasDL"] = true
	vls["creditScore"] = 600
	vls["maritalStatus"] = "single"
	vls["income"] = 11000
	vls["eligible"] = false
	vls["creditLimit"] = 0

	tpl, err := model.NewTuple("applicant", vls)
	assert.Nil(t, err)
	assert.NotNil(t, tpl)

	presentCreditLimit, err := tpl.GetDouble("creditLimit")
	assert.Nil(t, err)
	assert.Equal(t, float64(0), presentCreditLimit)

	tuples := make(map[model.TupleType]model.Tuple)
	tuples[tpl.GetTupleType()] = tpl

	dtable.Apply(nil, tuples)

	newCreditLimit, err := tpl.GetDouble("creditLimit")
	assert.Nil(t, err)
	assert.NotEqual(t, newCreditLimit, presentCreditLimit)

}

func TestDtableCSV(t *testing.T) {
	err := model.RegisterTupleDescriptors(string(tupleDescriptor))
	assert.Nil(t, err)

	dtable, err := LoadDecisionTableFromFile("test_dtable.csv")
	assert.Nil(t, err)

	err = dtable.Compile()
	assert.Nil(t, err)

	vls := make(map[string]interface{})

	vls["name"] = "test"
	vls["age"] = 20
	vls["hasDL"] = true
	vls["creditScore"] = 600
	vls["maritalStatus"] = "single"
	vls["income"] = 11000
	vls["eligible"] = false
	vls["creditLimit"] = 0

	tpl, err := model.NewTuple("applicant", vls)
	assert.Nil(t, err)
	assert.NotNil(t, tpl)

	presentCreditLimit, err := tpl.GetDouble("creditLimit")
	assert.Nil(t, err)
	assert.Equal(t, float64(0), presentCreditLimit)

	tuples := make(map[model.TupleType]model.Tuple)
	tuples[tpl.GetTupleType()] = tpl

	dtable.Apply(nil, tuples)

	newCreditLimit, err := tpl.GetDouble("creditLimit")
	assert.Nil(t, err)
	assert.NotEqual(t, newCreditLimit, presentCreditLimit)

}
