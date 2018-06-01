import Vue from 'vue';
import VueRouter from 'vue-router'
import BootstrapVue from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import Index from './Index.vue';
import RecentAddedExamples from './RecentAddedExamples.vue';

Vue.use(VueRouter)
Vue.use(BootstrapVue);

new Vue({
  el: '#index',
  components: {
    "index": Index
  }
});

new Vue({
  el: '#recent-added-examples',
  components: {
    "recent-added-examples": RecentAddedExamples 
  }
});

import routes from './routes'
const router = new VueRouter({
  routes: routes
});

const app = new Vue({
  el: '#app',
  router
});
