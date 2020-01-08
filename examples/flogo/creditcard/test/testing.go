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
		response, err := client.Get("http://localhost:7777/test/applicant?name=JohnDoe&gender=Male&age=20&address=BoltonUK&hasDL=false&ssn=1231231234&income=45000&maritalStatus=single&creditScore=500")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output := tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Saving tuple.")

	request = func() {
		response, err := client.Get("http://localhost:7777/test/applicant?name=JaneDoe&gender=Female&age=38&address=BoltonUK&hasDL=false&ssn=2424354532&income=32000&maritalStatus=single&creditScore=650")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output = tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Saving tuple.")

	request = func() {
		response, err := client.Get("http://localhost:7777/test/applicant?name=PrakashY&gender=Male&age=30&address=RedwoodShore&hasDL=true&ssn=2345342132&income=150000&maritalStatus=married&creditScore=750")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output = tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Saving tuple.")

	request = func() {
		response, err := client.Get("http://localhost:7777/test/applicant?name=SandraW&gender=Female&age=26&address=RedwoodShore&hasDL=true&ssn=3213214321&income=50000&maritalStatus=single&creditScore=625")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output = tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Saving tuple.")

	request = func() {
		response, err := client.Get("http://localhost:7777/test/processapplication?start=true&ssn=1231231234")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output = tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Saving tuple.")

	request = func() {
		response, err := client.Get("http://localhost:7777/test/processapplication?start=true&ssn=2345342132")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output = tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Saving tuple.")

	request = func() {
		response, err := client.Get("http://localhost:7777/test/processapplication?start=true&ssn=3213214321")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output = tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Saving tuple.")

	request = func() {
		response, err := client.Get("http://localhost:7777/test/processapplication?start=true&ssn=2424354532")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)
	}
	output = tests.CaptureStdOutput(request)
	assert.Contains(t, output, "Saving tuple.")
}
