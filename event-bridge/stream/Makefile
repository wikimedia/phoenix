
include ../../env/config.mk

BINARY  := stream
SOURCES := stream.go

# Configuration
LDFLAGS  = -X main.awsAccount=$(PHX_ACCOUNT_ID)
LDFLAGS += -X main.awsRegion=$(PHX_DEFAULT_REGION)
LDFLAGS += -X main.snsTopic=$(PHX_SNS_EVENT_STREAMS_BRIDGE)


all: clean
	go build -ldflags '$(LDFLAGS)' -o $(BINARY) $(SOURCES)

clean:
	rm -f $(BINARY)

.PHONY: all clean
