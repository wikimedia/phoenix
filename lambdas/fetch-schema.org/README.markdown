[![this is fine](https://img.shields.io/badge/Dev%20status-Works%20For%20Me-red.svg)](../../docs/Status.md#works-for-me)

fetch-schema.org
================

An AWS Lambda triggered when new Parsoid HTML content is uploaded to S3.  Retrieves properties from
Wikidata, creates schema.org structured data (in JSON-LD format), and uploads it to S3.

Deployment (requires [`aws`][1]):

```
$ make deploy
```


Troubleshooting
---------------

Set the `LOG_LEVEL` environment variable (to one of `ERROR`, `WARN`, `INFO`, or `DEBUG`).


[1]: https://aws.amazon.com/cli/
