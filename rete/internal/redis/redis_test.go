package redis

import (
	"context"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/common"
	"github.com/project-flogo/rules/rete/internal/types"
)

func Drain(port string) {
	for {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("", port), time.Second)
		if conn != nil {
			conn.Close()
		}
		if err != nil && strings.Contains(err.Error(), "connect: connection refused") {
			break
		}
	}
}

func Pour(port string) {
	for {
		conn, _ := net.Dial("tcp", net.JoinHostPort("", port))
		if conn != nil {
			conn.Close()
			break
		}
	}
}

func TestMain(m *testing.M) {
	run := func() int {
		command := exec.Command("docker", "run", "-p", "6384:6379", "-d", "redis")
		hash, err := command.Output()
		if err != nil {
			panic(err)
		}
		Pour("6384")

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
			Drain("6384")
		}()

		return m.Run()
	}
	os.Exit(run())
}

type testNetwork struct {
	Prefix string
}

func (n *testNetwork) AddRule(model.Rule) error {
	return nil
}

func (n *testNetwork) String() string {
	return ""
}

func (n *testNetwork) RemoveRule(string) model.Rule {
	return nil
}

func (n *testNetwork) GetRules() []model.Rule {
	return nil
}

func (n *testNetwork) Assert(ctx context.Context, rs model.RuleSession, tuple model.Tuple, changedProps map[string]bool, mode common.RtcOprn) error {
	return nil
}

func (n *testNetwork) Retract(ctx context.Context, rs model.RuleSession, tuple model.Tuple, changedProps map[string]bool, mode common.RtcOprn) error {
	return nil
}

func (n *testNetwork) GetAssertedTuple(ctx context.Context, rs model.RuleSession, key model.TupleKey) model.Tuple {
	return nil
}

func (n *testNetwork) RegisterRtcTransactionHandler(txnHandler model.RtcTransactionHandler, txnContext interface{}) {

}

func (n *testNetwork) SetTupleStore(tupleStore model.TupleStore) {

}

func (n *testNetwork) GetHandleWithTuple(ctx context.Context, tuple model.Tuple) types.ReteHandle {
	return nil
}

func (n *testNetwork) AssertInternal(ctx context.Context, tuple model.Tuple, changedProps map[string]bool, mode common.RtcOprn) error {
	return nil
}

func (n *testNetwork) RetractInternal(ctx context.Context, tuple model.Tuple, changedProps map[string]bool, mode common.RtcOprn) error {
	return nil
}

func (n *testNetwork) GetPrefix() string {
	return n.Prefix
}

func (n *testNetwork) GetIdGenService() types.IdGen {
	return nil
}

func (n *testNetwork) GetLockService() types.LockService {
	return nil
}

func (n *testNetwork) GetJtService() types.JtService {
	return nil
}

func (n *testNetwork) GetHandleService() types.HandleService {
	return nil
}

func (n *testNetwork) GetJtRefService() types.JtRefsService {
	return nil
}

func (n *testNetwork) GetTupleStore() model.TupleStore {
	return nil
}

func TestLockServiceImpl(t *testing.T) {
	fini := make(chan bool, 1)
	go func() {
		config := common.Config{
			IDGens: common.Service{
				Redis: redisutils.RedisConfig{
					Network: "tcp",
					Address: ":6384",
				},
			},
		}
		network := &testNetwork{
			Prefix: "a",
		}
		serviceA, serviceB := NewLockServiceImpl(network, config), NewLockServiceImpl(network, config)
		done := make(chan bool, 8)
		for i := 0; i < 100; i++ {
			go func() {
				serviceA.Lock()
				time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
				serviceA.Unlock()
				done <- true
			}()
			go func() {
				serviceB.Lock()
				time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
				serviceB.Unlock()
				done <- true
			}()
		}
		for i := 0; i < 100; i++ {
			<-done
			<-done
		}
		fini <- true
	}()
	select {
	case <-time.After(60 * time.Second):
		t.Fatal("test took too long")
	case <-fini:
	}
}
