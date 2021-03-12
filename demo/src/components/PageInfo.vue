<template>
  <div id="page-info">
    <h1>
      {{ pagedata.title }} - {{ pagedata.part.title }}
    </h1>
    <v-list two-line subheader>
      <v-list-item>
        <v-list-item-avatar>
          <v-icon>fa-calendar</v-icon>
        </v-list-item-avatar>
        <v-list-item-content>
          <v-list-item-title>Last modified</v-list-item-title>
          <v-list-item-subtitle v-text="humanReadable"></v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>
    </v-list>

    <v-card
    class="mx-auto"
    max-width="80%"
    outlined
    >
    <v-card-text @click.stop="onContentClick" v-html="pagedata.part.content"></v-card-text>
    </v-card>
  </div>
</template>

<script>
import moment from 'moment'

export default {
  name: 'PageInfo',
  props: { pagedata: Object },
  computed: {
    humanReadable() {
      return moment(this.pagedata.modified).format('MMMM Do YYYY, h:mm:ss a')
    }
  },
  methods: {
    onContentClick(e) {
      // HACK: We're getting internal links from the content; this is a bit of a
      // cheat to get internal links to open in Simple English WP
      if (e.target.tagName.toLowerCase() === 'a') {
        let href = e.target.getAttribute('href')
        if (e.target.getAttribute('rel') === 'mw:WikiLink') {
          href = href.replace(/^(\.\/)/, '')
          e.target.href = `http://simple.wikipedia.com/wiki/${href}`
        }
        // Change the target to open in new window
        e.target.target = '_blank'

        // Continue the link clicking...
      }
    }
  }
  // data: () => ({})
}
</script>
