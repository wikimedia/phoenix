<template>
  <div id="form-fetch-keyword">
    <v-container>
      <v-form v-model="valid">
          <v-row>
            <v-select
              v-model="keyword"
              label="Wikidata keyword"
              :items=keywordItems
              autocomplete
              @change="fetchKeywordSections"
              required
            ></v-select>
          </v-row>
      </v-form>
      <v-alert v-show="error">{{ error }}</v-alert>
      <v-progress-circular v-if="loading"
        indeterminate
        color="purple"
      ></v-progress-circular>
      <p v-if="result" class="font-italic text-center">
        These are the sections related to the requested keyword, with their salience percentage.<br />You can examine the network request and payload in the 'payload' tab.
      </p>

      <v-tabs
        v-if="result"
        dark
        background-color="primary"
      >
        <v-tab href="#tab-section">Sections</v-tab>
        <v-tab href="#tab-payload">Payload</v-tab>
        <v-tab-item value="tab-section">
          <SectionCollection v-if="result" :keyword="keyword" :sections="result.sections" />
        </v-tab-item>
        <v-tab-item value="tab-payload">
            <v-card
              v-if="payload"
              class="mx-auto"
              outlined
            >
              <v-card-title class="headline">Query</v-card-title>
              <v-card-text>
                <pre>{{ query }}</pre>
              </v-card-text>
            </v-card>
            <br />
            <v-card
              class="mx-auto"
              outlined
            >
              <v-card-title class="headline">Payload</v-card-title>
              <v-card-text>
                <pre>{{ payload }}</pre>
              </v-card-text>
            </v-card>
        </v-tab-item>
      </v-tabs>
    </v-container>
  </div>
</template>

<script>
import axios from 'axios'
import SectionCollection from './SectionCollection'

export default {
  name: 'FormFetchKeyword',
  components: { SectionCollection },
  data: () => ({
    error: null,
    valid: false,
    loading: false,
    result: null,
    keyword: '',
    keywordItems: [
      { text: 'Q503 (Banana)', value: 'Q503' },
      { text: 'Q89 (Apple)', value: 'Q89' },
      { text: 'Q1458083 (Science)', value: 'Q1458083' },
      { text: 'Q413 (Physics)', value: 'Q413' },
      { text: 'Q7191 (Nobel prize)', value: 'Q7191' },
      { text: 'Q7446056 (Outer space)', value: 'Q7446056' },
      { text: 'Q6508 (Astronomy)', value: 'Q6508' },
      { text: 'Q1970530 (Ecology)', value: 'Q1970530' },
      { text: 'Q4057308 (Materials)', value: 'Q4057308' }
    ]
  }),
  methods: {
    fetchKeywordSections() {
      const query = `{
  nodes(keyword: "${this.keyword}") {
    name
    id
    isPartOf {
      id
      name
    }
    unsafe
    keywords {
      id
      salience
    }
  }
}`
      this.loading = true
      console.log('fetchKeywordSections query', query)
      return this.fetch(query)
        .then(res => {
          this.payload = JSON.stringify(res.data.data, null, 2)
          this.query = query
          return res.data.data.nodes
        })
        .then(data => {
          const normalizedID = (e) => {
            return e.isPartOf[0].name.toLowerCase().replaceAll('_', ' ') + '|' +
              e.name.toLowerCase().replaceAll('_', ' ')
          }

          // FIXME: Normalize; we get some dupes in the shape of 'Related pages' and 'Related_pages'
          // from the same page (those are the same section) from the search; while that's
          // being looked into, we will normalize this in the frontend and remove these
          // as dupes
          const seenKeys = {}
          data = data.filter(entry => {
            if (seenKeys[normalizedID(entry)]) {
              return false
            } else {
              seenKeys[normalizedID(entry)] = true
              return true
            }
          })

          this.result = {
            sections: data
          }
        })
        .then(data => {
          this.loading = false
        })
    },
    fetch(query) {
      this.payload = null
      this.result = null
      this.error = null
      return axios
        .post(
          this.GRAPHQL_ENDPOINT,
          { query },
          {
            headers: { 'Content-Type': 'application/json' }
          }
        )
        .then(res => {
          console.log('res', res)
          if (res && res.data && res.data.data && (res.data.data.nodes)) {
            return res
          }
          return Promise.reject(new Error('Request not found'))
        })
        .catch(e => {
          this.error = e
          console.log(e)
          this.loading = false
        })
    }
  }
}
</script>
