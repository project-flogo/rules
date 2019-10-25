package tests

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	if code != 0 {
		os.Exit(code)
	}

	run := func() int {
		command := exec.Command("docker", "run", "-p", "6380:6379", "-d", "redis")
		hash, err := command.Output()
		if err != nil {
			panic(err)
		}
		Pour("6380")

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
			Drain("6380")
		}()

		return m.Run()
	}
	redis = true
	code = run()
	if code != 0 {
		os.Exit(code)
	}

	performance = true
	os.Exit(run())
}
