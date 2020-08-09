package main

import (
	"context"
	"encoding/json"
	"fmt"

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
	awsRegion       string = "us-east-1"
	awsAccount      string
	s3StorageBucket string = "storage.bucket.peter-test"
	s3RawBucket     string
	s3RawForlder    string

	debug bool = false
	log   *common.Logger
)

func keyf(msg *common.ChangeEvent) string {
	return fmt.Sprintf("%s/%s/%s-%d", s3RawForlder, msg.ServerName, msg.Title, msg.Revision)
}

func handleRequest(ctx context.Context, event events.SNSEvent) {
	s3client := s3.New(session.New(&aws.Config{
		Region: aws.String(awsRegion),
	}))

	repo := storage.Repository{
		Store:  s3client,
		Bucket: s3StorageBucket,
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
			continue
		}

		log.Debug("Create html doc")
		document, err := goquery.NewDocument(data.String())
		if err != nil {
			log.Error("Unable to create html document with error: %s", err)
			return
		}

		page, nodes, err := parseParsoidDocument(document)

		if err != nil {
			log.Error("Unable to parse parsoid documet with error: %s", err)
			return
		}

		saveError := repo.Apply(&storage.Update{
			Page:   *page,
			Nodes:  nodes,
			Abouts: map[string]common.Thing{},
		})

		if saveError != nil {
			log.Error("Unable to save to strage: %s", saveError)
			return
		}

		log.Debug("Save page successfully")
	}
}

func main() {
	lambda.Start(handleRequest)
}
