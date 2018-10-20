import Vue from 'vue';
import VueRouter from 'vue-router'
import BootstrapVue from 'bootstrap-vue'
import VueAnalytics from 'vue-analytics'
import VueLazyload from 'vue-lazyload'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

Vue.use(VueRouter)
Vue.use(BootstrapVue)
Vue.use(VueLazyload)

import routes from './routes'
const router = new VueRouter({
  routes: routes
});

Vue.use(VueAnalytics, {
  id: 'UA-591180-8',
  router
})

const app = new Vue({
  el: '#app',
  router
});
