package test

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	_ "github.com/project-flogo/core/data/expression/script"
	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/rules/ruleapi/tests"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	data, err := ioutil.ReadFile(filepath.FromSlash("../flogo.json"))
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
		response, err := client.Get("http://localhost:7777/test/student?grade=GRADE-C&name=s1&class=X-A&careRequired=false")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output := tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Saving tuple.")

	request = func() {
		response, err := client.Get("http://localhost:7777/test/student?grade=GRADE-B&name=s2&class=X-A&careRequired=false")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output = tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Saving tuple.")

	request = func() {
		response, err := client.Get("http://localhost:7777/test/studentanalysis?name=s1")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output = tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Saving tuple.")

	request = func() {
		response, err := client.Get("http://localhost:7777/test/studentanalysis?name=s2")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output = tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Saving tuple.")
}
