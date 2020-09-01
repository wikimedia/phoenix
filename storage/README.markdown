[![this is fine](https://img.shields.io/badge/Dev%20status-Works%20For%20Me-red.svg)](../../docs/Status.md#works-for-me)

# storage

## Running tests

### Mocked storage

To run the tests against mocked storage:

    $ go test

### S3

To run the tests against Amazon S3 storage:

    $ TESTS_USE_S3=1 go test
    $ TESTS_USE_S3=1 AWS_REGION=us-west-1 go test
    $ TESTS_USE_S3=1 AWS_BUCKET=my-bucket go test

When unset, `AWS_REGION` defaults to `us-east-2`
When unset, `AWS_BUCKET` defaults to `scpoc-structured-content-store`.

### DynamoDB

To run the tests against DynamoDB (indexing tests):

    $ TESTS_USE_DYNAMODB=1 AWS_DYNAMODB_PAGE_TITLES_TABLE=test-table go test

### Elasticsearch

To run the tests against Elasticsearch (indexing tests):

    $ curl -D - -XPUT -H 'Content-Type: application/json' -d @page_name_mapping.json http://localhost:9200/page_name; echo
    $ TESTS_USE_ELASTICSEARCH=1 go test
