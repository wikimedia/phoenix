<template>
  <v-card class="sectionbox mt-4" elevation="6">
    <v-card-title color="blue-grey darken-2" class="sectionbox-title">
      {{ sectionDisplayName }}
      <v-spacer></v-spacer>
      <v-chip outlined small color="red" class="mx-5" title="Salience">
        <v-icon small left color="red">mdi-scale</v-icon>
        {{percentSalience(keywordData.salience)}}
      </v-chip>
    </v-card-title>
    <v-card-subtitle>
      From page: <span class="sectionbox-page-name">{{ pageTitle }}</span>
      <v-btn small text icon class="ml-2" @click.stop="onPageInfoClick()"><v-icon >mdi-information</v-icon></v-btn>
    </v-card-subtitle>
    <v-card-text class="sectionbox-content" @click.stop="onContentClick"><slot name="content"></slot></v-card-text>
    <v-divider class="mt-2" v-if="otherKeywords.length"></v-divider>
    <v-card-subtitle v-if="otherKeywords.length">Keywords:</v-card-subtitle>
    <v-card-subtitle v-if="!otherKeywords.length">No more keywords for this section.</v-card-subtitle>
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
        <v-icon small left>mdi-information</v-icon> {{getKeywordLabel(k.id)}} ({{percentSalience(k.salience)}})
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
    keywordLabels: Object,
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
        return k.id && k.id !== this.currKeyword
      })
    },
    pageTitle () {
      return this.sectiondata.isPartOf[0].name
    },
    sectionDisplayName () {
      if (this.sectiondata.name === '__intro__') {
        return `${this.pageTitle} (introduction)`
      }
      return this.sectiondata.name.replaceAll('_', ' ')
    }
  },
  methods: {
    getKeywordLabel(keyword) {
      return this.keywordLabels[keyword] || keyword
    },
    onKeywordClick(keyword) {
      this.$emit('keywordClick', keyword)
    },
    onPageInfoClick() {
      this.$emit('pageInfoClick', this.pageTitle)
    },
    fixedSalience(s) {
      return Number.parseFloat(s).toFixed(2)
    },
    percentSalience(s) {
      const percent = Number(s) * 100
      return (Number.parseFloat(percent).toFixed(2)) + '%'
    },
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
}
</script>

<style lang="less">
.sectionbox-content {
  max-height: 300px;
  overflow-y: auto;
  overflow-x: hidden;
}
</style>
