[![this is fine](https://img.shields.io/badge/Dev%20status-Works%20For%20Me-red.svg)](../../docs/Status.md#works-for-me)

allpages
========

Iterate MediaWiki titles, print them to standard out.

```console
$ ./allpages -h
Usage of ./allpages:
  -from string
        Start iterating from closest matching title
  -server-name string
        Wiki server name (default "simple.wikipedia.org")
$ # Output format is: <revision> <title>
$ ./allpages -from E
6173530 E. E. Cummings
6624167 E. L. Doctorow
5665538 E. R. Braithwaite
6786663 E. Ruth Anderson
6942382 E
6622285 E-book reader
6871148 E-mart
...
$
```


Known issues
------------

- [ ] API access is done in a very bespoke manner; At some point, it would make sense to use a
      generic Action API abstraction (or build one, if necessary)
