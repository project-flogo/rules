package main

import (
	"flag"
	"fmt"

	"github.com/project-flogo/rules/examples/ordermanagement/audittrail"
	"github.com/project-flogo/rules/ruleapi"

	"github.com/project-flogo/core/support/log"
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

	log.RootLogger().Info("***** Order Management App *****")

	log.RootLogger().Info("--- Loading Tuple Schema ---")
	err = loadTupleSchema()
	if err != nil {
		panic(err)
	}

	log.RootLogger().Info("--- Adding Rules ---")
	ruleSession, err = createAndLoadRuleSession()
	if err != nil {
		panic(err)
	}

	log.RootLogger().Info("--- Starting Kinesis Publisher ---")
	audittrail.SetupKinesisPubSub(config, *awsStreamName)

	log.RootLogger().Info("--- Starting WebSocket Publisher ---")
	wsPublisher := audittrail.Create(8686, *awsStreamName)
	go wsPublisher.Start()

	log.RootLogger().Info("--- Starting Rules Engine ---")
	ruleSession.Start(nil)

	log.RootLogger().Info("--- Starting MQTT Subscriber ---")
	setupFlogoMQTTTriggers()
}

// loads the tuple schema
func loadTupleSchema() error {
	content := getFileContent(tupleSchemaPath)
	err := model.RegisterTupleDescriptors(string(content))
	if err != nil {
		return err
	}
	return nil
}

// create rulesession and load rules in it
func createAndLoadRuleSession() (model.RuleSession, error) {
	content := getFileContent(ruleDefinitionPath)
	return ruleapi.GetOrCreateRuleSessionFromConfig("oms_session", string(content))
}

// Get file content
func getFileContent(filePath string) string {
	absPath := common.GetAbsPathForResource(filePath)
	return common.FileToString(absPath)
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
