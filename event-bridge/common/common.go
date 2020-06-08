package common

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/wikimedia/phoenix/env"
)

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
	return p.SendMessage(&Message{Title: title, ServerName: serverName, Revision: revision})
}

// SendMessage publishes an SNS message
func (p *Publisher) SendMessage(msg *Message) (*sns.PublishOutput, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling SNS event: %s", err)
	}

	input := &sns.PublishInput{Message: aws.String(string(b)), TopicArn: aws.String(env.SNSEventStreamsBridge().ARN())}

	result, err := p.client.Publish(input)
	if err != nil {
		return nil, fmt.Errorf("Error publishing to SNS: %s", err)
	}
	return result, nil
}

// NewPublisher creates a Publisher
func NewPublisher() (*Publisher, error) {
	config := &aws.Config{Region: aws.String(env.SNSEventStreamsBridge().AWSConfig().Region())}
	sess, err := session.NewSession(config)
	if err != nil {
		return nil, fmt.Errorf("Error creating AWS session: %s", err)
	}

	return &Publisher{client: sns.New(sess)}, nil
}
