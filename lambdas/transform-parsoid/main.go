package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/wikimedia/phoenix/common"
	"github.com/wikimedia/phoenix/storage"
)

// Name to assign (unnamed) lead/intro sections.
const leadSectionName = "__intro__"

var (
	// These values are passed in at build-time w/ -ldflags (see: Makefile)
	awsRegion                 string
	awsAccount                string
	dynamoDBPageTitles        string
	dynamoDBNodeNames         string
	s3StructuredContentBucket string
	s3RawBucket               string
	s3RawIncomeFolder         string
	s3RawLinkedFolder         string
	snsNodePublished          string

	debug bool = false
	log   *common.Logger
)

func keyf(msg *common.ChangeEvent) string {
	return fmt.Sprintf("%s/%s/%s-%d", s3RawIncomeFolder, msg.ServerName, msg.Title, msg.Revision)
}

func readLinkedData(s3client *s3.S3, msg *common.ChangeEvent) (*common.Thing, error) {
	meta, err := s3client.GetObject(
		&s3.GetObjectInput{
			Bucket: aws.String(s3RawBucket),
			Key:    aws.String(fmt.Sprintf("%s/%s/%s-%d.json", s3RawLinkedFolder, msg.ServerName, msg.Title, msg.Revision)),
		})
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve linked-data (Wikidata) from S3: %s", err)
	}
	defer meta.Body.Close()

	body, err := ioutil.ReadAll(meta.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %w", err)
	}

	thing := &common.Thing{}
	if err := json.Unmarshal(body, thing); err != nil {
		return nil, fmt.Errorf("Error unmarshalling JSON: %w", err)
	}
	return thing, nil
}

// A helper function that returns a PostPutNodeCallback function conditional on a env variable.
func postPutNodeCallback(snsClient *sns.SNS) func(node common.Node) error {
	var disabled = false

	if env, ok := os.LookupEnv("DISABLE_PUT_NODE_CALLBACK"); ok {
		if v, err := strconv.ParseBool(env); err == nil {
			disabled = v
		}
	}

	// If not disabled (read: if enabled), return a callback that will deliver a "node published" SNS message.
	if !disabled {
		return func(node common.Node) error {
			var b []byte
			var err error
			var input *sns.PublishInput
			var output *sns.PublishOutput

			// JSON-encoded SNS message
			if b, err = json.Marshal(common.NodeStoredEvent{ID: node.ID}); err != nil {
				log.Error("Unable to marhsal SNS event to JSON: %s", err)
				return fmt.Errorf("Unable to marshal SNS event to JSON: %w", err)
			}

			input = &sns.PublishInput{
				Message:  aws.String(string(b)),
				TopicArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:%s", awsRegion, awsAccount, snsNodePublished)),
			}

			// Publish to SNS
			if output, err = snsClient.Publish(input); err != nil {
				log.Error("Failed to publish SNS event: %s", err)
				return fmt.Errorf("Failed to publish SNS event: %w", err)
			}

			log.Debug("Successfully published SNS message: %s (%s stored)", *output.MessageId, node.ID)

			return nil
		}
	}

	// If disabled, return a no-op callback.
	return func(node common.Node) error {
		return nil
	}
}

func handleRequest(ctx context.Context, event events.SNSEvent) {
	awsSession := session.New(&aws.Config{Region: aws.String(awsRegion)})
	s3client := s3.New(awsSession)
	snsClient := sns.New(awsSession)

	repo := storage.Repository{
		Store:  s3client,
		Index:  &storage.DynamoDBIndex{Client: dynamodb.New(awsSession), TitlesTable: dynamoDBPageTitles, NamesTable: dynamoDBNodeNames},
		Bucket: s3StructuredContentBucket,
	}

	for _, record := range event.Records {
		msg := &common.ChangeEvent{}
		if err := json.Unmarshal([]byte(record.SNS.Message), msg); err != nil {
			log.Error("Unable to deserialize message payload:", err)
			continue
		}

		log.Debug("Processing change event: %+v", msg)

		data, err := s3client.GetObject(
			&s3.GetObjectInput{
				Bucket: aws.String(s3RawBucket),
				Key:    aws.String(keyf(msg)),
			})

		if err != nil {
			log.Error("Unable to retrieve %s/%s from S3: %s", s3RawBucket, keyf(msg), err)
			continue
		}

		defer data.Body.Close()

		log.Debug("Creating html document...")

		document, err := goquery.NewDocumentFromReader(data.Body)
		if err != nil {
			log.Error("Unable to create html document with error: %s", err)
			continue
		}

		log.Debug("Parsing html parsoid document...")

		page, nodes, err := parseParsoidDocument(document)

		if err != nil {
			log.Error("Unable to parse parsoid document (%+v) with error: %s (+%v)", msg, err)
			continue
		}

		log.Debug("Loading JSON-LD output from S3...")

		thing, err := readLinkedData(s3client, msg)
		if err != nil {
			log.Error("Unable to load linked data with error: %s", err)
			continue
		}

		log.Debug("Saving document in canonical format...")

		saveError := repo.Apply(&storage.Update{
			Page:   *page,
			Nodes:  nodes,
			Abouts: map[string]common.Thing{"//schema.org": *thing},

			// Send events for each node published
			PostPutNodeCallback: postPutNodeCallback(snsClient),
		})

		if saveError != nil {
			log.Error("Unable to save to storage: %s", saveError)
			continue
		}

		log.Debug("Page saved successfully")
	}
}

func init() {
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
	log.Debug("DynamoDB page titles table .......: %s", dynamoDBPageTitles)
	log.Debug("DynamoDB node names table ........: %s", dynamoDBNodeNames)
	log.Debug("S3 structured content bucket .....: %s", s3StructuredContentBucket)
	log.Debug("S3 raw content bucket ............: %s", s3RawBucket)
	log.Debug("S3 raw content incoming folder ...: %s", s3RawIncomeFolder)
	log.Debug("S3 raw content linked folder .....: %s", s3RawLinkedFolder)
	log.Debug("SNS node published topic .........: %s", snsNodePublished)
}

func main() {
	lambda.Start(handleRequest)
}
