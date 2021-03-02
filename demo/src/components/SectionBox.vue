<template>
  <v-card class="sectionbox mt-4" elevation="6">
    <v-card-title color="blue-grey darken-2" class="sectionbox-title">
      {{sectiondata.name}}
      <v-spacer></v-spacer>
      <v-chip outlined small color="red" class="mx-5" title="Salience">
        <v-icon small left color="red">mdi-scale</v-icon>
        {{fixedSalience(keywordData.salience)}}
      </v-chip>
    </v-card-title>
    <v-card-subtitle>From page: <span class="sectionbox-page-name">{{ pageTitle }}</span></v-card-subtitle>
    <v-card-text class="sectionbox-content"><slot name="content"></slot></v-card-text>
    <v-divider class="mt-2" v-if="otherKeywords.length"></v-divider>
    <v-card-subtitle>Keywords:</v-card-subtitle>
    <v-card-text v-if="otherKeywords.length">
      <v-chip
        color="secondary"
        class="mx-1 my-1"
        small
        label
        v-for="k in otherKeywords"
        :key="k.id"
        @click.stop="onKeywordClick(k.id)"
      >
        <v-icon small left>mdi-information</v-icon> {{k.id}} ({{percentSalience(k.salience)}})
      </v-chip>
    </v-card-text>
  </v-card>
</template>

<script>
export default {
  name: 'SectionBox',
  props: {
    currKeyword: String,
    sectiondata: Object,
    isPartOf: String
  },
  data: () => ({
  }),
  computed: {
    keywordData() {
      return this.sectiondata.keywords.filter(k => {
        return k.id === this.currKeyword
      })[0]
    },
    otherKeywords() {
      return this.sectiondata.keywords.filter(k => {
        return k.id !== this.currKeyword
      })
    },
    pageTitle () {
      return this.sectiondata.isPartOf[0].name
    }
  },
  methods: {
    onKeywordClick(keyword) {
      this.$emit('keywordClick', keyword)
    },
    fixedSalience(s) {
      return Number.parseFloat(s).toFixed(2)
    },
    percentSalience(s) {
      const percent = Number(s) * 100
      return (Number.parseFloat(percent).toFixed(2)) + '%'
    }
  }
}
</script>

<style lang="less">
.sectionbox-content {
  max-height: 300px;
  overflow-y: auto;
  overflow-x: hidden;
}
</style>
