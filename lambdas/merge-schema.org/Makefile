
include ../../env/config.mk

GOOS    ?= linux
BINARY  ?= main
SOURCES := main.go

# Configuration
LDFLAGS  = -X main.awsRegion=$(PHX_DEFAULT_REGION)
LDFLAGS += -X main.rawContentBucket=$(PHX_S3_RAW_CONTENT_BUCKET)
LDFLAGS += -X main.incomingFolder=$(PHX_S3_RAW_CONTENT_INCOMING)
LDFLAGS += -X main.ldFolder=$(PHX_S3_RAW_CONTENT_WD_LINKED)
LDFLAGS += -X main.htmlFolder=$(PHX_S3_RAW_CONTENT_LINKED_HTML)


build: clean
	GOOS=$(GOOS) go build -ldflags '$(LDFLAGS)' -o $(BINARY) $(SOURCES)
	zip function.zip main

deploy: build
	aws lambda update-function-code --function-name $(PHX_LAMBDA_MERGE_SCHEMAORG) --zip-file fileb://function.zip

clean:
	rm -f function.zip $(BINARY)
	
test:
	go test

.PHONY: build clean deploy
