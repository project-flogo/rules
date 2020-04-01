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

	request := func() {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:7777/moveevent?packageid=PACKAGE1&targetstate=sitting", nil)
		assert.Nil(t, err)
		response, err := client.Do(req)
		assert.Nil(t, err)
		_, err = ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
	}

	outpt := tests.CaptureStdOutput(request)
	var result string
	if strings.Contains(outpt, "target state [sitting]") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	result = ""

	request1 := func() {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:7777/moveevent?packageid=PACKAGE1&targetstate=sitting", nil)
		assert.Nil(t, err)
		response, err := client.Do(req)
		assert.Nil(t, err)
		_, err = ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
	}

	outpt = tests.CaptureStdOutput(request1)
	if strings.Contains(outpt, "target state [sitting]") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	result = ""

	request2 := func() {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:7777/moveevent?packageid=PACKAGE1&targetstate=moving", nil)
		assert.Nil(t, err)
		response, err := client.Do(req)
		assert.Nil(t, err)
		_, err = ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
	}

	outpt = tests.CaptureStdOutput(request2)
	if strings.Contains(outpt, "target state [moving]") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	result = ""

	request3 := func() {
		req, err := http.NewRequest(http.MethodGet, "http://localhost:7777/moveevent?packageid=PACKAGE2&targetstate=normal", nil)
		assert.Nil(t, err)
		response, err := client.Do(req)
		assert.Nil(t, err)
		_, err = ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
	}

	outpt = tests.CaptureStdOutput(request3)
	if strings.Contains(outpt, "Tuple inserted successfully") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	result = ""
}

func TestStateMachineJSON(t *testing.T) {
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
