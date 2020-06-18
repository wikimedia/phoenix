package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wikimedia/phoenix/common"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	// These values are passed in at build-time w/ -ldflags (see: Makefile)
	awsRegion        string
	rawContentBucket string
	incomingFolder   string
	ldFolder         string
	htmlFolder       string

	// See init()
	s3client *s3.S3
	log      *common.Logger
)

func handleRequest(ctx context.Context, event events.SNSEvent) {
	for _, record := range event.Records {
		msg := &common.ChangeEvent{}
		if err := json.Unmarshal([]byte(record.SNS.Message), msg); err != nil {
			log.Error("Unable to deserialize change event message: %s", err)
			continue
		}

		log.Debug("Processing change event: %+v", msg)

		// Download the schema.org JSON-LD
		meta, err := s3client.GetObject(
			&s3.GetObjectInput{
				Bucket: aws.String(rawContentBucket),
				Key:    aws.String(fmt.Sprintf("%s/%s/%s-%d.json", ldFolder, msg.ServerName, msg.Title, msg.Revision)),
			})
		if err != nil {
			log.Error("Unable to retrieve linked-data (Wikidata) from S3: %s", err)
			continue
		}

		log.Debug("Linked-data download complete: %+v", meta)

		// Download the Parsoid HTML
		data, err := s3client.GetObject(
			&s3.GetObjectInput{
				Bucket: aws.String(rawContentBucket),
				Key:    aws.String(fmt.Sprintf("%s/%s/%s-%d", incomingFolder, msg.ServerName, msg.Title, msg.Revision)),
			})
		if err != nil {
			log.Error("Unable to retrieve HTML document from S3: %s", err)
			continue
		}

		log.Debug("HTML download complete: %+v", data)

		// Smash the two together

		// First: Copy the linked-data JSON from the S3 response to string
		buf := new(strings.Builder)
		if _, err = io.Copy(buf, meta.Body); err != nil {
			log.Error("Unable to copy linked-data from S3 response: %s", err)
			continue
		}

		// Next create a script node (<script>...</script>), and append the linked-data as a text node
		script := html.Node{Type: html.ElementNode, DataAtom: atom.Script, Data: "script"}
		script.Attr = []html.Attribute{{Key: "type", Val: "application/ld+json"}}
		script.AppendChild(&html.Node{Type: html.TextNode, Data: buf.String()})

		// Parse the HTML document
		// Note: This should probably use the tokenizer instead of parser, and stream the data from/to S3
		// as the transform is performed.
		doc, err := html.Parse(data.Body)
		if err != nil {
			log.Error("Unable to parse HTML document: %s", err)
		}

		// Find the head (<head>...</head>), and append the script node created above.
		var f func(*html.Node) bool
		f = func(n *html.Node) bool {
			if n.Type == html.ElementNode && n.Data == "head" {
				n.AppendChild(&script)
				// Our work is done, so signal a bail out.
				return false
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if !f(c) {
					break
				}
			}
			return true
		}

		f(doc)

		// (Re)render the document, writing output to a buffer

		var rendered *bytes.Buffer = new(bytes.Buffer)
		html.Render(rendered, doc)

		// Upload transformed HTML
		s3res, err := s3client.PutObject(
			&s3.PutObjectInput{
				Body:        aws.ReadSeekCloser(bytes.NewReader(rendered.Bytes())),
				Bucket:      aws.String(rawContentBucket),
				Key:         aws.String(fmt.Sprintf("%s/%s/%s-%d.html", htmlFolder, msg.ServerName, msg.Title, msg.Revision)),
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

		log.Debug("HTML upload complete: %+v", s3res)
	}
}

func init() {
	// Initialize the AWS S3 client
	s3client = s3.New(session.New(&aws.Config{Region: aws.String(awsRegion)}))

	// Determine logging level
	var level string = "ERROR"
	if v, ok := os.LookupEnv("LOG_LEVEL"); ok {
		level = v
	}

	// Initialize the logger
	log = common.NewLogger(level)
	log.Warn("%s LOGGING ENABLED (use LOG_LEVEL env var to configure)", common.LevelString(log.Level))

	log.Debug("AWS region ...........: %s", awsRegion)
	log.Debug("S3 bucket ............: %s", rawContentBucket)
	log.Debug("Incoming folder ......: %s", incomingFolder)
	log.Debug("Linked data folder ...: %s", ldFolder)
	log.Debug("Upload folder ........: %s", htmlFolder)
}

func main() {
	lambda.Start(handleRequest)
}
