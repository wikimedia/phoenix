<template>
  <div id="form-fetch-keyword">
    <v-form v-model="valid">
      <v-container>
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
      </v-container>
    </v-form>
    <v-alert v-show="error">{{ error }}</v-alert>
    <v-progress-circular v-if="loading"
      indeterminate
      color="purple"
    ></v-progress-circular>

    <v-container v-if="result">
      <v-tabs>
        <v-tab href="#tab-section">Sections</v-tab>
        <v-tab href="#tab-payload">Payload</v-tab>
        <v-tab-item value="tab-section">
          <SectionCollection v-if="result" :keyword="keyword" :sections="result.sections" />
        </v-tab-item>
        <v-tab-item value="tab-payload">
          <v-container v-if="payload">
            <v-card
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
          </v-container>
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
  nodes(relatedTo: "${this.keyword}") {
    name
    isPartOf
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
          // 'http://ec2-3-133-13-197.us-east-2.compute.amazonaws.com:8080',
          '/graphql',
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
