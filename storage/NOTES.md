# NOTES

Stored objects (articles, sections, and metadata) are supposed to have a Phoenix-specific surrogate
lookup key. Generating one that is stable though would require storing a mapping of ID to
MediaWiki page ID (MediaWiki is after all, canonical). Instead of generating one though, we can
_fake_ an ID by using a hash of the wiki name and page ID.

## Storage and indexing

The tentative/proposed approach is to store the objects that make up the document graph as
JSON-encoded files in S3, and use Elasticsearch to index them for queries. This is NOT an
endorsement for this approach if this ends up being moved from proof-of-concept to production; It
should be Good Enough for our purposes here, and optimizes for flexability as we iterate to figure
things out.

### Objects (JSON)

#### page

Keyed by: `{bucket}/{wiki}/content/{hash({wiki}+{pageid})}.json`

```json
{
  "identifier": "/content/ddc93d9352c6de4f",
  "_source": {
    "id": 31769,
    "revision": 6949779,
    "tid": "5d6d1580-bf8f-11ea-a5d3-5d162f2fa7a5",
    "authority": "simple.wikipedia.org"
  },
  "name": "San Antonio",
  "url": "//simple.wikipedia.org/wiki/San_Antonio",
  "dateModified": "2020-07-10T16:04:16-05:00",
  "hasPart": ["/content/ea6c21de-c2f0-11ea-84e5-54e1ad382fa6"],
  "about": {
    "schema.org": "/data/ea6c221b-c2f0-11ea-84e5-54e1ad382fa6"
  }
}
```

#### section

Keyed by: `{bucket}/{wiki}/content/{uuid}.json`

```json
{
  "id": "/content/d00d0dad-c2df-11ea-a292-54e1ad382fa6",
  "isPartOf": ["/content/ddc93d9352c6de4f"],
  "dateModified": "2020-07-10T14:01:51-05:00",
  "unsafe": "<p>San Antonio is a large city in southern Texas, USA...</p>"
}
```

#### metadata

Keyed by: `{bucket}/{wiki}/metadata/{uuid}.json`

```json
{
  "alternateName": "Alamo City",
  "description": "second-most populous city in Texas, United States of America",
  "name": "San Antonio",
  "sameAs": "https://www.wikidata.org/wiki/Q975",
  "@context": "https://schema.org",
  "@type": "Thing"
}
```

### Sequencing of operations

1. Retrieve current page object
1. Write metadata & section objects
1. Write page object
1. Delete previous section & metadata objects (as referenced in old page object)
1. Index new document(s)

### Caveats

The document structure we're utilizing implies that sections can exist independantly of the
document, and can be referenced by more than one document. So, while we're using a structure that
might support it, we'll never be in a position to implement these semantics so long as we are a
transform of content from MediaWiki. As a result, the storage semantics above assume that metadata
and section objects move in lock-step with pages and are (re)generated and overwritten on each
update. Likewise, sections reference parents (plural), even if it's only ever just the one.
