
include ../../env/config.mk

GOOS    := linux
BINARY  := main
SOURCES := main.go nodeParser.go pageParser.go parser.go

# Configuration
LDFLAGS  = -X main.awsAccount=$(PHX_ACCOUNT_ID)
LDFLAGS += -X main.awsRegion=$(PHX_DEFAULT_REGION)

LDFLAGS += -X main.s3RawBucket=$(PHX_S3_RAW_CONTENT_BUCKET)
LDFLAGS += -X main.s3RawFolder=$(PHX_S3_RAW_CONTENT_INCOMING)

LDFLAGS += -X main.s3StructuredContentBucket=$(PHX_S3_CANONICAL_CONTENT_BUCKET)



build: clean
	GOOS=$(GOOS) go build -ldflags '$(LDFLAGS)' -o $(BINARY) $(SOURCES)
	zip function.zip main

deploy: build
	aws lambda update-function-code --function-name $(PHX_LAMBDA_TRANSFORM_PARSOID) --zip-file fileb://function.zip

clean:
	rm -f function.zip $(BINARY)

.PHONY: build clean deploy
