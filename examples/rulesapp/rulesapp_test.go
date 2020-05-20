package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuleApp(t *testing.T) {
	err := example()
	assert.Nil(t, err)
}
