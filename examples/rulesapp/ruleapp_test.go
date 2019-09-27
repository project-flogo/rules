package main

import (
	"strings"
	"testing"

	"github.com/project-flogo/rules/ruleapi/tests"
	"github.com/stretchr/testify/assert"
)

func TestRuleApp(t *testing.T) {
	request := func() {
		tests.Command("go", "run", "main.go")
	}
	output := tests.CaptureStdOutput(request)
	var result string
	if strings.Contains(output, "Rule fired") && strings.Contains(output, "Loaded tuple descriptor") {
		result = "success"
	}
	assert.Equal(t, "success", result)
}
