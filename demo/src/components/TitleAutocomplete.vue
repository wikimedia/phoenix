<template>
  <form autocomplete="off">
    <v-autocomplete
        v-model="chosenTitle"
        :items="suggestedTitles"
        :loading="loadingSuggestedTitles"
        :search-input.sync="search"
        hide-no-data
        item-text="name"
        item-value="id"
        label="Title"
        prepend-icon="mdi-article"
        placeholder="Start typing to Search"
        clearable
        return-object
    ></v-autocomplete>
    <span v-if="error">{{ error }}</span>
  </form>
</template>

<script>
import axios from 'axios'
import debounce from 'debounce'

export default {
  name: 'TitleAutocomplete',
  data() {
    return {
      chosenTitle: {},
      suggestedTitles: [],
      search: '',
      error: '',
      loadingSuggestedTitles: false
    }
  },
  methods: {
    onSubmit: () => {
      // console.log(this)
      // this.$emit('chosen', this.chosenTitle)
    },
    makeSearch: async (value, self) => {
      self.error = ''

      // Empty value
      if (!value) {
        self.suggestedTitles = []
        self.chosenTitle = {}
      }

      // Still loading
      if (self.loadingSuggestedTitles) {
        return
      }
      self.loadingSuggestedTitles = true
      await axios
        .get(
          'https://simple.wikipedia.org/w/api.php',
          {
            params: {
              action: 'query',
              prop: ['info', 'pageprops', 'description'],
              generator: 'prefixsearch',
              gpssearch: self.search,
              gpsnamespace: 0, // Main namespace
              gpslimit: 5,
              ppprop: 'disambiguation',
              redirects: true,
              origin: '*',
              format: 'json'
            }
          },
          {
            headers: {
              'Content-Type': 'application/json; charset=UTF-8',
              origin: '*'
            }
          }
        )
        .then(response => {
          return Object.values(response.data.query.pages)
        })
        .then(pages => {
          self.suggestedTitles = pages.map(data => {
            return {
              name: data.title,
              id: data.title
            }
          })
        })
        .catch(error => {
          self.error = 'Error fetching page suggestions. Please try again.'
          console.log('Error fetching title suggestions', error)
        })
        .finally(() => (self.loadingSuggestedTitles = false))
    }
  },
  watch: {
    search (value) {
      if (!value) {
        this.suggestedTitles = []
        return
      }

      // Debounce
      debounce(this.makeSearch, 200)(value, this)
    },
    chosenTitle (value) {
      this.$emit('chosen', this.chosenTitle && this.chosenTitle.id)
    }
  }
}
</script>
