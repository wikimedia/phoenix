fetch-changed
=============

An AWS Lambda subscribed to an SNS topic that receives page edit events.  The function downloads the
corresponding Parsoid HTML to an S3 bucket.

Deployment (requires `aws`):

```
$ make deploy
```


Troubleshooting
---------------

Set the `WMDEBUG` environment variable (to something other than `0` or `false`) to enable
verbose debug logging.


Issues
------

- [ ] Hardcodes runtime values (S3 bucket, region, etc)
