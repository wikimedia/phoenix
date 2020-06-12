
GOOS        := linux
GOARCH      := amd64
CGO_ENABLED := 0
BINARY      := main
SOURCES     := main.go mw.go schema_org.go

build: clean
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build -o $(BINARY) $(SOURCES)
	zip function.zip main

deploy: build
	aws lambda update-function-code --function-name scpoc-lambda-fetch-schemaorg --zip-file fileb://function.zip

clean:
	rm -f function.zip $(BINARY)

.PHONY: build clean deploy
