module github.com/wikimedia/phoenix/service

go 1.14

replace github.com/wikimedia/phoenix/storage => ../storage

require (
	github.com/aws/aws-sdk-go v1.36.27
	github.com/elastic/go-elasticsearch/v7 v7.10.0
	github.com/google/uuid v1.1.4 // indirect
	github.com/gorilla/handlers v1.5.1
	github.com/graph-gophers/graphql-go v0.0.0-20201113091052-beb923fada29
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/rs/cors v1.7.0
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/wikimedia/phoenix/common v0.0.0-20210113223703-1e9b4f02ef22
	github.com/wikimedia/phoenix/storage v0.0.0-20210113223703-1e9b4f02ef22
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b // indirect
	golang.org/x/text v0.3.5 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)
