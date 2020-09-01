import { format, subDays, isAfter, toDate, parseISO } from 'date-fns'
import { Example } from '~/models/Example'
import Tweet from '~/models/Tweet'

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
  if (e.ReferringTweets) {
    e.ReferringTweets.Tweets = (e.ReferringTweets.Tweets || []).map(function(t: Tweet) {
      t.CreatedAt = format(parseISO(t.CreatedAt), "yyyy/MM/dd HH:mm");
      return t;
    });
  }
  return e;
}