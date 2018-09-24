package audittrail

import (
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
)

// KPInstance is an instanceof Kinesis publisher
var KPInstance KinesisPublisher

// Entry to record the trail of flow
type Entry struct {
	OrderID     string
	TimeStamp   time.Time
	Status      string
	RuleName    string
	Description string
}

// ConnectionConfig stores the local connection details
type ConnectionConfig struct {
	AccessKey  string
	SecretKey  string
	RegionName string
}

// KinesisPublisher interface providing basic communication methods
type KinesisPublisher interface {
	Connect()
	CreateStream(streamName string) error
	Publish(streamName string, content interface{}) (bool, error)
	GetRecords(streamName string) error
	Disconnect() error
	StreamExists(streamName string) (bool, error)
	DescribeStream(streamName string) (*kinesis.DescribeStreamOutput, error)
}

const (
	// Since a single shard is used, a single static key
	partitionKey = "key1"
)

// Kinesis Publisher implementation to publish/get data to/from kinesis
type kinesisPublisherImpl struct {
	config        ConnectionConfig
	kinesisClient *kinesis.Kinesis
}

// New creates a new kinesis publisher instance
func New(config ConnectionConfig) (KinesisPublisher, error) {
	return &kinesisPublisherImpl{config: config}, nil
}

func (kp *kinesisPublisherImpl) CreateStream(streamName string) error {
	_, err := kp.kinesisClient.CreateStream(&kinesis.CreateStreamInput{
		ShardCount: aws.Int64(1),
		StreamName: aws.String(streamName),
	})
	if err != nil {
		panic(err)
	}

	if err := kp.kinesisClient.WaitUntilStreamExists(&kinesis.DescribeStreamInput{StreamName: aws.String(streamName)}); err != nil {
		panic(err)
	}
	return nil
}

func (kp *kinesisPublisherImpl) Publish(streamName string, content interface{}) (bool, error) {
	auditTrailItem, ok := content.(Entry)
	if !ok {
		log.Fatal("Invalid content type, expected content type AuditTrailItem")
	}
	auditTrailItem.TimeStamp = time.Now()
	data, err := json.Marshal(auditTrailItem)
	if err != nil {
		log.Fatalf("Error marshalling %s", auditTrailItem)
	}
	_, err = kp.kinesisClient.PutRecord(&kinesis.PutRecordInput{
		Data:         data,
		StreamName:   aws.String(streamName),
		PartitionKey: aws.String(partitionKey),
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (kp *kinesisPublisherImpl) GetRecords(streamName string) error {
	streamDescriptionOutput, err := kp.DescribeStream(streamName)
	if err != nil {
		return err
	}

	for _, shard := range streamDescriptionOutput.StreamDescription.Shards {
		iteratorOutput, err := kp.kinesisClient.GetShardIterator(&kinesis.GetShardIteratorInput{
			ShardId:           shard.ShardId,
			ShardIteratorType: aws.String("TRIM_HORIZON"),
			StreamName:        aws.String(streamName),
		})
		if err != nil {
			return err
		}
		shardIterator := iteratorOutput.ShardIterator

		for {
			resp, err := kp.kinesisClient.GetRecords(&kinesis.GetRecordsInput{
				ShardIterator: shardIterator,
			})
			if err != nil {
				return err
			}

			for _, record := range resp.Records {
				auditTrailItem := Entry{}
				err := json.Unmarshal([]byte(record.Data), &auditTrailItem)
				if err != nil {
					log.Printf("Error unmarshalling %s", record.Data)
					return err
				}
				AuditTrailChannel <- auditTrailItem
			}

			if resp.NextShardIterator == nil || shardIterator == resp.NextShardIterator {
				return nil
			}

			shardIterator = resp.NextShardIterator
		}
	}
	return nil
}

func (kp *kinesisPublisherImpl) Disconnect() error {
	return nil
}

func (kp *kinesisPublisherImpl) Connect() {
	creds := credentials.NewStaticCredentials(kp.config.AccessKey, kp.config.SecretKey, "")
	s := session.New(&aws.Config{Region: aws.String(kp.config.RegionName), Credentials: creds})
	kp.kinesisClient = kinesis.New(s)
}

func (kp *kinesisPublisherImpl) StreamExists(streamName string) (bool, error) {
	_, err := kp.DescribeStream(streamName)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "ResourceNotFoundException" {
				return false, nil
			}
		}
		return false, err
	}
	return true, nil
}

func (kp *kinesisPublisherImpl) DescribeStream(streamName string) (*kinesis.DescribeStreamOutput, error) {
	return kp.kinesisClient.DescribeStream(&kinesis.DescribeStreamInput{StreamName: aws.String(streamName)})
}

// SetupKinesisPubSub sets up connection to kinesis and creates the stream if does not exist
func SetupKinesisPubSub(config ConnectionConfig, streamName string) KinesisPublisher {
	kp, err := New(config)
	if err != nil {
		panic(err)
	}

	kp.Connect()

	exists, err := kp.StreamExists(streamName)
	if err != nil {
		panic(err)
	}

	if !exists {
		err = kp.CreateStream(streamName)
		if err != nil {
			panic(err)
		}
	}
	KPInstance = kp
	return kp
}

// PublishAuditTrailItem creates Audit Trail entry and publishes it
func PublishAuditTrailItem(streamName string, orderID string, status string, ruleName string, description string) {
	auditTrailItem := Entry{
		OrderID:     orderID,
		Status:      status,
		RuleName:    ruleName,
		Description: description,
		TimeStamp:   time.Now(),
	}
	KPInstance.Publish(streamName, auditTrailItem)
}
