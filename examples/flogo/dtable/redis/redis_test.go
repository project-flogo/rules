package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/project-flogo/rules/examples/flogo/dtable/test"
	"github.com/project-flogo/rules/ruleapi/tests"

	"github.com/stretchr/testify/assert"
)

func TestDTableRedis(t *testing.T) {
	command := exec.Command("docker", "run", "-p", "6381:6379", "-d", "redis")
	hash, err := command.Output()
	if err != nil {
		assert.Nil(t, err)
	}
	tests.Pour("6381")

	defer func() {
		command := exec.Command("docker", "stop", strings.TrimSpace(string(hash)))
		err := command.Run()
		if err != nil {
			assert.Nil(t, err)
		}
		command = exec.Command("docker", "rm", strings.TrimSpace(string(hash)))
		err = command.Run()
		if err != nil {
			assert.Nil(t, err)
		}
		tests.Drain("6381")
	}()
	os.Setenv("STORECONFIG", filepath.FromSlash("../rsconfig.json"))
	test.Test(t)
}
