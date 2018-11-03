import RecentAddedExamples from './RecentAddedExamples.vue';
import SearchExamples from './SearchExamples.vue';
import ListExample from './ListExample.vue';
import ExampleDetail from './ExampleDetail.vue';

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
  },
  {
    path: '/example/:url',
    component: ExampleDetail 
  }
]
