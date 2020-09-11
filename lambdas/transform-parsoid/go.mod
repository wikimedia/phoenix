module transform-parsoid

go 1.14

replace github.com/wikimedia/phoenix/storage => ../../storage

require (
	github.com/PuerkitoBio/goquery v1.5.1
	github.com/andybalholm/cascadia v1.2.0 // indirect
	github.com/aws/aws-lambda-go v1.19.1
	github.com/aws/aws-sdk-go v1.34.21
	github.com/wikimedia/phoenix/common v0.0.0-20200910210446-f96abf625df6
	github.com/wikimedia/phoenix/storage v0.0.0-20200911175454-f97fac6083b7
	golang.org/x/net v0.0.0-20200904194848-62affa334b73 // indirect
)
