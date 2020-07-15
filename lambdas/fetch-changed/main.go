package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wikimedia/phoenix/common"
)

var (
	// These values are passed in at build-time w/ -ldflags (see: Makefile)
	awsAccount string
	awsRegion  string
	snsTopic   string
	s3Bucket   string
	s3Folder   string

	debug bool = false
	log   *common.Logger
)

func keyf(msg *common.ChangeEvent) string {
	return fmt.Sprintf("%s/%s/%s-%d", s3Folder, msg.ServerName, msg.Title, msg.Revision)
}

func urlf(msg *common.ChangeEvent) string {
	return fmt.Sprintf("https://%s/api/rest_v1/page/html/%s/%d", msg.ServerName, url.PathEscape(msg.Title), msg.Revision)
}

func getPage(msg *common.ChangeEvent) ([]byte, error) {
	res, err := http.Get(urlf(msg))
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func handleRequest(ctx context.Context, event events.SNSEvent) {
	s3client := s3.New(session.New(&aws.Config{
		Region: aws.String(awsRegion),
	}))

	snsChangePub := common.NewChangeEventPublisher(awsAccount, awsRegion, snsTopic)

	for _, record := range event.Records {
		msg := &common.ChangeEvent{}
		if err := json.Unmarshal([]byte(record.SNS.Message), msg); err != nil {
			log.Error("Unable to deserialize message payload:", err)
			continue
		}

		log.Debug("Processing change event: %+v", msg)

		page, err := getPage(msg)
		if err != nil {
			log.Error("Unable to retrieve %s (%s)\n", urlf(msg), err)
			continue
		}

		log.Debug("Retrieved %s, %d bytes", urlf(msg), len(page))

		result, err := s3client.PutObject(
			&s3.PutObjectInput{
				Body:        aws.ReadSeekCloser(bytes.NewReader(page)),
				Bucket:      aws.String(s3Bucket),
				Key:         aws.String(keyf(msg)),
				ContentType: aws.String("text/html"),
				Metadata: map[string]*string{
					"title":       aws.String(msg.Title),
					"server_name": aws.String(msg.ServerName),
					"revision":    aws.String(fmt.Sprintf("%d", msg.Revision)),
				},
			})
		if err != nil {
			log.Error("Unable to upload HTML document to S3: %s", err)
			continue
		}

		log.Debug("HTML upload complete: %+v", result)

		output, err := snsChangePub.Send(msg)
		if err != nil {
			log.Error("Unable to send SNS change event: %s", err)
			continue
		}

		log.Debug("SNS change event sent: %+v", output)
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
	log.Debug("SNS topic ............: %s", snsTopic)
	log.Debug("S3 bucket ............: %s", s3Bucket)
	log.Debug("S3 folder ............: %s", s3Folder)
}

func main() {
	lambda.Start(handleRequest)
}
