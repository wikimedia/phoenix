package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanity(t *testing.T) {

	assert.Equal(t, AWSRegion, SNSEventStreamsBridge().AWSConfig().Region(), "Regions should match")
	assert.Equal(t, AWSRegion, SNSRawContentIncoming().AWSConfig().Region(), "Regions should match")
	assert.Equal(t, AWSRegion, S3RawContentStorage().AWSConfig().Region(), "Regions should match")

	assert.Equal(t, AWSAccountID, SNSEventStreamsBridge().AWSConfig().Account(), "Account IDs should match")
	assert.Equal(t, AWSAccountID, SNSRawContentIncoming().AWSConfig().Account(), "Account IDs should match")
	assert.Equal(t, AWSAccountID, S3RawContentStorage().AWSConfig().Account(), "Account IDs should match")

	assert.NotEmpty(t, SNSEventStreamsBridge().ARN())
	assert.NotEmpty(t, SNSEventStreamsBridge().Name())

	assert.NotEmpty(t, SNSRawContentIncoming().ARN())
	assert.NotEmpty(t, SNSRawContentIncoming().Name())

	assert.NotEmpty(t, S3RawContentStorage().ARN())
	assert.NotEmpty(t, S3RawContentStorage().Name())

}
