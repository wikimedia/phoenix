package main

import (
	"encoding/json"
	"fmt"

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

type message struct {
	Title      string `json:"title"`
	ServerName string `json:"server_name"`
	Revision   int    `json:"revision"`
}

func main() {
	var count int64 = 0
	events := streams.NewClient().Match("namespace", namespace).Match("wiki", wiki).Match("type", "edit")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		fmt.Println("Error creating new AWS session:", err)
		return
	}

	client := sns.New(sess)

	events.RecentChanges(func(event streams.RecentChangeEvent) {
		fmt.Printf("Change event captured!\n")
		fmt.Printf("  Title ............: %s\n", event.Title)
		fmt.Printf("  Server name ......: %s\n", event.ServerName)
		fmt.Printf("  Wiki .............: %s\n", event.Wiki)
		fmt.Printf("  Namespace ........: %d\n", event.Namespace)
		fmt.Printf("  Type .............: %s\n", event.Type)
		fmt.Printf("  Revision .........: %d\n", event.Revision.New)

		count++

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

		fmt.Printf("%s queued!\n", *result.MessageId)
	})

	fmt.Printf("\nExiting (%d raw events processed)", count)
}
