
include ../../env/config.mk

GOOS        := linux
GOARCH      := amd64
CGO_ENABLED := 0
BINARY      := main
SOURCES     := main.go mw.go schema_org.go

# Configuration
LDFLAGS  = -X main.awsRegion=$(PHX_DEFAULT_REGION)
LDFLAGS += -X main.awsAccount=$(PHX_ACCOUNT_ID)
LDFLAGS += -X main.s3Bucket=$(PHX_S3_RAW_CONTENT_BUCKET)
LDFLAGS += -X main.s3Folder=$(PHX_S3_RAW_CONTENT_WD_LINKED)
LDFLAGS += -X main.snsTopic=$(PHX_SNS_RAW_CONTENT_WD_LINKED)


build: clean
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build -ldflags '$(LDFLAGS)' -o $(BINARY) $(SOURCES)
	zip function.zip main

deploy: build
	aws lambda update-function-code --function-name $(PHX_LAMBDA_FETCH_SCHEMAORG) --zip-file fileb://function.zip

clean:
	rm -f function.zip $(BINARY) -ldflags '$(LDFLAGS)'

test:
	go test

.PHONY: build clean deploy
