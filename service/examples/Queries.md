# Sample Queries

## Pages

Querying for a page by its name, returning the first 3 child nodes:

    {
      page(name: { authority: "simple.wikipedia.org", name: "Banana" }) {
        name
        dateModified
        hasPart(limit: 3, offset: 0) {
          id
          name
          dateModified
        }
        about {
          key
          val
        }
      }
    }

Querying for a page by its ID, returning all but the first 3 child nodes:

    {
      page(name: { id: "/page/c52e8f8b0808caa" }) {
        name
        dateModified
        hasPart(offset: 3) {
          id
          name
          dateModified
        }
        about {
          key
          val
        }
      }
    }

## Nodes

Querying a node by its name:

    {
      node(name: { authority: "simple.wikipedia.org", pageName: "Banana", name: "Fruit" } ) {
        dateModified
        name
        unsafe
      }
    }

Querying a node by its id:

    {
      node(name: { id: "/node/5507c30ba578cdbe" } ) {
        dateModified
        name
        unsafe
      }
    }

## Linked data

Query a page with a specific `about` (by its key):

    {
      page(name: { authority: "simple.wikipedia.org", name: "Banana" }) {
        name
        dateModified
        hasPart(limit: 3, offset: 0) {
          id
          name
          dateModified
        }
        about(key: \"//schema.org\") {
          val
        }
      }
    }

