import Vue from 'vue'
import VueRouter from 'vue-router'
import Home from '../views/Home.vue'
import FetchPart from '../views/FetchPart.vue'
import FetchTopic from '../views/FetchTopic.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/fetchbyname',
    name: 'FetchPart',
    component: FetchPart
  },
  {
    path: '/fetchbytopic',
    name: 'FetchTopic',
    component: FetchTopic
  },
  {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () =>
      import(/* webpackChunkName: "about" */ '../views/About.vue')
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
