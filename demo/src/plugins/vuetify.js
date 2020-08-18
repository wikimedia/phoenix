import Vue from 'vue'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import '@fortawesome/fontawesome-free/css/all.css'

Vue.use(Vuetify, {
  iconfont: 'fa'
})

const theme = {
  primary: '#673ab7', // '#9c27b0',
  secondary: '#673ab7',
  accent: '#434343',
  error: '#f44336',
  warning: '#ff9800',
  info: '#00bcd4',
  success: '#4caf50'
}

export default new Vuetify({
  iconfont: 'fa',
  theme: {
    themes: {
      dark: theme,
      light: theme
      // light: {
      //   primary: '#9c27b0',
      //   secondary: '#424242',
      //   accent: '#82B1FF',
      //   error: '#FF5252',
      //   info: '#2196F3',
      //   success: '#4CAF50',
      //   warning: '#FFC107'
      // }
    }
  }
})
