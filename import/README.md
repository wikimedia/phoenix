# README

This utility is meant to perform a one-off import of [Rosette](http://www.rosette.com)-generated related topics
for Phoenix documents. It works by iterating Phoenix nodes (which at the time of this writing correspond to
top-level sections) from a JSON-formatted list, submits the node text to Rosette for analysis, stores the
corresponding related topics to the content store (S3), and finally, indexes them in Elasticsearch. It does
this sequentially (no concurrency), and with rate-limiting, to accommodate the constraints imposed by Rosette.

See also: [phoenix/issues/91](https://github.com/wikimedia/phoenix/issues/91)

## Usage

    Usage of ./rosette:
      -debug-log string
    	    enable debug logging to file (default "/dev/null")
      -limit int
    	    number of items to process (default -1)
      -resume string
    	    node ID to resume from

The utility expects a JSON-formatted, AWS CLI-generated, DynamoDB table scan of `scpoc-dynamodb-node-names` on
standard-in. For example:

    $ ./rosette <(zcat node-names-scan_2021-01-27T16:35:06-06:00.json.gz)

## Gotchas

- Likely requires that you have the AWS CLI installed (or the contents of `~/.aws` setup appropriately, at least).

## What's here

| File                                              | Description                                                                  |
| ------------------------------------------------- | ---------------------------------------------------------------------------- |
| node-names-scan_2021-01-27T16:35:06-06:00.json.gz | Snapshot of the `scpoc-dynamodb-node-names` DynamoDB table (484,971 entries) |
