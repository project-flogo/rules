package main

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/rules/ruleapi/tests"
	"github.com/stretchr/testify/assert"
)

func testApplication(t *testing.T, e engine.Engine) {
	tests.Drain("7777")
	err := e.Start()
	assert.Nil(t, err)
	defer func() {
		err := e.Stop()
		assert.Nil(t, err)
	}()
	tests.Pour("7777")

	transport := &http.Transport{
		MaxIdleConns: 1,
	}
	defer transport.CloseIdleConnections()
	client := &http.Client{
		Transport: transport,
	}

	// check for samename condition
	request := func() {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:7777/test/n1?name=Bob", nil)
		assert.Nil(t, err)
		response, err := client.Do(req)
		assert.Nil(t, err)
		_, err = ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
	}

	outpt := tests.CaptureStdOutput(request)

	var result string
	if strings.Contains(outpt, "Rule fired") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	result = ""

	// check for tuples match n1 and n2
	request1 := func() {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:7777/test/n2?name=Bob", nil)
		assert.Nil(t, err)
		response, err := client.Do(req)
		assert.Nil(t, err)
		_, err = ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
	}

	outpt = tests.CaptureStdOutput(request1)
	if strings.Contains(outpt, "Rule fired") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	result = ""

	//  check  for name mismatch
	request2 := func() {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:7777/test/n1?name=Tom", nil)
		assert.Nil(t, err)
		response, err := client.Do(req)
		assert.Nil(t, err)
		_, err = ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
	}

	outpt = tests.CaptureStdOutput(request2)
	if !strings.Contains(outpt, "Rule fired") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	result = ""

	//  Already asserted tuple check
	request3 := func() string {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:7777/test/n1?name=Bob", nil)
		assert.Nil(t, err)
		response, err := client.Do(req)
		assert.Nil(t, err)
		body, err := ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
		return string(body)
	}

	body := request3()
	if strings.Contains(body, "already asserted") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	result = ""
}

func TestSimpleJSON(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Handler Routing JSON integration test in short mode")
	}

	data, err := ioutil.ReadFile(filepath.FromSlash("./flogo.json"))
	assert.Nil(t, err)
	cfg, err := engine.LoadAppConfig(string(data), false)
	assert.Nil(t, err)
	e, err := engine.New(cfg)
	assert.Nil(t, err)
	testApplication(t, e)
}
