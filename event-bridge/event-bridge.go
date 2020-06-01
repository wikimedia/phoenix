package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/eevans/wikimedia/streams"
)

// AWS connection info
const arn string = "arn:aws:sns:us-east-2:113698225543:scpoc-event-streams-bridge"
const region string = "us-east-2"

// Event streams attributes
const namespace int = 0
const wiki string = "simplewiki"

var since = flag.String("since", "", "Offset to begin stream from (ISO8601 timestamp or milliseconds since epoch)")

type message struct {
	Title      string `json:"title"`
	ServerName string `json:"server_name"`
	Revision   int    `json:"revision"`
}

func main() {
	flag.Parse()

	events := streams.NewClient().Match("namespace", namespace).Match("wiki", wiki).Match("type", "edit")

	if *since != "" {
		// Validate that *since parses as RFC3339 (ISO8601)
		if _, err := time.Parse(time.RFC3339, *since); err != nil {
			// If not RFC3339/ISO8601, is it at least a number?
			if _, err = strconv.Atoi(*since); err != nil {
				fmt.Fprint(os.Stderr, "Invalid timestamp argument for `-since`!")
				os.Exit(1)
			}
		}

		events.Since = *since
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating new AWS session: %s", err)
		os.Exit(1)
	}

	client := sns.New(sess)

	err = events.RecentChanges(func(event streams.RecentChangeEvent) {
		fmt.Printf("Change event captured!\n")
		fmt.Printf("  Title ............: %s\n", event.Title)
		fmt.Printf("  Server name ......: %s\n", event.ServerName)
		fmt.Printf("  Wiki .............: %s\n", event.Wiki)
		fmt.Printf("  Namespace ........: %d\n", event.Namespace)
		fmt.Printf("  Type .............: %s\n", event.Type)
		fmt.Printf("  Revision .........: %d\n", event.Revision.New)
		fmt.Printf("  Timestamp ........: %s\n", event.Meta.Dt)

		msg := message{event.Title, event.ServerName, event.Revision.New}

		b, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("Error marshalling SNS event:", err)
			return
		}

		input := &sns.PublishInput{Message: aws.String(string(b)), TopicArn: aws.String(arn)}

		result, err := client.Publish(input)
		if err != nil {
			fmt.Println("Error publishing to SNS:", err)
			return
		}

		fmt.Printf("Queued as %s\n", *result.MessageId)
	})

	fmt.Println()

	fmt.Printf("RecentChanges: %s\n", err)
	fmt.Printf("Last timestamp: %s\n", events.LastTimestamp())
}
