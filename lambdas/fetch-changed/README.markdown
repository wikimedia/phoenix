[![this is fine](https://img.shields.io/badge/Dev%20status-Works%20For%20Me-red.svg)](../../docs/Status.md#works-for-me)

fetch-changed
=============

An AWS Lambda subscribed to an SNS topic that receives page edit events.  The function downloads the
corresponding Parsoid HTML to an S3 bucket.

Deployment (requires [`aws`][1]):

```
$ make deploy
```


Troubleshooting
---------------

Set the `LOG_LEVEL` environment variable (to one of `ERROR`, `WARN`, `INFO`, or `DEBUG`).


[1]: https://aws.amazon.com/cli/
