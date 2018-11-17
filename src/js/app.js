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

function getDomain(example) {
  var url = example.FinalUrl;
  return url.replace('http://','').replace('https://','').split(/[/?#]/)[0];
}

function truncate(str, length, omission) {
  str = str ? str : '';
  var length = length ? parseInt(length, 10) : 20;
  var omission = omission ? omission.toString() : '...';

  if (str.length <= length) {
    return str;
  }
  else {
    return str.substring(0, length) + omission;
  }
}

Vue.filter('getTitle', function(example, length, omission) {
  var title = example.Title ? example.Title : example.Url;
  return truncate(title, length, omission);
})

Vue.filter('getTweetTitle', function(example) {
  var title = example.Title.replace(/\r?\n/g, ' ');
  var result = title.match(/^((.*? on Twitter)|(.*?さんのツイート)): \"(.*?)\"$/);
  return truncate(result[4], 200, '...');
})

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

Vue.filter('getTwitterId', function(example) {
  var domain = getDomain(example);
  var url = example.FinalUrl;
  var paths = url.replace('http://','').replace('https://','').split(/[/?#]/);
  if (paths.length === 0) {
    return;
  } else if ('twitter.com' === domain) {
    return paths[1];
  } else {
    return;
  }
})

Vue.filter('getDescription', function(example, length, omission) {
  var title = example.Title ? example.Title : example.Url;
  var body = example.Body ? example.Body : title;
  var desc = example.OgDescription ? example.OgDescription : (example.Description ? example.Description : body);
  return truncate(desc, length, omission);
})

const app = new Vue({
  el: '#app',
  router
});
