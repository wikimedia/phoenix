package env

import "fmt"

// AWSAccountID is the Amazon Web Services account identifier
const AWSAccountID string = "113698225543"

// AWSRegion is the Amazon Web Services region identifier
const AWSRegion string = "us-east-2"

// Config returns an AWSConfig containing high-level Amazon Web Services configuration information
func Config() *AWSConfig {
	return &AWSConfig{account: AWSAccountID, region: AWSRegion}
}

// SNSEventStreamsBridge returns metadata for the topic that receives change events from the
// Wikimedia Event Streams service
func SNSEventStreamsBridge() *AWSResource {
	return &AWSResource{
		name:      "scpoc-event-streams-bridge",
		awsConfig: Config(),
		kind:      "sns",
	}
}

// SNSRawContentIncoming returns metadata for the topic that receives events when new HTML content
// is added to the "raw content" bucket
func SNSRawContentIncoming() *AWSResource {
	return &AWSResource{
		name:      "scpoc-sns-raw-content-incoming",
		awsConfig: Config(),
		kind:      "sns",
	}
}

// S3RawContentStorage returns metadata for the "raw content" S3 bucket
func S3RawContentStorage() *AWSResource {
	return &AWSResource{
		name:      "scpoc-raw-content-store",
		awsConfig: Config(),
		kind:      "s3",
	}
}

// AWSConfig represents high-level Amazon Web Services configuration
type AWSConfig struct {
	region  string
	account string
}

// Region returns the AWS region identifer
func (conf *AWSConfig) Region() string {
	return conf.region
}

// Account returns the AWS account identifier
func (conf *AWSConfig) Account() string {
	return conf.account
}

// AWSResource represents generic AWS resources
type AWSResource struct {
	name      string
	awsConfig *AWSConfig
	kind      string
}

// Name returns the resource name
func (rsrc *AWSResource) Name() string {
	return rsrc.name
}

// ARN returns the ARN value for the resource
func (rsrc *AWSResource) ARN() string {
	return fmt.Sprintf("arn:aws:%s:%s:%s:%s", rsrc.kind, rsrc.awsConfig.Region(), rsrc.awsConfig.Account(), rsrc.Name())
}

// AWSConfig returns the corresponding AWSConfig for this resource
func (rsrc *AWSResource) AWSConfig() *AWSConfig {
	return rsrc.awsConfig
}
