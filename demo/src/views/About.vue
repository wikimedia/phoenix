<template>
  <div class="about">
    <v-container>
      <h1>About Phoenix</h1>
      <br>
      <h2>What is Phoenix?</h2>
      <p>
        <a href="https://github.com/wikimedia/phoenix" target="_blank">Phoenix</a> is an experimental service
        demonstrating the value of a structured content store. In addition to its functional capabilities,
        Phoenix explores emerging architecture patterns that will direct modernization efforts at Wikimedia.
        To learn more about Phoenix and the architecture process, read the
        <a href="https://www.mediawiki.org/wiki/Architecture_Repository/Strategy/Goals_and_initiatives/Structured_content_proof_of_value" target="_blank">artifact</a>.
      </p>
      <h2>How does it work?</h2>
      <p>
        Phoenix consumes a limited set of articles from
        <a href="https://simple.wikipedia.org" target="_blank">Simple English Wikipedia</a>
        and structures the content into sections. The content store is updated in response to changes to the
        original articles.
      </p>
      <p>
        Linking between sections and keywords from <a href="https://www.wikidata.org" target="_blank">Wikidata</a>
        is provided by <a href="https://www.rosette.com/" target="_blank">Rosette</a>. Integration with Rosette
        is intended to be a short-term feature for experimental use only. Phoenix is capable of updating keywords
        in response to changes to the original content, but this feature is currently enabled for only a
        <a href="https://github.com/wikimedia/phoenix/blob/master/event-bridge/stream/allowed.yaml" target="_blank">limited set of content</a>.
      </p>
      <p>
        This site connects to Phoenix using a GraphQL API. You can try out the API using the
        <a href="/sandbox">API sandbox</a>. To explore the schema, select <i>Docs</i> on the
        right side of the sandbox.
      </p>
      <h3>Example queries</h3>
      <pre
        class="body-2"
      >
  # List sections in an article
  {
    page(name: { authority: "simple.wikipedia.org", name: "Banana"} ) {
      name
      dateModified
      hasPart(offset: 0) {
        name
      }
    }
  }

  # Request a specific section by name
  {
    node(name: { authority: "simple.wikipedia.org", pageName: "Banana", name: "Fruit" } ) {
      dateModified
      name
      unsafe
    }
  }

  # List sections related to a Wikidata keyword
  {
    nodes(keyword: "Q503") {
      name
      isPartOf { name }
      keywords {
        id
        salience
      }
    }
  }
      </pre>
    </v-container>
  </div>
</template>
