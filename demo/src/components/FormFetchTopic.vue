<template>
  <div id="form-fetch-page">
    <v-form>
      <v-container>
        <v-row>
          <v-select
            v-model="topicName"
            label="Topic"
            :items=topiclist
            autocomplete
            @change="fetchPiecesForTopic"
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

    <v-container v-if="ready" class="topic-results">
      <div class="topic-result-single" v-for="result in results" :key="result.title">
        <h2>{{result.page}}</h2>
        <div v-html="result.content"></div>
      </div>
    </v-container>
  </div>
</template>

<script>
import axios from 'axios'
import scientists from '../data/scientists-topics.json'
import fruits from '../data/fruits-topics.json'

export default {
  name: 'FormFetchTopics',
  data: () => ({
    loading: false,
    ready: false,
    error: null,
    results: [],
    topicName: null,
    topiclist: [
      { text: 'Physics (Q1457258)', value: 'scientists|Physics (Q1457258)' },
      { text: 'Fruit (Q7134786)', value: 'fruits|Fruit (Q7134786)' },
      { text: 'Physical sciences (Q7153079)', value: 'scientists|Physical sciences (Q7153079)' },
      { text: 'Scholars and academics (Q7005672)', value: 'scientists|Scholars and academics (Q7005672)' },
      { text: 'Scientists (Q7043111)', value: 'scientists|Scientists (Q7043111)' },
      { text: 'Food and drink (Q5645580)', value: 'fruits|Food and drink (Q5645580)' },
      { text: 'Cuisine (Q9703849)', value: 'fruits|Cuisine (Q9703849)' },
      { text: 'Scientists (Q7043111)', value: 'scientists|Scientists (Q7043111)' }
    ]
  }),
  methods: {
    fetchPiecesForTopic() {
      const promises = []
      const topic = this.topicName.split('|')
      const partsForTopics = {
        fruits,
        scientists
      }
      const pieces = partsForTopics[topic[0]][topic[1]]
      const numPieces = Math.min(Object.keys(pieces).length, 8)
      let counter = 0
      this.loading = true
      this.ready = false
      this.results = []
      console.log('fetchPiecesForTopic', numPieces, pieces)
      // Fetch info per piece, limit at 8
      // eslint-disable-next-line no-unused-vars
      for (const [_, pieceData] of Object.entries(pieces)) {
        if (counter >= numPieces) {
          break
        }
        const p = this.fetch(pieceData.page, pieceData.part)
          .then(res => {
            console.log(`${pieceData.page} - ${pieceData.part}`, res)
            return {
              title: `${pieceData.page} - ${pieceData.part}`,
              page: pieceData.page,
              part: pieceData.part,
              content: res.unsafe
            }
          })

        promises.push(p)
        counter++
      }

      Promise.all(promises).then(results => {
        console.log('promises', results)
        this.results = results
        this.loading = false
        this.ready = true
      })
    },
    fetch(pageName, partName) {
      console.log('fetch', pageName, partName)
      return axios
        .post(
          'http://localhost:8080/query',
          {
            query: `{
  node(name: { authority: "simple.wikipedia.org", pageName: "${pageName}", name: "${partName}" } ) {
    name
    unsafe
  }
}`
          },
          {
            headers: { 'Content-Type': 'application/json' }
          }
        )
        .then(res => {
          if (res && res.data && res.data.data && (res.data.data.page || res.data.data.node)) {
            return res.data.data.node
          }
          return Promise.reject(new Error('Request not found'))
        })
        .catch(e => {
          console.log('err')
          this.error = e
          console.log(e)
        })
    }
  }
}
</script>

<style>
.topic-results {
  display: flex;
  flex-wrap: wrap;
}

.topic-result-single {
  margin: 0.5em;
  padding: 0.5em;
  border: 1px solid #ccc;
  border-radius: 0.5em;
  overflow-y: auto;
  height: 350px;
  max-width: 30%;
}

.topic-result-single h2 {
  color: #3f51b5;
}
</style>
