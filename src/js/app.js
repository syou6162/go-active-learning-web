import Vue from 'vue';
import BootstrapVue from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import Index from './Index.vue';
import RecentAddedExamples from './RecentAddedExamples.vue';

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
