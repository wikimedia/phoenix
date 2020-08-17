<template>
  <div id="page-info">
    <h1>
      {{ pagedata.title }}
    </h1>
    <v-list two-line subheader>
      <v-subheader inset>Metadata</v-subheader>
      <v-list-item>
        <v-list-item-avatar>
          <v-icon>fa-calendar</v-icon>
        </v-list-item-avatar>
        <v-list-item-content>
          <v-list-item-title>Last modified</v-list-item-title>
          <v-list-item-subtitle v-text="humanReadable"></v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>

      <v-subheader inset>Parts</v-subheader>
      <v-list-item v-for="part in pagedata.parts" :key="part">
        <v-list-item-avatar>
          <v-icon>fa-puzzle-piece</v-icon>
        </v-list-item-avatar>
        <v-list-item-content>{{ part }}</v-list-item-content>
      </v-list-item>
    </v-list>
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
