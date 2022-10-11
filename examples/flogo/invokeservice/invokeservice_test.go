package main

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/rules/ruleapi/tests"

	"github.com/stretchr/testify/assert"
)

func TestInvokeService(t *testing.T) {

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
		response, err := client.Get("http://localhost:7777/test/n1?name=function")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output := tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Rule[n1.name == function] fired. serviceFunctionAction() function got invoked.")

	request = func() {
		response, err := client.Get("http://localhost:7777/test/n1?name=activity")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output = tests.CaptureStdOutput(request)
	assert.Contains(t, output, "")

}
