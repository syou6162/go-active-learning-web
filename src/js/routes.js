import RecentAddedExamples from './RecentAddedExamples.vue';
import List from './List.vue';

export default [
  {
    path: '/list/:listname',
    component: List
  },
  {
    path: '/recent-added-examples',
    component: RecentAddedExamples
  }
]
