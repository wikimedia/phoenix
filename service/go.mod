module github.com/wikimedia/phoenix/service

go 1.14

replace github.com/wikimedia/phoenix/storage => /home/eevans/dev/src/git/phoenix/storage

replace github.com/wikimedia/phoenix/common => /home/eevans/dev/src/git/phoenix/common

require (
	github.com/aws/aws-sdk-go v1.36.1
	github.com/gorilla/handlers v1.5.1
	github.com/graph-gophers/graphql-go v0.0.0-20201113091052-beb923fada29
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/rs/cors v1.7.0
	github.com/wikimedia/phoenix/common v0.0.0-20201201202245-9b0069be3ccb
	github.com/wikimedia/phoenix/storage v0.0.0-20201201202245-9b0069be3ccb
)
