<template>
  <div id="section-collection">
      <v-chip
        class="ma-2"
        dark
        label
        large
        color="pink"
        @click.stop="onKeywordInfoButtonClick(keyword)"
      >
        <v-avatar left>
          <v-icon>mdi-information</v-icon>
        </v-avatar>
        <strong>Keyword: {{keyword}}</strong>
      </v-chip>
    <v-container>
      <masonry
        :cols="{default: 3, 1000: 2, 600: 1}"
        :gutter="20"
        >
        <SectionBox
          v-for="sect in sections"
          :key="sect.id"
          :sectiondata="sect"
          :currKeyword=keyword
          @keywordClick="onSectionKeywordClicked"
          @pageInfoClick="onPageInfoClicked"
        >
          <div slot="content" v-html="sect.unsafe"></div>
        </SectionBox>
      </masonry>
    </v-container>

    <v-dialog
      class="wikidata-dialog"
      v-model="wikidataDialog"
      scrollable
      :fullscreen="$vuetify.breakpoint.smAndDown"
    >
      <v-card>
        <v-card-title>
          Wikidata item: {{currWikidataItem}}
        </v-card-title>
        <v-card-text class="wikidataDialog-text">
          <iframe :src="currWikidataUrl"></iframe>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="primary"
            text
            @click="wikidataDialog = false"
          >
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog
      class="page-dialog"
      v-model="pageViewDialog"
      scrollable
      :fullscreen="$vuetify.breakpoint.smAndDown"
    >
      <v-card>
        <v-card-title>
          Wikipedia page: {{viewPageTitle}}
        </v-card-title>
        <v-card-text class="wikidataDialog-text">
          <iframe :src="viewSimpleEnglishArticle"></iframe>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="primary"
            text
            @click="pageViewDialog = false"
          >
            Close
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

  </div>
</template>

<script>
import SectionBox from './SectionBox'

export default {
  name: 'SectionCollection',
  components: { SectionBox },
  props: {
    keyword: String,
    sections: Array
  },
  data: () => ({
    wikidataDialog: false,
    pageViewDialog: false,
    currWikidataItem: null,
    viewPageTitle: null
  }),
  computed: {
    wikidataUrl() {
      console.log('sections', this.sections)
      return `https://www.wikidata.org/wiki/${this.keyword}`
    },
    currWikidataUrl() {
      return `https://m.wikidata.org/wiki/${this.currWikidataItem}`
    },
    viewSimpleEnglishArticle () {
      const urlTitle = encodeURIComponent(this.viewPageTitle)
      return `https://simple.wikipedia.org/wiki/${urlTitle}`
    }
  },
  methods: {
    onSectionKeywordClicked(keyword) {
      // Bubble up
      this.$emit('keywordClick', keyword)
    },
    onPageInfoClicked(pageTitle) {
      this.viewPageTitle = pageTitle
      this.pageViewDialog = true
    },
    onKeywordInfoButtonClick(keyword) {
      this.currWikidataItem = keyword
      this.wikidataDialog = true
    }
  }
}
</script>

<style lang="less">
.wikidataDialog-text iframe {
  width: 100%;
  height: 500px;
}
</style>
