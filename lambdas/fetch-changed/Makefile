
GOOS    := linux
BINARY  := main
SOURCES := main.go

build: clean
	GOOS=$(GOOS) go build -o $(BINARY) $(SOURCES)
	zip function.zip main

deploy: build
	aws lambda update-function-code --function-name scpoc-fetch-changed --zip-file fileb://function.zip

clean:
	rm -f function.zip $(BINARY)

.PHONY: build clean deploy
