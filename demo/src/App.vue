<template>
  <v-app id="phoenix">
    <v-navigation-drawer v-model="drawer" app clipped dark color="accent">
      <v-list subheader dense>
        <v-list-item link to="/">
          <v-list-item-action>
            <v-icon>fa-anchor</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Home</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item link to="/about">
          <v-list-item-action>
            <v-icon>fa-info-circle</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>About</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item link :href="sandboxLink">
          <v-list-item-action>
            <v-icon>fa-flask</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>API sandbox</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

        <v-divider></v-divider>
        <v-subheader inset>Examples</v-subheader>
        <v-list-item link to="/fetchbyname">
          <v-list-item-action>
            <v-icon>fa-file</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Page sections</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-list-item link to="/fetchbykeyword">
          <v-list-item-action>
            <v-icon>fa-file</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Collections</v-list-item-title>
          </v-list-item-content>
        </v-list-item>

      </v-list>
    </v-navigation-drawer>

    <v-app-bar
      app
      clipped-left
      color="primary"
      dark
    >
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Phoenix : Demo</v-toolbar-title>
    </v-app-bar>

    <v-main>
      <v-container>
            <router-view />
      </v-container>
    </v-main>

    <v-footer app>
      <span>&copy; {{ new Date().getFullYear() }}</span>
    </v-footer>
  </v-app>
</template>

<script>
export default {
  props: {
    source: String
  },

  data: () => ({
    drawer: null
  }),
  computed: {
    sandboxLink: function () {
      const prefix = window.location.hostname === 'wikimedia.github.io'
        ? '/phoenix' : ''

      return prefix + '/sandbox'
    }
  }

//   created () {
//     this.$vuetify.theme.dark = false
//   }
}
</script>
