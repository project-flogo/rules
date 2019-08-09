package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/project-flogo/core/engine"
	"github.com/project-flogo/rules/ruleapi/tests"
	"github.com/stretchr/testify/assert"
)

const (
	kafkaConn = "localhost:9092"
	topic     = "orderinfo"
)

func initProducer() (sarama.SyncProducer, error) {
	// setup sarama log to stdout
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	// producer config
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	// sync producer
	prd, err := sarama.NewSyncProducer([]string{kafkaConn}, config)

	return prd, err
}

func publish(message string, producer sarama.SyncProducer) {
	// publish sync
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	p, o, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Error publish: ", err.Error())
	}

	fmt.Println("Partition: ", p)
	fmt.Println("Offset: ", o)
}

func testApplication(t *testing.T, e engine.Engine) {
	err := e.Start()
	assert.Nil(t, err)
	defer func() {
		tests.Command("docker", "rm", "-f", "$(docker ps -aq)")
		err := e.Stop()
		assert.Nil(t, err)
	}()
	tests.Pour("9092")

	producer, err := initProducer()
	if err != nil {
		fmt.Println("Error producer: ", err.Error())
		os.Exit(1)
	}

	request := func() {
		publish(`{"type":"grocery","totalPrice":"2001.0"}`, producer)
	}
	outpt := tests.CaptureOutput(request)

	fmt.Println("********************", outpt, "888888888888888888888")
	var result string
	if strings.Contains(outpt, "Rule fired") {
		result = "success"
	}
	assert.Equal(t, "success", result)
	outpt = ""
	result = ""

}

func TestSimpleKafkaJSON(t *testing.T) {
	cmd := exec.Command("docker-compose")
	err := cmd.Run()
	if err == nil {
		t.Skip("skipping test - docker-compose not found")
	}

	data, err := ioutil.ReadFile(filepath.FromSlash("./flogo.json"))
	assert.Nil(t, err)
	tests.Command("docker-compose", "up", "-d")
	time.Sleep(50 * time.Second)
	cfg, err := engine.LoadAppConfig(string(data), false)
	assert.Nil(t, err)
	e, err := engine.New(cfg)
	assert.Nil(t, err)
	testApplication(t, e)
}
