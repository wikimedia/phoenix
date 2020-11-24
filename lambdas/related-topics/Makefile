
include ../../env/config.mk

GOOS    := linux
BINARY  := main
SOURCES := main.go

# Configuration
LDFLAGS  = -X main.awsAccount=$(PHX_ACCOUNT_ID)
LDFLAGS += -X main.awsRegion=$(PHX_DEFAULT_REGION)
LDFLAGS += -X main.s3StructuredContentBucket=$(PHX_S3_STRUCTURED_CONTENT_BUCKET)


build: clean
	GOOS=$(GOOS) go build -ldflags '$(LDFLAGS)' -o $(BINARY) $(SOURCES)
	zip function.zip main

deploy: build
	aws lambda update-function-code --function-name $(PHX_LAMBDA_RELATED_TOPICS) --zip-file fileb://function.zip

clean:
	rm -f function.zip $(BINARY)

test:
	go test

.PHONY: build clean deploy
