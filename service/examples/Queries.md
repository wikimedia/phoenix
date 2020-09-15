Sample Queries
==============

Pages
-----

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


Nodes
-----
