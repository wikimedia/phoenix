
# Repository configuration
include ../../env/config.mk
# User/dev overrides
include ../../.config.mk

GOOS    := linux
BINARY  := main
SOURCES := main.go

# Configuration
LDFLAGS  = -X main.awsAccount=$(PHX_ACCOUNT_ID)
LDFLAGS += -X main.awsRegion=$(PHX_DEFAULT_REGION)
LDFLAGS += -X main.snsTopic=$(PHX_SNS_RAW_CONTENT_INCOMING)
LDFLAGS += -X main.s3Bucket=$(PHX_S3_RAW_CONTENT_BUCKET)
LDFLAGS += -X main.s3Folder=$(PHX_S3_RAW_CONTENT_INCOMING)


build: clean
	GOOS=$(GOOS) go build -ldflags '$(LDFLAGS)' -o $(BINARY) $(SOURCES)
	zip function.zip main

deploy: build
	aws lambda update-function-code --function-name $(PHX_LAMBDA_FETCH_CHANGED) --zip-file fileb://function.zip

clean:
	rm -f function.zip $(BINARY)

test:
	go test

.PHONY: build clean deploy
