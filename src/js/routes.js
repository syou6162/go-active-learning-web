import RecentAddedExamples from './RecentAddedExamples.vue';
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
    path: '/recent-added-examples',
    component: RecentAddedExamples
  }
]
