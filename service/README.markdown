[![this is fine](https://img.shields.io/badge/Dev%20status-Works%20For%20Me-red.svg)](../../docs/Status.md#works-for-me)

# service

Running:

```sh-session
$ go run service.go
```

```sh-session
$ # Meanwhile, in an adjacent terminal...
$ # Query by page ID
$ curl -XPOST -d @query_by_id.json localhost:8080/query | json_pp
{
  "data" : {
    "page" : {
       "about" : [
          {
             "key" : "//schema.org",
             "val" : "/data/b6f7c05a-d367-11ea-af5c-2b020c033632"
          },
          {
             "key" : "//purl.org/dc/terms",
             "val" : "/data/b6f7c05a-d367-11ea-af5c-2b020c033632"
          }
       ],
       "dateModified" : "2020-07-31T16:07:04-05:00",
       "hasPart" : [
          "/page/385d6436a06b99d",
          "/page/644ed20cc75621c",
          "/page/42945840d44937c"
       ],
       "id" : "abcdefghijklmn",
       "name" : "Foobar",
       "url" : "//en.wikipedia.org/wiki/Foobar"
    }
  }
}
$ # Query by page name (Foobar)
$ curl -XPOST -d @query_by_name.json localhost:8080/query | json_pp
{
  "data" : {
    "page" : {
       "about" : [
          {
             "key" : "//schema.org",
             "val" : "/data/b6f7c05a-d367-11ea-af5c-2b020c033632"
          },
          {
             "key" : "//purl.org/dc/terms",
             "val" : "/data/b6f7c05a-d367-11ea-af5c-2b020c033632"
          }
       ],
       "dateModified" : "2020-07-31T16:07:04-05:00",
       "hasPart" : [
          "/page/385d6436a06b99d",
          "/page/644ed20cc75621c",
          "/page/42945840d44937c"
       ],
       "id" : "abcdefghijklmn",
       "name" : "Foobar",
       "url" : "//en.wikipedia.org/wiki/Foobar"
    }
  }
}
$
```
