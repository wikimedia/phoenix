module github.com/wikimedia/phoenix/lambdas/related-topics

go 1.15

replace github.com/wikimedia/phoenix/storage => ../../storage

require (
	github.com/aws/aws-lambda-go v1.22.0
	github.com/aws/aws-sdk-go v1.36.25
	github.com/elastic/go-elasticsearch/v7 v7.10.0
	github.com/google/uuid v1.1.4 // indirect
	github.com/wikimedia/phoenix/common v0.0.0-20210106213327-5044c4eca381
	github.com/wikimedia/phoenix/storage v0.0.0-20210106213327-5044c4eca381
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b // indirect
	golang.org/x/text v0.3.5 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
