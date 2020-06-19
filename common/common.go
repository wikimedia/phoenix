package common

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// ChangeEventPublisher wraps the AWS SDK for publishing Message structs
type ChangeEventPublisher struct {
	client   *sns.SNS
	account  string
	region   string
	topic    string
	topicARN string
}

// Send publishes an SNS message for a ChangeEvent
func (p *ChangeEventPublisher) Send(msg *ChangeEvent) (*sns.PublishOutput, error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling SNS event: %s", err)
	}

	input := &sns.PublishInput{Message: aws.String(string(b)), TopicArn: aws.String(p.arn())}

	result, err := p.client.Publish(input)
	if err != nil {
		return nil, fmt.Errorf("Error publishing to SNS: %w", err)
	}
	return result, nil
}

func (p *ChangeEventPublisher) arn() string {
	return fmt.Sprintf("arn:aws:sns:%s:%s:%s", p.region, p.account, p.topic)
}

// NewChangeEventPublisher returns an initialized ChangeEventPublisher
func NewChangeEventPublisher(account string, region string, topic string) *ChangeEventPublisher {
	config := &aws.Config{Region: aws.String(region)}
	sess, err := session.NewSession(config)
	if err != nil {
		panic(err)
	}

	return &ChangeEventPublisher{client: sns.New(sess), account: account, region: region, topic: topic}
}
