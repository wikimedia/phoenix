package common

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/wikimedia/phoenix/env"
)

// Publisher wraps the AWS SDK for publishing Message structs
type Publisher struct {
	client   *sns.SNS
	topicARN string
}

// Send publishes an SNS message
func (p *Publisher) Send(msg *ChangeEvent) (*sns.PublishOutput, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling SNS event: %s", err)
	}

	input := &sns.PublishInput{Message: aws.String(string(b)), TopicArn: aws.String(p.topicARN)}

	result, err := p.client.Publish(input)
	if err != nil {
		return nil, fmt.Errorf("Error publishing to SNS: %w", err)
	}
	return result, nil
}

// NewPublisher creates a Publisher
func NewPublisher(topicARN string) *Publisher {
	config := &aws.Config{Region: aws.String(env.SNSEventStreamsBridge().AWSConfig().Region())}
	sess, err := session.NewSession(config)
	if err != nil {
		panic(err)
	}

	return &Publisher{client: sns.New(sess), topicARN: topicARN}
}
