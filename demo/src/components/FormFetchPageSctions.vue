<template>
  <div id="form-fetch-page">
    <v-form v-model="valid">
      <v-container>
        <v-row>
          <v-select
            v-model="pageName"
            label="Page name"
            :items=pages
            autocomplete
            @change="fetchParts"
            required
          ></v-select>
          <v-select
            v-model="partName"
            label="Parts"
            :items=selectPartsItems
            :disabled="!selectPartsItems.length"
            :loading="loading"
            @change="fetchContent"
            autocomplete
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
        <v-tab href="#tab-page">Page</v-tab>
        <v-tab href="#tab-payload">Payload</v-tab>
        <v-tab-item value="tab-page">
          <PageInfo v-if="result" :pagedata="result" />
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
import PageInfo from './PageInfo'
import axios from 'axios'

export default {
  name: 'FormFetchPageSctions',
  components: { PageInfo },
  data: () => ({
    error: null,
    payload: null,
    query: null,
    result: null,
    valid: false,
    loading: false,
    selectPartsItems: [],
    pageName: '',
    partName: '',
    pages: [
      { text: 'Philadelphia', value: 'Philadelphia' },
      // { text: 'Albert Einstein', value: 'Albert Einstein' },
      { text: 'Banana', value: 'Banana' },
      { text: 'Apple', value: 'Apple' }
    ]
  }),
  methods: {
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
          // console.log('fetchContent data', data)
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
          // 'http://ec2-3-133-13-197.us-east-2.compute.amazonaws.com:8080',
          '/graphql',
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
