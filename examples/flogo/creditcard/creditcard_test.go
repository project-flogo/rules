package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	_ "github.com/project-flogo/core/data/expression/script"
	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/rules/ruleapi/tests"

	"github.com/stretchr/testify/assert"
)

func TestCreditCard(t *testing.T) {

	data, err := ioutil.ReadFile(filepath.FromSlash("flogo.json"))
	assert.Nil(t, err)
	cfg, err := engine.LoadAppConfig(string(data), false)
	assert.Nil(t, err)
	e, err := engine.New(cfg)
	assert.Nil(t, err)
	tests.Drain("7777")
	err = e.Start()
	assert.Nil(t, err)
	defer func() {
		err := e.Stop()
		assert.Nil(t, err)
	}()
	tests.Pour("7777")

	client := &http.Client{}
	request := func() {
		req, err := http.NewRequest("PUT", "http://localhost:7777/newaccount", bytes.NewBuffer([]byte(`{"Name":"Test","Age":"26","Income":"60100","Address":"TEt","Id":"12312","Gender":"male","maritalStatus":"single"}`)))
		req.Header.Set("Content-Type", "application/json")
		response, err := client.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output := tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Rule fired: NewUser")

	request = func() {
		req, err := http.NewRequest("PUT", "http://localhost:7777/credit", bytes.NewBuffer([]byte(`{"Id":12312,"creditScore":680}`)))
		req.Header.Set("Content-Type", "application/json")
		response, err := client.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output = tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Rule fired: Rejected")

}
