event-bridge
============

Send [Wikimedia `recentchange` events][1] to an SNS topic.


Known issues
------------

- [ ] Hardcodes runtime values (S3 bucket & region, target wiki & namespace, etc)
- [ ] Should block indefinitely; Exits after a short time (exit status 0)


[1]: https://wikitech.wikimedia.org/wiki/Event_Platform/EventStreams
