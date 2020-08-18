<template>
  <div id="form-fetch-page">
    <v-form v-model="valid">
      <v-container>
        <v-row>
          <v-text-field
            v-model="pageName"
            :rules="pageRules"
            label="Page name"
            required
          ></v-text-field>
          <v-btn
            color="secondary"
            :disbled="!valid"
            :loading="loading"
            @click="fetch"
            >Fetch</v-btn
          >
        </v-row>
      </v-container>
    </v-form>
    <v-alert v-show="error">{{ error }}</v-alert>
    <v-container v-if="result">
      <v-tabs>
        <v-tab href="#tab-page">Page</v-tab>
        <v-tab href="#tab-payload">Payload</v-tab>
        <v-tab-item value="tab-page">
          <PageInfo v-if="result" :pagedata="result" />
        </v-tab-item>
        <v-tab-item value="tab-payload">
          <v-container v-if="payload">
            <v-card>
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
  name: 'FormFetchByPageName',
  components: { PageInfo },
  data: () => ({
    error: null,
    payload: null,
    result: null,
    valid: false,
    loading: false,
    pageName: '',
    pageRules: [
      v => !!v || 'Page name is required',
      v => /^[^<>{}[\]|#]+$/.test(v) || 'Page name must be valid'
    ]
  }),
  methods: {
    fetch() {
      this.loading = true
      this.payload = null
      this.result = null
      this.error = null
      axios
        .post(
          'http://localhost:8080/query',
          {
            query:
              'query Page($name: String!) { pageByName(name: $name) { name dateModified hasPart about { key val } } }',
            // TODO: The name should come from the search input
            variables: { name: 'Foobar' }
          },
          {
            headers: { 'Content-Type': 'application/json' }
          }
        )
        .then(res => {
          const data = res.data.data.pageByName
          this.payload = JSON.stringify(res, null, 2)
          this.result = {
            title: data.name,
            modified: data.dateModified,
            parts: data.hasPart
          }
          console.log(res)
          console.log(this.result)
        })
        .catch(e => {
          this.error = e
          console.log(e)
        })
        .then(() => {
          // always trigger
          this.loading = false
        })
    }
  }
}
</script>
