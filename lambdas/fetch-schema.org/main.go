package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/wikimedia/phoenix/common"
)

const userAgent string = "Phoenix_lambda/0.0.0"

var (
	s3client  *s3.S3
	snsclient *common.Publisher
	log       *common.Logger

	awsRegion string
	s3Bucket  string
	s3Folder  string
	snsTopic  string
)

// Convenience for performing HTTP GETs and returning the entire page as []byte
func request(page string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", page, nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating HTTP request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP GET error: %w", err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP status %d (expected %d)", res.StatusCode, http.StatusOK)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading HTTP response body: %w", err)
	}

	return body, nil
}

// Handle upload to AWS S3
func putObject(msg *common.ChangeEvent, thing *Thing) (*s3.PutObjectOutput, error) {
	var b []byte
	var err error
	var s3res *s3.PutObjectOutput

	// Serialize the Thing to JSON
	if b, err = json.Marshal(thing); err != nil {
		return nil, fmt.Errorf("Unable to marshal JSON: %w", err)
	}

	// Upload to S3
	s3res, err = s3client.PutObject(
		&s3.PutObjectInput{
			Body:        aws.ReadSeekCloser(bytes.NewReader(b)),
			Bucket:      aws.String(s3Bucket),
			Key:         aws.String(fmt.Sprintf("%s/%s/%s-%d.json", s3Folder, msg.ServerName, msg.Title, msg.Revision)),
			ContentType: aws.String("application/json"),
			Metadata: map[string]*string{
				"title":       aws.String(msg.Title),
				"server_name": aws.String(msg.ServerName),
				"revision":    aws.String(fmt.Sprintf("%d", msg.Revision)),
			},
		})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return nil, fmt.Errorf("%s: %s (%+v)", aerr.Code(), aerr.Message(), aerr.OrigErr())
		}

		return nil, err
	}

	return s3res, nil
}

func handleRequest(ctx context.Context, event events.SNSEvent) {
	for _, record := range event.Records {
		var msg *common.ChangeEvent
		var wdItem string
		var err error
		var thing *Thing
		var s3res *s3.PutObjectOutput
		var snsres *sns.PublishOutput

		msg = &common.ChangeEvent{}
		if err := json.Unmarshal([]byte(record.SNS.Message), msg); err != nil {
			log.Error("Unable to deserialize change message payload: %s", err)
			continue
		}

		log.Debug("Processing change event: %+v", msg)

		// Get wikibase_item
		if wdItem, err = wikibaseItemID(msg.ServerName, msg.Title); err != nil {
			log.Error("Unable to retrieve wikibase_item: %s", err)
			continue
		}

		log.Debug("Found wikibase_item: %s", wdItem)

		// Query Wikidata & create JSON+LD output
		if thing, err = schemaOrgItem(wdItem); err != nil {
			log.Error("Unable to query schema.org item attributes: %s", err)
			continue
		}

		log.Debug("Mapped %s to schema.org/Thing: %+v", wdItem, thing)

		if s3res, err = putObject(msg, thing); err != nil {
			log.Error("Unable to upload JSON object to S3: %s", err)
			continue
		}

		log.Debug("Uploaded JSON-LD to S3: %+v", s3res)

		// Publish SNS event
		snsres, err = snsclient.SendChangeEvent(msg)
		if err != nil {
			log.Error("Unable to send SNS change event: %s", err)
			continue
		}

		log.Debug("Published SNS event: %+v", snsres)

	}
}

func init() {
	// AWS S3 client obj
	s3client = s3.New(session.New(&aws.Config{Region: aws.String(awsRegion)}))

	// AWS SNS client obj
	snsclient = common.NewPublisher(snsTopic)

	// Determine logging level
	var level string = "ERROR"
	if v, ok := os.LookupEnv("LOG_LEVEL"); ok {
		level = v
	}

	// Initialize the logger
	log = common.NewLogger(level)
	log.Warn("%s LOGGING ENABLED (use LOG_LEVEL env var to configure)", common.LevelString(log.Level))
}

func main() {
	lambda.Start(handleRequest)
}
