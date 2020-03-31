package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/project-flogo/rules/ruleapi/tests"
	"github.com/stretchr/testify/assert"
)

var redis = false

func TestMain(m *testing.M) {
	code := m.Run()
	if code != 0 {
		os.Exit(code)
	}

	run := func() int {
		command := exec.Command("docker", "run", "-p", "6383:6379", "-d", "redis")
		hash, err := command.Output()
		if err != nil {
			panic(err)
		}
		tests.Pour("6383")

		defer func() {
			command := exec.Command("docker", "stop", strings.TrimSpace(string(hash)))
			err := command.Run()
			if err != nil {
				panic(err)
			}
			command = exec.Command("docker", "rm", strings.TrimSpace(string(hash)))
			err = command.Run()
			if err != nil {
				panic(err)
			}
			tests.Drain("6383")
		}()

		return m.Run()
	}
	redis = true
	os.Exit(run())
}

func TestRuleApp(t *testing.T) {
	os.Setenv("name", "Smith")
	err := example(redis)
	assert.Nil(t, err)
}
