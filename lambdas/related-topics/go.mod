module github.com/wikimedia/phoenix/lambdas/related-topics

go 1.15

replace github.com/wikimedia/phoenix/storage => /home/eevans/dev/src/git/phoenix/storage

replace github.com/wikimedia/phoenix/common => /home/eevans/dev/src/git/phoenix/common

require (
	github.com/aws/aws-lambda-go v1.20.0
	github.com/aws/aws-sdk-go v1.36.1
	github.com/wikimedia/phoenix/common v0.0.0-20201201202245-9b0069be3ccb
	github.com/wikimedia/phoenix/storage v0.0.0-20201201202245-9b0069be3ccb
)
