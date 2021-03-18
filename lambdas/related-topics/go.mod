module github.com/wikimedia/phoenix/lambdas/related-topics

go 1.15

replace github.com/wikimedia/phoenix/storage => ../../storage

require (
	github.com/andybalholm/cascadia v1.2.0 // indirect
	github.com/aws/aws-lambda-go v1.22.0
	github.com/aws/aws-sdk-go v1.36.31
	github.com/elastic/go-elasticsearch/v7 v7.10.0
	github.com/google/uuid v1.2.0 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/wikimedia/phoenix/common v0.0.0-20210122212136-06a4785bb422
	github.com/wikimedia/phoenix/rosette v0.0.0-20210223213805-8f53ab7329f0
	github.com/wikimedia/phoenix/storage v0.0.0-20210122183222-d75f3fd4ef67
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/text v0.3.5 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
