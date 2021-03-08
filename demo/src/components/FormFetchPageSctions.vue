<template>
  <div id="form-fetch-page">
    <v-container>
      <p class="font-italic text-center">
        Search for a page in Simple English Wikipedia, and fetch a specific section directly
      </p>
      <v-form v-model="valid" autocomplete="off">
          <v-row>
            <v-col>
              <TitleAutocomplete @chosen="onTitleChosen" />
            </v-col>
            <v-col>
              <v-select
                v-model="partName"
                label="Sections"
                :items=selectPartsItems
                :disabled="!selectPartsItems.length"
                :loading="loading"
                @change="fetchContent"
                autocomplete
                required
              ></v-select>
            </v-col>
          </v-row>
      </v-form>
      <v-alert v-show="error">{{ error }}</v-alert>
      <v-progress-circular v-if="loading"
        indeterminate
        color="purple"
        class="text-center"
      ></v-progress-circular>

      <div v-if="result">
        <p class="font-italic text-center">
          View the requested section, or examine the network request and payload.
        </p>
        <v-tabs
          dark
          background-color="primary"
        >
          <v-tab href="#tab-page">Section</v-tab>
          <v-tab href="#tab-payload">Payload</v-tab>
          <v-tab-item value="tab-page">
            <PageInfo v-if="result" :pagedata="result" />
          </v-tab-item>
          <v-tab-item value="tab-payload">
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
          </v-tab-item>
        </v-tabs>
      </div>
    </v-container>
  </div>
</template>

<script>
import PageInfo from './PageInfo'
import TitleAutocomplete from './TitleAutocomplete'
import axios from 'axios'

export default {
  name: 'FormFetchPageSctions',
  components: { PageInfo, TitleAutocomplete },
  data: () => ({
    error: null,
    payload: null,
    query: null,
    result: null,
    valid: false,
    loading: false,
    selectPartsItems: [],
    pageName: '',
    partName: ''
  }),
  methods: {
    reset () {
      this.selectPartsItems = []
      this.partName = ''
      this.pageName = ''
      this.payload = null
      this.query = null
      this.result = null
      this.valid = false
    },
    onTitleChosen (chosenTitle) {
      console.log('chosenTitle', chosenTitle)
      if (!chosenTitle) {
        this.reset()
      } else {
        this.pageName = chosenTitle
        this.fetchParts()
      }
    },
    fetchParts() {
      const query = `{
        page(name: { authority: "simple.wikipedia.org", name: "${this.pageName}"} ) {
          name
          dateModified
          hasPart(offset: 0) {
            name
          }
        }
      }`
      this.selectPartsItems = []
      this.partName = ''
      this.loading = true
      return this.fetch(query)
        .then(res => {
          return res.data.data.page
        })
        .then(data => {
          this.selectPartsItems = data.hasPart.map(part => part.name)
          this.loading = false
        })
    },
    fetchContent() {
      const query = `{
        node(name: { authority: "simple.wikipedia.org", pageName: "${this.pageName}", name: "${this.partName}" } ) {
          dateModified
          name
          unsafe
        }
      }`
      // window.console.log('fetchContent query', query)
      return this.fetch(query)
        .then(res => {
          this.payload = JSON.stringify(res, null, 2)
          this.query = query
          return res.data.data.node
        })
        .then(data => {
          this.result = {
            title: this.pageName,
            modified: data.dateModified,
            part: {
              title: data.name,
              content: data.unsafe
            }
          }
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
          if (res && res.data && res.data.data && (res.data.data.page || res.data.data.node)) {
            return res
          }
          return Promise.reject(new Error('Request not found'))
        })
        .catch(e => {
          this.error = e
          window.console.log(e)
          this.loading = false
        })
    }
  }
}
</script>
