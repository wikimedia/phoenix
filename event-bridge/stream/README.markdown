[![this is fine](https://img.shields.io/badge/Dev%20status-Works%20For%20Me-red.svg)](../../docs/Status.md#works-for-me)

stream
======

Stream [Wikimedia change events][1] to an SNS topic.


Known issues
------------

- [ ] Hardcodes runtime values (S3 bucket & region, target wiki & namespace, etc)
- [ ] Should block indefinitely; Exits 15 minutes (see [T242767][2])
- [ ] No tests


[1]: https://wikitech.wikimedia.org/wiki/Event_Platform/EventStreams
[2]: https://phabricator.wikimedia.org/T242767
