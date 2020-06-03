event-bridge
============

Support for sending Wikimedia change events to the Project Phoenix.


|      | Description |
| ---- | ----------- |
| `stream`  | Subscribes to the [`recentchange` event stream][1], and publishes change events to an SNS topic |
| `publish` | Publishes ad hoc change events to an SNS topic |


[1]: https://wikitech.wikimedia.org/wiki/Event_Platform/EventStreams