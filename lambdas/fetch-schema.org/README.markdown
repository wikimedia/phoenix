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

Set the `WMDEBUG` environment variable (to something other than `0` or `false`) to enable
verbose debug logging.


[1]: https://aws.amazon.com/cli/