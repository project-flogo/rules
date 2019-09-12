package tests

import (
	"testing"

	"github.com/project-flogo/rules/ruleapi"
)

func TestClearSessions(t *testing.T) {
	ruleapi.GetOrCreateRuleSession("test")
	ruleapi.ClearSessions()
}
