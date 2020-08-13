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
          <v-btn color="secondary" :disbled="!valid"
            :loading="loading"
            @click="fetch"
          >Fetch</v-btn>
        </v-row>
      </v-container>
    </v-form>
    <v-container v-show="result">{{result}}</v-container>
    <v-alert v-show="error">{{error}}</v-alert>
  </div>
 </template>

<script>
import axios from 'axios'

export default {
  name: 'FormFetchByPageName',
  data: () => ({
    error: null,
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
    fetch () {
      this.loading = true
      this.result = null
      this.error = null
      axios.post('http://localhost:8080/query', {
        query: 'query Page($name: String!) { pageByName(name: $name) { name dateModified hasPart about { key val } } }',
        // TODO: The name should come from the search input
        variables: { name: 'Foobar' }
      },
      {
        headers: { 'Content-Type': 'application/json' }
      })
        .then(res => {
          // TODO: Actually do something with this result
          this.result = res
          this.loading = false
          console.log(res)
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
