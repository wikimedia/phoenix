module transform-parsoid

go 1.14

replace github.com/wikimedia/phoenix/storage => /home/eevans/dev/src/git/phoenix/storage

replace github.com/wikimedia/phoenix/common => /home/eevans/dev/src/git/phoenix/common

require (
	github.com/PuerkitoBio/goquery v1.6.0
	github.com/andybalholm/cascadia v1.2.0 // indirect
	github.com/aws/aws-lambda-go v1.20.0
	github.com/aws/aws-sdk-go v1.35.34
	github.com/elastic/go-elasticsearch/v7 v7.10.0 // indirect
	github.com/wikimedia/phoenix/common v0.0.0-20201109145749-c23218c68a2d
	github.com/wikimedia/phoenix/storage v0.0.0-20201109145749-c23218c68a2d
)
