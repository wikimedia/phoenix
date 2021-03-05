<template>
  <div class="home">
    <v-img contain alt="Phoenix image" src="../assets/phoenix.png" max-height="150px" />
    <h1 class="text-center">Welcome to Phoenix</h1>
    <p class="font-italic text-center">
      Exploring emerging patterns for building systems that manage and distribute content
    </p>
    <v-container>
      <p>
        <a href="https://github.com/wikimedia/phoenix" target="_blank">Phoenix</a> is an experimental
        service that consumes and structures Wikimedia content. It demonstrates how the information on
        a wiki page can be stored as a set of parts. By associating parts with relevant keywords, we can
        create relationships between pieces of knowledge across traditional page boundaries. This site
        lets you explore these relationships using content from
        <a href="https://simple.wikipedia.org" target="_blank">Simple English Wikipedia</a>
        and keywords from <a href="https://www.wikidata.org" target="_blank">Wikidata</a>.
      </p>
      <b>
        Phoenix is not a production-ready system.
        The API used in this demo should be considered unstable and may change without notice.
      </b>
      <br><br>
      <v-row>
        <v-col cols="12" md="6" lg="6" xl="6">
          <v-card>
            <v-card-title>Get only the information you need</v-card-title>
            <v-card-text>
              Fetch just the part of a page that you need, without having to request the entire page.
            </v-card-text>
            <v-card-actions>
              <v-btn
                link to="/fetchbyname"
                text
              >
                Try it out
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
        <v-col cols="12" md="6" lg="6" xl="6">
          <v-card>
            <v-card-title>Explore collections by keyword</v-card-title>
            <v-card-text>
              Make connections across pages by searching for sections related to Wikidata items.
            </v-card-text>
            <v-card-actions>
              <v-btn
                link to="/fetchbykeyword"
                text
              >
                Try it out
              </v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
// @ is an alias to /src

export default {
  name: 'Home'
  // components: {
  //   HelloWorld
  // }
}
</script>
