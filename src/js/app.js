import Vue from 'vue';
import VueRouter from 'vue-router'
import BootstrapVue from 'bootstrap-vue'
import VueAnalytics from 'vue-analytics'
import VueLazyload from 'vue-lazyload'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import vueHeadful from 'vue-headful';

Vue.use(VueRouter)
Vue.use(BootstrapVue)
Vue.use(VueLazyload)
Vue.component('vue-headful', vueHeadful);

import routes from './routes'
const router = new VueRouter({
  mode: 'history',
  routes: routes
});

Vue.use(VueAnalytics, {
  id: 'UA-591180-8',
  router
})

import { truncate, getDomain, getTwitterId, getTweetTitle } from './util';

Vue.filter('getTitle', function(example, length, omission) {
  var title = example.Title ? example.Title : example.Url;
  return truncate(title, length, omission);
})

Vue.filter('getTweetTitle', getTweetTitle)
Vue.filter('getDomain', getDomain) 

Vue.filter('getUserName', function(example) {
  var domain = getDomain(example);
  var url = example.FinalUrl;
  var paths = url.replace('http://','').replace('https://','').split(/[/?#]/);
  if (paths.length === 0) {
    return;
  } else if ('twitter.com' === domain) {
    return '(@' + paths[1] + ')';
  } else if ('github.com' === domain) {
    return '(@' + paths[1] + ')';
  } else if ('qiita.com' === domain) {
    return '(@' + paths[1] + ')';
  } else if ('www.slideshare.net' === domain) {
    return '(id:' + paths[1] + ')';
  } else if ('speakerdeck.com' === domain) {
    return '(id:' + paths[1] + ')';
  } else {
    return;
  }
})

Vue.filter('getTwitterId', getTwitterId)

Vue.filter('getDescription', function(example, length, omission) {
  var title = example.Title ? example.Title : example.Url;
  var body = example.Body ? example.Body : title;
  var desc = example.OgDescription ? example.OgDescription : (example.Description ? example.Description : body);
  return truncate(desc, length, omission);
})

Vue.filter('getDescriptionForSearchEngine', function(example) {
  var tweets = example.ReferringTweets.map(t => "@" + t.ScreenName + "「" + t.FullText.substr(0, 100) + "...」").slice(0, 3);
  var bookmarks = example.HatenaBookmark.bookmarks.filter(function(b) {
    return b.comment !== "";
  }).map(b => "id:" + b.user + "「"+ b.comment + "」").slice(0, 3);
  return tweets.join("\n") + bookmarks.join("\n");
})

const app = new Vue({
  el: '#app',
  router
});
