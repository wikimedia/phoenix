<template>
  <div id="form-fetch-keyword">
    <v-container>
      <v-form v-model="valid">
          <v-row>
            <v-select
              v-model="keyword"
              label="Choose a Wikidata keyword"
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
      <p v-if="result">
        These are the sections related to the requested keyword, with their salience percentage:
        <v-chip outlined small color="red" title="Salience">
          <v-icon small left color="red">mdi-scale</v-icon>
          95%
        </v-chip>

        <ul>
          <li>You can examine the network request and payload in the 'payload' tab.</li>
          <li>You can click any of the keywords in each section to fetch sections for that keyword.</li>
        </ul>
      </p>

      <v-tabs
        v-if="result"
        dark
        background-color="primary"
      >
        <v-tab href="#tab-section">Sections</v-tab>
        <v-tab href="#tab-payload">Payload</v-tab>
        <v-tab-item value="tab-section">
          <SectionCollection v-if="result" :keyword="keyword" :sections="result.sections" :keywordLabels="result.keywordLabels" @keywordClick="onKeywordClick" />
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
import WikidataFetcher from '../tools/WikidataFetcher'
import SectionCollection from './SectionCollection'

export default {
  name: 'FormFetchKeyword',
  components: { SectionCollection },
  data: () => ({
    wikidataFetcher: null,
    error: null,
    valid: false,
    loading: false,
    result: null,
    keyword: '',
    keywordItems: [
      { text: 'Q2934 (Goat)', value: 'Q2934' },
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
  created() {
    this.wikidataFetcher = new WikidataFetcher('simple', 'en')
  },
  methods: {
    async onKeywordClick (keyword) {
      this.loading = true
      let keywordModel = this.keywordItems.filter(data => {
        return data.value === keyword
      })[0]
      // Add it to the keywordItems if it's not already ther
      if (!keywordModel) {
        keywordModel = {
          value: keyword,
          text: keyword
        }
        // Try to get the human-readable name
        await this.fetchWikidataItemLabels(keyword)
          .then(items => {
            if (items[keyword]) {
              keywordModel.text = `${keyword} (${items[keyword]})`
            }
          })
          .catch(err => {
            // If this failed, skip it
            console.log('Wikidata label fetch failed', err)
          })
        // Keyword isn't already in the dropdown, add it:
        this.keywordItems.push(keywordModel)
      }

      // Select the keyword
      this.keyword = keywordModel.value
      // Trigger another lookup
      this.fetchKeywordSections()
    },
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
      return this.fetch(query)
        .then(res => {
          this.payload = JSON.stringify(res.data.data, null, 2)
          this.query = query
          return res.data.data.nodes
        })
        .then(async data => {
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

          // Fetch wikidata keyword names
          const allKeywords = []
          // get all keywords
          data.forEach(section => {
            section.keywords.forEach(key => {
              if (allKeywords.indexOf(key.id) === -1) {
                allKeywords.push(key.id)
              }
            })
          })
          // Get all names resolved
          return this.wikidataFetcher.fetchWikidataItemLabels(allKeywords)
            .then(keyNames => {
              this.result = {
                sections: data,
                keywordLabels: keyNames
              }
            })
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
    },
    fetchWikidataItemLabels (id) {
      return this.wikidataFetcher.fetchWikidataItemLabels(id)
    }
  }
}
</script>
