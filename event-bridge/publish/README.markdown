[![this is fine](https://img.shields.io/badge/Dev%20status-Works%20For%20Me-red.svg)](https://media.giphy.com/media/z9AUvhAEiXOqA/source.gif)

publish
=======

Publish ad-hoc change events to an SNS topic.

    $ ./publish
    Usage: ./publish <server> [<title> <revision>]

To publish a single change event, invoke `publish` with arguments for the server name,
the document title, and revision ID.

    $ ./publish simple.wikipedia.org Banana 6934315
    Queued "Banana" as 094fad70-df19-5637-9b00-fa20209833ed

To publish an arbitrary number of change events, invoke `publish` with a single argument
of the server name, and provide the revision ID and title as space-separate values on 
standard-in.

    $ cat events.txt
    6949779 San Antonio
    6948419 Detroit
    6948410 Seattle
    $ ./publish simple.wikipedia.org  < events.txt 
    Queued "San Antonio" as 1f17eae2-0469-514a-9fc1-bdb304db8173
    Queued "Detroit" as 4cabb130-b427-57aa-a107-e661b4959793
    Queued "Seattle" as cff171fb-37ff-5044-9342-26aa374d507d


Known issues
------------

- [ ] Hardcodes runtime values (S3 bucket & region, target wiki & namespace, etc)
- [ ] No tests

