package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wikimedia/phoenix/common"
	"github.com/wikimedia/phoenix/storage"
)

var (
	content storage.Repository
	log     *common.Logger

	// These values are passed in at build-time w/ -ldflags (see: Makefile)
	awsAccount                string
	awsRegion                 string
	s3StructuredContentBucket string
)

// ~~~~~~~~~~
// TODO: Only a stub; Process `Node.Unsafe` to create actual related topics structure
// ~~~~~~~~~~
func getTopics(node *common.Node) ([]common.RelatedTopic, error) {
	// Faked/hard-coded topic(s)
	var topics = []common.RelatedTopic{
		{
			ID:       "Q1",
			Label:    "totality consisting of space, time, matter and energy",
			Salience: .99,
		},
		{
			ID:       "Q2",
			Label:    "third planet from the Sun in the Solar System",
			Salience: .99,
		},
		{
			ID:       "Q3",
			Label:    "matter capable of extracting energy from the environment for replication",
			Salience: .99,
		},
	}

	return topics, nil
}

func handleRequest(ctx context.Context, event events.SNSEvent) {
	for _, record := range event.Records {
		var err error
		var msg = &common.NodeStoredEvent{}
		var node *common.Node
		var topics []common.RelatedTopic

		// Deserialize message
		if err = json.Unmarshal([]byte(record.SNS.Message), msg); err != nil {
			log.Error("Unable to deserialize message payload: %s", err)
			continue
		}

		log.Debug("Processing Node published event: %+v", msg)

		// Retrieve Node object from storage
		if node, err = content.GetNode(msg.ID); err != nil {
			log.Error("Failed to retreive S3 object for node (ID=%s): %s", msg.ID, err)
			continue
		}

		log.Debug("Processing Node.Unsafe='%.24s...'", node.Unsafe)

		// Fetch related-topics
		if topics, err = getTopics(node); err != nil {
			log.Error("Unable to retrieve related topics for %s: %s", msg.ID, err)
			continue
		}

		// Store related topics
		if err = content.PutTopics(node, topics); err != nil {
			log.Error("Failed to store related-topics: %s", err)
		} else {
			// Index...
		}
	}
}

func init() {
	// AWS S3 client obj
	content = storage.Repository{
		Store:  s3.New(session.New(&aws.Config{Region: aws.String(awsRegion)})),
		Bucket: s3StructuredContentBucket,
	}

	// Determine logging level
	var level string = "ERROR"
	if v, ok := os.LookupEnv("LOG_LEVEL"); ok {
		level = v
	}

	// Initialize the logger
	log = common.NewLogger(level)
	log.Warn("%s LOGGING ENABLED (use LOG_LEVEL env var to configure)", common.LevelString(log.Level))

	log.Debug("AWS account ......................: %s", awsAccount)
	log.Debug("AWS region .......................: %s", awsRegion)
	log.Debug("S3 structured content bucket .....: %s", s3StructuredContentBucket)
}

func main() {
	lambda.Start(handleRequest)
}
