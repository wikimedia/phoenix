package common

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

const topic string = "arn:aws:sns:us-east-2:113698225543:scpoc-event-streams-bridge"
const region string = "us-east-2"

// Message represents the JSON object published to SNS.
type Message struct {
	Title      string `json:"title"`
	ServerName string `json:"server_name"`
	Revision   int    `json:"revision"`
}

// Publisher wraps the AWS SDK for publishing Message structs
type Publisher struct {
	client *sns.SNS
}

// Send publishes an SNS message
func (p *Publisher) Send(serverName string, title string, revision int) (*sns.PublishOutput, error) {
	msg := Message{Title: title, ServerName: serverName, Revision: revision}

	b, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling SNS event: %s", err)
	}

	input := &sns.PublishInput{Message: aws.String(string(b)), TopicArn: aws.String(topic)}

	result, err := p.client.Publish(input)
	if err != nil {
		return nil, fmt.Errorf("Error publishing to SNS: %s", err)
	}
	return result, nil
}

// NewPublisher creates a Publisher
func NewPublisher() (*Publisher, error) {
	config := &aws.Config{Region: aws.String(region)}
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, fmt.Errorf("Error creating AWS session: %s", err)
	}

	return &Publisher{client: sns.New(sess)}, nil
}
