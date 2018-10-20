import RecentAddedExamples from './RecentAddedExamples.vue';
import SearchExamples from './SearchExamples.vue';
import ListExample from './ListExample.vue';

export default [
  {
    path: '/',
    redirect: '/list/general'
  },
  {
    path: '/list/:listname',
    component: ListExample
  },
  {
    path: '/search',
    component: SearchExamples
  },
  {
    path: '/recent-added-examples',
    component: RecentAddedExamples
  }
]
