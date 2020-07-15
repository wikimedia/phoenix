[![this is fine](https://img.shields.io/badge/Dev%20status-Works%20For%20Me-red.svg)](../../docs/Status.md#works-for-me)

merge-schema.org
================

An AWS Lambda triggered when new linked data is uploaded to S3.  The corresponding HTML
document is updated to include the linked data (as JSON-LD), and uploaded to S3.

Deployment (requires [`aws`][1]):

```
$ make deploy
```


Troubleshooting
---------------

Set the `LOG_LEVEL` environment variable (to one of `ERROR`, `WARN`, `INFO`, or `DEBUG`).


[1]: https://aws.amazon.com/cli/