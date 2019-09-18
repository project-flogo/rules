package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/rules/ruleapi/tests"
	"github.com/stretchr/testify/assert"
)

// Response is a reply form the server
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

var (
	payload = `{
		"Name": "Sam",
		"Age": "26",
		"Income": "50100",
		"Address": "SFO",
		"Id": "4"
	  }`
	payload1 = `{
		"Id": "4",
		"creditScore": "850"
	  }`
	payload2 = `{
		"Name": "Sam1",
		"Age": "17",
		"Income": "50100",
		"Address": "SFO",
		"Id": "5"
	  }`
	payload3 = `{
		"Name": "Sam2",
		"Age": "26",
		"Income": "5100",
		"Address": "SFO",
		"Id": "6"
	  }`
	payload4 = `{
		"Name": "Sam3",
		"Age": "26",
		"Income": "5100",
		"Address": "",
		"Id": "7"
	  }`
	payload5 = `{
		"Name": "Sam5",
		"Age": "32",
		"Income": "75000",
		"Address": "SFO",
		"Id": "8"
	  }`
	payload6 = `{
		"Id": "8",
		"creditScore": "760"
	  }`
	payload7 = `{
		"Name": "Sam4",
		"Age": "41",
		"Income": "30000",
		"Address": "SFO",
		"Id": "9"
	  }`
	payload8 = `{
		"Id": "9",
		"creditScore": "720"
	  }`
)

func testapplication(t *testing.T, e engine.Engine) {
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

	payload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	payload1, err := json.Marshal(payload1)
	if err != nil {
		panic(err)
	}

	payload2, err := json.Marshal(payload2)
	if err != nil {
		panic(err)
	}

	payload3, err := json.Marshal(payload3)
	if err != nil {
		panic(err)
	}

	payload4, err := json.Marshal(payload4)
	if err != nil {
		panic(err)
	}

	payload5, err := json.Marshal(payload5)
	if err != nil {
		panic(err)
	}

	payload6, err := json.Marshal(payload6)
	if err != nil {
		panic(err)
	}

	payload7, err := json.Marshal(payload7)
	if err != nil {
		panic(err)
	}

	payload8, err := json.Marshal(payload8)
	if err != nil {
		panic(err)
	}

	//  valid new user details
	request := func() {
		req, err := http.NewRequest(http.MethodPut, "http://localhost:7777/newaccount", bytes.NewBuffer(payload))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		response, err := client.Do(req)
		assert.Nil(t, err)
		_, err = ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
		assert.Nil(t, err)
	}
	outpt := tests.CaptureStdOutput(request)
	var result string
	if strings.Contains(outpt, "Rule fired") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	outpt = ""
	result = ""

	//  user with credit score > 850
	request1 := func() {
		req, err := http.NewRequest(http.MethodPut, "http://localhost:7777/credit", bytes.NewBuffer(payload1))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		response, err := client.Do(req)
		assert.Nil(t, err)
		_, err = ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
		assert.Nil(t, err)
	}
	outpt = tests.CaptureStdOutput(request1)
	if strings.Contains(outpt, "Rule fired") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	outpt = ""
	result = ""

	//  New user with age < 17
	request2 := func() {
		req, err := http.NewRequest(http.MethodPut, "http://localhost:7777/newaccount", bytes.NewBuffer(payload2))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		response, err := client.Do(req)
		assert.Nil(t, err)
		_, err = ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
		assert.Nil(t, err)
	}
	outpt = tests.CaptureStdOutput(request2)
	if strings.Contains(outpt, "Applicant is not eligible to apply for creditcard") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	outpt = ""
	result = ""

	//  New user with income < 10k
	request3 := func() {
		req, err := http.NewRequest(http.MethodPut, "http://localhost:7777/newaccount", bytes.NewBuffer(payload3))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		response, err := client.Do(req)
		assert.Nil(t, err)
		_, err = ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
		assert.Nil(t, err)
	}
	outpt = tests.CaptureStdOutput(request3)
	if strings.Contains(outpt, "Applicant is not eligible to apply for creditcard") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	outpt = ""
	result = ""

	//  New user with empty address
	request4 := func() {
		req, err := http.NewRequest(http.MethodPut, "http://localhost:7777/newaccount", bytes.NewBuffer(payload4))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		response, err := client.Do(req)
		assert.Nil(t, err)
		_, err = ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
		assert.Nil(t, err)
	}
	outpt = tests.CaptureStdOutput(request4)
	if strings.Contains(outpt, "Applicant is not eligible to apply for creditcard") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	outpt = ""
	result = ""

	request5 := func() {
		req, err := http.NewRequest(http.MethodPut, "http://localhost:7777/newaccount", bytes.NewBuffer(payload5))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		response, err := client.Do(req)
		assert.Nil(t, err)
		_, err = ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
		assert.Nil(t, err)
	}
	outpt = tests.CaptureStdOutput(request5)
	if strings.Contains(outpt, "Rule fired") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	outpt = ""
	result = ""

	//  user with credit score 760
	request6 := func() {
		req, err := http.NewRequest(http.MethodPut, "http://localhost:7777/credit", bytes.NewBuffer(payload6))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		response, err := client.Do(req)
		assert.Nil(t, err)
		_, err = ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
		assert.Nil(t, err)
	}
	outpt = tests.CaptureStdOutput(request6)
	if strings.Contains(outpt, "Rule fired") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	outpt = ""
	result = ""

	request7 := func() {
		req, err := http.NewRequest(http.MethodPut, "http://localhost:7777/newaccount", bytes.NewBuffer(payload7))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		response, err := client.Do(req)
		assert.Nil(t, err)
		_, err = ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
		assert.Nil(t, err)
	}
	outpt = tests.CaptureStdOutput(request7)
	if strings.Contains(outpt, "Rule fired") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	outpt = ""
	result = ""

	//  New user with credit score < 720
	request8 := func() {
		req, err := http.NewRequest(http.MethodPut, "http://localhost:7777/credit", bytes.NewBuffer(payload8))
		assert.Nil(t, err)
		req.Header.Set("Content-Type", "application/json")
		response, err := client.Do(req)
		assert.Nil(t, err)
		_, err = ioutil.ReadAll(response.Body)
		assert.Nil(t, err)
		response.Body.Close()
		assert.Nil(t, err)
	}
	outpt = tests.CaptureStdOutput(request8)
	if strings.Contains(outpt, "Rule fired: Rejected") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	outpt = ""
	result = ""
}

func TestCreditCardJSON(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Handler Routing JSON integration test in short mode")
	}

	data, err := ioutil.ReadFile(filepath.FromSlash("./flogo.json"))
	assert.Nil(t, err)
	cfg, err := engine.LoadAppConfig(string(data), false)
	assert.Nil(t, err)
	e, err := engine.New(cfg)
	assert.Nil(t, err)
	testapplication(t, e)
}
