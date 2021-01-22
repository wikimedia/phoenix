[![this is fine](https://img.shields.io/badge/Dev%20status-Works%20For%20Me-red.svg)](../../docs/Status.md#works-for-me)

# transform-parsoid

An AWS Lambda that decomposes Parsoid HTML documents into graphs of JSON objects, and stores
them to S3.

## Disabling outgoing node storage events

By default, an SNS message is sent for each new `Node` object stored (at the time of this
writing, used exclusively for related-topics processing of section data). To disable
publishing of these events, set the `DISABLE_PUT_NODE_CALLBACK` environment var to `true`.
