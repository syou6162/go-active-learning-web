var moment = require('moment-timezone');

export function NewExample(e, opts = {}) {
  var isNewDayThreshold = moment().add(-1 * (opts["IsNewDayThreshold"] || 1), "days");
  var createdAt = moment(e.CreatedAt);
  var updatedAt = moment(e.UpdatedAt);
  e.CreatedAt = createdAt;
  e.UpdatedAt = updatedAt;
  e.IsNew = createdAt.isAfter(isNewDayThreshold);
  e.HatenaBookmark.bookmarks = (e.HatenaBookmark.bookmarks || []).reverse();
  return e;
}

export function getDomain(example) {
  var url = example.FinalUrl;
  return url.replace('http://','').replace('https://','').split(/[/?#]/)[0];
}

export function truncate(str, length, omission) {
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

export function filterBookmarksWithComment(example) {
  return example.HatenaBookmark.bookmarks.filter(function(b) {
    return b.comment !== "";
  });
}

import Amplify, {
  Auth,
  API,
} from 'aws-amplify';

Amplify.Logger.LOG_LEVEL = 'DEBUG';
Amplify.configure({
  Auth: {
    region: 'us-east-1',
    identityPoolId: 'us-east-1:a61e4d83-cd2a-4bd3-8995-19f9e2726f76',
    userPoolId: 'us-east-1_4Jft7xFAo',
    userPoolWebClientId: '4ia5ifrn456rqg4vr6dqfh7e68',
  },
});

export function signUp(username, password) {
  return Auth.signUp({
    username,
    password,
  });
}

export function signIn(username, password) {
  return Auth.signIn(username, password);
}

export function signOut() {
  return Auth.signOut();
}

export function IsAdmin() {
  return Auth.currentSession().then(_ => true).catch(_ => false);
}
