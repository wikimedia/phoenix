package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wikimedia/phoenix/common"
	"github.com/wikimedia/phoenix/storage"
)

var (
	// These values are passed in at build-time w/ -ldflags (see: Makefile)
	awsRegion                 string
	awsAccount                string
	s3StructuredContentBucket string
	s3RawBucket               string
	s3RawFolder               string

	debug bool = false
	log   *common.Logger
)

func keyf(msg *common.ChangeEvent) string {
	return fmt.Sprintf("%s/%s/%s-%d", s3RawFolder, msg.ServerName, msg.Title, msg.Revision)
}

func handleRequest(ctx context.Context, event events.SNSEvent) {
	s3client := s3.New(session.New(&aws.Config{
		Region: aws.String(awsRegion),
	}))

	repo := storage.Repository{
		Store:  s3client,
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
			log.Error("Unable to retrieve HTML document from S3: %s", err)
			log.Error("Bucket %s, Key: %s", s3RawBucket, keyf(msg))
			continue
		}

		log.Debug("Create html doc")
		defer data.Body.Close()
		document, err := goquery.NewDocumentFromReader(data.Body)
		if err != nil {
			log.Error("Unable to create html document with error: %s", err)
			return
		}

		log.Debug("Parse html parsoid document")
		page, nodes, err := parseParsoidDocument(document)

		if err != nil {
			log.Error("Unable to parse parsoid documet with error: %s", err)
			return
		}

		log.Debug("Save canonical data")
		saveError := repo.Apply(&storage.Update{
			Page:   *page,
			Nodes:  nodes,
			Abouts: map[string]common.Thing{},
		})

		if saveError != nil {
			log.Error("Unable to save to storage: %s", saveError)
			return
		}

		log.Debug("Save page successfully")
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

	log.Debug("AWS account ..........: %s", awsAccount)
	log.Debug("AWS region ...........: %s", awsRegion)
	log.Debug("S3 structured content bucket ............: %s", s3StructuredContentBucket)
	log.Debug("S3 raw bucket ............: %s", s3RawBucket)
	log.Debug("S3 raw income folder ............: %s", s3RawFolder)
}

func main() {
	lambda.Start(handleRequest)
}
