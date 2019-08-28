import Vue from 'vue';
import { truncate, getDomain, filterBookmarksWithComment } from '~/plugins/util';
import Example from '~/models/Example'

Vue.filter('getTitle', function(example: Example, length: number, omission: string) {
  var title = example.Title ? example.Title : example.Url;
  return truncate(title, length, omission);
})

Vue.filter('getDomain', getDomain) 

Vue.filter('getUserName', function(example: Example) {
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

Vue.filter('getDescription', function(example: Example, length: number, omission: string) {
  var title = example.Title ? example.Title : example.Url;
  var body = example.Body ? example.Body : title;
  var desc = example.OgDescription ? example.OgDescription : (example.Description ? example.Description : body);
  return truncate(desc, length, omission);
})

Vue.filter('getDescriptionForSearchEngine', function(example: Example) {
  var tweets = example.ReferringTweets.map(t => "@" + t.ScreenName + "「" + t.FullText.substr(0, 100) + "...」").slice(0, 3);
  var bookmarks = filterBookmarksWithComment(example).map(b => "id:" + b.user + "「"+ b.comment + "」").slice(0, 3);
  return tweets.join("\n") + bookmarks.join("\n");
})

Vue.filter('getExampleUrl', function(example: Example) {
  return '/example/' + example.Id;
})
