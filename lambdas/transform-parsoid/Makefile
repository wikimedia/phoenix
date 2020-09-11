
include ../../env/config.mk

GOOS    := linux
BINARY  := main
SOURCES := main.go nodeParser.go pageParser.go parser.go

# Configuration
LDFLAGS  = -X main.awsAccount=$(PHX_ACCOUNT_ID)
LDFLAGS += -X main.awsRegion=$(PHX_DEFAULT_REGION)
LDFLAGS += -X main.dynamoDBPageTitles=$(PHX_DYNAMODB_PAGE_TITLES)
LDFLAGS += -X main.dynamoDBNodeNames=$(PHX_DYNAMODB_NODE_NAMES)
LDFLAGS += -X main.s3RawBucket=$(PHX_S3_RAW_CONTENT_BUCKET)
LDFLAGS += -X main.s3RawIncomeFolder=$(PHX_S3_RAW_CONTENT_INCOMING)
LDFLAGS += -X main.s3RawLinkedFolder=$(PHX_S3_RAW_CONTENT_WD_LINKED)
LDFLAGS += -X main.s3StructuredContentBucket=$(PHX_S3_STRUCTURED_CONTENT_BUCKET)


build: clean
	GOOS=$(GOOS) go build -ldflags '$(LDFLAGS)' -o $(BINARY) $(SOURCES)
	zip function.zip main

deploy: build
	aws lambda update-function-code --function-name $(PHX_LAMBDA_TRANSFORM_PARSOID) --zip-file fileb://function.zip

clean:
	rm -f function.zip $(BINARY)

test:
	go test

.PHONY: build clean deploy
