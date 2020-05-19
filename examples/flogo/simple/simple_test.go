package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/project-flogo/core/data/expression/script"
	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/rules/ruleapi/tests"

	"github.com/stretchr/testify/assert"
)

func TestSimpleApp(t *testing.T) {

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
		response, err := client.Get("http://localhost:7777//test/n1?name=Bob")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output := tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Rule fired: [n1.name == Bob]")

	request = func() {
		response, err := client.Get("http://localhost:7777//test/n2?name=Bob")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output = tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Rule fired: [n1.name == Bob && n1.name == n2.name]")

	request = func() {
		response, err := client.Get("http://localhost:7777//test/n1?name=testprop")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output = tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Rule fired: [flogo property example]")

	// set env variable used for the testcase
	defaultVal := os.Getenv("simpleappname")
	os.Setenv("simpleappname", "test1234")
	defer func() {
		os.Setenv("simpleappname", defaultVal)
	}()

	request = func() {
		response, err := client.Get("http://localhost:7777//test/n1?name=test1234")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output = tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Rule fired: [env variable example]")

}
