
# Repository configuration
include ../env/config.mk
# User/dev overrides
include ../.config.mk

GOOS    := linux
BINARY  := rosette
SOURCES := rosette.go

# Configuration
LDFLAGS += -X main.awsRegion=$(PHX_DEFAULT_REGION)
LDFLAGS += -X main.s3StructuredContentBucket=$(PHX_S3_STRUCTURED_CONTENT_BUCKET)
LDFLAGS += -X main.esEndpoint=$(PHX_SEARCH_ENDPOINT)
LDFLAGS += -X main.esIndex=$(PHX_SEARCH_IDX_TOPICS)
LDFLAGS += -X main.esUsername=$(PHX_SEARCH_USERNAME)
LDFLAGS += -X main.esPassword=$(PHX_SEARCH_PASSWORD)
LDFLAGS += -X main.rosetteAPIKey=$(PHX_ROSETTE_API_KEY)


build: clean
	GOOS=$(GOOS) go build -ldflags '$(LDFLAGS)' -o $(BINARY) $(SOURCES)

clean:
	rm -f $(BINARY)

test:
	go test

.PHONY: build clean deploy
