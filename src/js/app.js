import Vue from 'vue'
import Index from './Index.vue';
import RecentAddedExamples from './RecentAddedExamples.vue';

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
