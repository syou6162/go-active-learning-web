import RecentAddedExamples from './RecentAddedExamples.vue';
import List from './List.vue';

export default [
  {
    path: '/',
    redirect: '/list/general'
  },
  {
    path: '/list/:listname',
    component: List
  },
  {
    path: '/recent-added-examples',
    component: RecentAddedExamples
  }
]
