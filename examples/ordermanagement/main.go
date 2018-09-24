package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/project-flogo/rules/examples/ordermanagement/audittrail"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/project-flogo/rules/common"
	"github.com/project-flogo/rules/common/model"
)

//go:generate go run $GOPATH/src/github.com/TIBCOSoftware/flogo-lib/flogo/gen/gen.go $GOPATH

const (
	msgValueField      = "message"
	tupleSchemaPath    = "src/github.com/project-flogo/rules/examples/ordermanagement/schema/oms_schema.json"
	ruleDefinitionPath = "src/github.com/project-flogo/rules/examples/ordermanagement/schema/rule_definition.json"
)

// cli arguments
var (
	awsAccesskey  = flag.String("accesskey", "", "Access Key for AWS profile")
	awsSecretkey  = flag.String("secretkey", "", "Secret Key for AWS profile")
	awsRegion     = flag.String("region", "", "AWS region")
	awsStreamName = flag.String("streamName", "oms_audittrail", "AWS kinesis stream name")
)

// RuleSession to maintain all the rule definitions
var ruleSession model.RuleSession

func main() {
	flag.Parse()

	err := validate()
	if err != nil {
		panic(err)
	}

	config := audittrail.ConnectionConfig{
		AccessKey:  *awsAccesskey,
		SecretKey:  *awsSecretkey,
		RegionName: *awsRegion,
	}

	logger.Info("***** Order Management App *****")

	logger.Info("--- Loading Tuple Schema ---")
	err = loadTupleSchema()
	if err != nil {
		panic(err)
	}

	logger.Info("--- Adding Rules ---")
	ruleSession, err = loadRulesInRuleSession()
	if err != nil {
		panic(err)
	}

	logger.Info("--- Starting Kinesis Publisher ---")
	audittrail.SetupKinesisPubSub(config, *awsStreamName)

	logger.Info("--- Starting WebSocket Publisher ---")
	wsPublisher := audittrail.Create(8686, *awsStreamName)
	go wsPublisher.Start()

	logger.Info("--- Starting Rules Engine ---")
	ruleSession.Start(nil)

	logger.Info("--- Starting MQTT Subscriber ---")
	setupFlogoMQTTTriggers()
}

// loads the tuple schema
func loadTupleSchema() error {
	tupleDescFileAbsPath := common.GetAbsPathForResource(tupleSchemaPath)

	dat, err := ioutil.ReadFile(tupleDescFileAbsPath)
	if err != nil {
		panic(err)
	}
	err = model.RegisterTupleDescriptors(string(dat))
	if err != nil {
		return err
	}
	return nil
}

// validates for any missing variables
func validate() error {
	if *awsAccesskey == "" || *awsSecretkey == "" || *awsRegion == "" {
		return fmt.Errorf("One or more of the required AWS credential details are missing (AccessKey/SecretKey/Region)")
	}

	if *awsStreamName == "" {
		return fmt.Errorf("Missing Kinesis stream name")
	}
	return nil
}
