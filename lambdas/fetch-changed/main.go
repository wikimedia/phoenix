package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wikimedia/phoenix/common"
	"github.com/wikimedia/phoenix/env"
)

const folderName string = "incoming"

var debug bool = false

func keyf(msg *common.ChangeEvent) string {
	return fmt.Sprintf("%s/%s/%s-%d", folderName, msg.ServerName, msg.Title, msg.Revision)
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

func logDebug(format string, v ...interface{}) {
	if debug {
		fmt.Printf("[DEBUG] %s\n", fmt.Sprintf(strings.TrimSuffix(format, "\n"), v...))
	}
}

func handleRequest(ctx context.Context, event events.SNSEvent) {
	s3client := s3.New(session.New(&aws.Config{
		Region: aws.String(env.S3RawContentStorage().AWSConfig().Region()),
	}))

	snsPub := common.NewPublisher(env.SNSRawContentIncoming().ARN())

	for _, record := range event.Records {
		snsRecord := record.SNS

		msg := &common.ChangeEvent{}
		if err := json.Unmarshal([]byte(snsRecord.Message), msg); err != nil {
			fmt.Println("[Error] Unable to deserialize message payload:", err)
			continue
		}

		logDebug("%+v", msg)

		page, err := getPage(msg)
		if err != nil {
			fmt.Printf("[Error] Unable to retrieve %s (%s)\n", urlf(msg), err)
			continue
		}

		logDebug("Retrieved %s, %d bytes", urlf(msg), len(page))

		input := &s3.PutObjectInput{
			Body:   aws.ReadSeekCloser(bytes.NewReader(page)),
			Bucket: aws.String(env.S3RawContentStorage().Name()),
			Key:    aws.String(keyf(msg)),
			Metadata: map[string]*string{
				"title":       aws.String(msg.Title),
				"server_name": aws.String(msg.ServerName),
				"revision":    aws.String(fmt.Sprintf("%d", msg.Revision)),
			},
		}

		logDebug("%+v", input)

		result, err := s3client.PutObject(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Printf("[Error] %s: %s (%+v)\n", aerr.Code(), aerr.Message(), aerr.OrigErr())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			continue
		}

		logDebug("%+v", result)

		output, err := snsPub.SendChangeEvent(msg)
		if err != nil {
			fmt.Printf("[Error] Unable to send SNS change event: %s", err)
			continue
		}

		logDebug("%+v", output)
	}
}

func init() {
	if val, ok := os.LookupEnv("WMDEBUG"); ok {
		if val != "0" || strings.ToLower(val) != "false" {
			debug = true
			fmt.Println("DEBUG LOGGING ENABLED (unset WMDEBUG env var to disable)")
		}
	}
}

func main() {
	lambda.Start(handleRequest)
}
