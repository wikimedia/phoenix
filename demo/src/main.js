import Vue from 'vue'
import App from './App.vue'
import './registerServiceWorker'
import router from './router'
import vuetify from './plugins/vuetify'
import 'roboto-fontface/css/roboto/roboto-fontface.css'
import '@mdi/font/css/materialdesignicons.css'
import VueMasonry from 'vue-masonry-css'
const getGraphQLEndpoint = () => {
  switch (process.env.NODE_ENV) {
    case 'localgraphql':
      // Run the demo against a local run graphql service in localhost
      return 'http://localhost:8080/query'
    case 'production':
      // Run the demo from the netlify /graphql redirect
      return '/graphql'
    default:
      // Anything else, assume local development against aws endpoint
      return 'http://ec2-3-133-13-197.us-east-2.compute.amazonaws.com:8080'
  }
}

Vue.use(VueMasonry)

Vue.config.productionTip = false
Vue.prototype.GRAPHQL_ENDPOINT = getGraphQLEndpoint()

new Vue({
  router,
  vuetify,
  render: h => h(App)
}).$mount('#app')
