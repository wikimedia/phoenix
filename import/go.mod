module github.com/wikimedia/phoenix/import

go 1.15

replace github.com/wikimedia/phoenix/rosette => ../rosette

require (
	github.com/aws/aws-sdk-go v1.36.31
	github.com/elastic/go-elasticsearch/v7 v7.10.0
	github.com/wikimedia/phoenix/common v0.0.0-20210122212136-06a4785bb422
	github.com/wikimedia/phoenix/rosette v0.0.0-00010101000000-000000000000
	github.com/wikimedia/phoenix/storage v0.0.0-20210122212136-06a4785bb422
)
