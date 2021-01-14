
# Repository configuration
include ../env/config.mk
# User/dev overrides
include ../.config.mk

GOOS        ?= linux
GOARCH      ?= amd64
CGO_ENABLED := 0
BINARY      := service
SOURCES     := service.go

# Configuration
LDFLAGS  = -X main.awsRegion=$(PHX_DEFAULT_REGION)
LDFLAGS += -X main.dynamoDBPageTitles=$(PHX_DYNAMODB_PAGE_TITLES)
LDFLAGS += -X main.dynamoDBNodeNames=$(PHX_DYNAMODB_NODE_NAMES)
LDFLAGS += -X main.s3Bucket=$(PHX_S3_STRUCTURED_CONTENT_BUCKET)


build: clean
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build -ldflags '$(LDFLAGS)' -o $(BINARY) $(SOURCES)

clean:
	rm -f $(BINARY)

deploy: build
	@echo "*****************************"
	@echo "*** No deploy configured! ***"
	@echo "*****************************"

test:
	go test

.PHONY: build clean deploy test
