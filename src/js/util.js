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

export function getTweetTitle(example) {
  var title = example.Title.replace(/\r?\n/g, ' ');
  var result = title.match(/^((.*? on Twitter)|(.*?さんのツイート)): \"(.*?)\"$/);
  return truncate(result[4], 200, '...');
}
 
export function filterBookmarksWithComment(example) {
  return example.HatenaBookmark.bookmarks.filter(function(b) {
    return b.comment !== "";
  });
}
