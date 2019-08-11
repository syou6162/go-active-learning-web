import { format, formatDistance, formatRelative, subDays, isAfter, toDate, parseISO } from 'date-fns'
import jaLocale from 'date-fns/locale/ja'
import Example from '~/models/Example'
import Bookmark from '~/models/Bookmark'

export function NewExample(e: Example, opts = {}): Example {
  var isNewDayThreshold = subDays(new Date(), opts["IsNewDayThreshold"] || 1);
  var createdAt = toDate(parseISO(e.CreatedAt));
  var updatedAt = toDate(parseISO(e.UpdatedAt));
  e.CreatedAt = format(createdAt, "yyyy/MM/dd HH:mm");
  e.UpdatedAt = format(updatedAt, "yyyy/MM/dd HH:mm");
  e.IsNew = isAfter(createdAt, isNewDayThreshold);
  if (e.HatenaBookmark) {
    e.HatenaBookmark.bookmarks = (e.HatenaBookmark.bookmarks || []).reverse();
  }
  return e;
}

export function getDomain(example: Example): string {
  var url = example.FinalUrl;
  return url.replace('http://','').replace('https://','').split(/[/?#]/)[0];
}

export function truncate(str: string, length: number, omission: string): string {
  str = str ? str : '';
  var l = length ? parseInt(length, 10) : 20;
  var omission = omission ? omission.toString() : '...';

  if (str.length <= l) {
    return str;
  }
  else {
    return str.substring(0, l) + omission;
  }
}

export function filterBookmarksWithComment(example: Example): Bookmark[] {
  return example.HatenaBookmark.bookmarks.filter(function(b: Bookmark) {
    return b.comment !== "";
  });
}
