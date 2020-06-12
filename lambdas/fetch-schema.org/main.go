package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/wikimedia/phoenix/common"
	"github.com/wikimedia/phoenix/env"
)

const folderName string = "schema.org"
const userAgent string = "Phoenix_lambda/0.0.0"

var s3client *s3.S3
var snsclient *common.Publisher
var debug bool = false

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

// Copy-pasta from fetch-changed
func logDebug(format string, v ...interface{}) {
	if debug {
		fmt.Printf("[DEBUG] %s\n", fmt.Sprintf(strings.TrimSuffix(format, "\n"), v...))
	}
}

// Handle upload to AWS S3
func putObject(msg *common.ChangeEvent, thing *Thing) (*s3.PutObjectOutput, error) {
	var b []byte
	var err error
	var input *s3.PutObjectInput
	var s3res *s3.PutObjectOutput

	// Serialize the Thing to JSON
	if b, err = json.MarshalIndent(thing, "", "  "); err != nil {
		return nil, fmt.Errorf("Unable to marshal JSON: %w", err)
	}

	// Upload to S3
	input = &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader(b)),
		Bucket: aws.String(env.S3RawContentStorage().Name()),
		Key:    aws.String(fmt.Sprintf("%s/%s/%s-%d.json", folderName, msg.ServerName, msg.Title, msg.Revision)),
		Metadata: map[string]*string{
			"title":       aws.String(msg.Title),
			"server_name": aws.String(msg.ServerName),
			"revision":    aws.String(fmt.Sprintf("%d", msg.Revision)),
		},
	}

	s3res, err = s3client.PutObject(input)
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
			fmt.Println("[Error] Unable to deserialize change message payload:", err)
			continue
		}

		logDebug("Change event received: %+v", msg)

		// Get wikibase_item
		if wdItem, err = wikibaseItemID(msg.ServerName, msg.Title); err != nil {
			fmt.Println("[Error] Unable to retrieve wikibase_item:", err)
			continue
		}

		logDebug("Found wikibase_item: %s", wdItem)

		// Query Wikidata & create JSON+LD output
		if thing, err = schemaOrgItem(wdItem); err != nil {
			fmt.Println("[Error] Unable to retrieve schema.org item:", err)
			continue
		}

		logDebug("Mapped %s to schema.org/Thing: %+v", wdItem, thing)

		if s3res, err = putObject(msg, thing); err != nil {
			fmt.Println("[Error] Unable to upload S3 object:", err)
			continue
		}

		logDebug("Uploaded JSON-LD: %+v", s3res)

		// Publish SNS event
		snsres, err = snsclient.SendChangeEvent(msg)
		if err != nil {
			fmt.Printf("[Error] Unable to send SNS change event: %s", err)
			continue
		}

		logDebug("Published SNS event: %+v", snsres)

	}
}

func init() {
	// AWS S3 client obj
	region := env.S3RawContentStorage().AWSConfig().Region()
	s3client = s3.New(session.New(&aws.Config{Region: aws.String(region)}))

	// AWS SNS client obj
	topic := env.SNSRawContentSchemaOrg().ARN()
	snsclient = common.NewPublisher(topic)

	// Enable debug output if env var set
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
