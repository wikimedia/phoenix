package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/wikimedia/phoenix/common"
	"github.com/wikimedia/phoenix/storage"
)

var (
	// These values are passed in at build-time w/ -ldflags (see: Makefile)
	awsAccount string
	awsRegion  string
	snsTopic   string
	s3Bucket   string
	s3Folder   string

	debug bool = false
	log   *common.Logger
)

// Use the existing endpoint for now for testing
func urlf(domain string, title string) string {
	return fmt.Sprintf("https://%s/api/rest_v1/page/html/%s", domain, url.PathEscape(title))
}

func requestParsoid(domain string, title string) (body io.ReadCloser, err error) {
	res, err := http.Get(urlf(domain, title))
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func main() {

	bodyReader, err := requestParsoid("simple.wikipedia.org", "Mars")
	if err != nil {
		return
	}
	defer bodyReader.Close()

	document, err := goquery.NewDocumentFromReader(bodyReader)
	if err != nil {
		return
	}

	s3client := s3.New(session.New(&aws.Config{
		Region: aws.String(awsRegion),
	}))

	repo := storage.Repository{
		Store:  s3client,
		Bucket: s3Bucket,
	}

	page, nodes, err := parseParsoidDocument(document)
	if err != nil {
		return
	}

	saveError := repo.Apply(&storage.Update{
		Page:   *page,
		Nodes:  nodes,
		Abouts: map[string]common.Thing{},
	})
	if saveError != nil {
		println(saveError)
	}
}
