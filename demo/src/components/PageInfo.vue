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
    max-width="350"
    outlined
    >
    <v-card-text v-html="pagedata.part.content"></v-card-text>
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
  }
  // data: () => ({})
}
</script>
