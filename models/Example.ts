import HatenaBookmark from '~/models/HatenaBookmark'
import Bookmark from '~/models/Bookmark'
import ReferringTweets from '~/models/Tweet'

export interface Example {
    Id: number
    Label: number
    Url: string
    FinalUrl: string
    Title: string
    Description: string
    OgDescription: string
    OgType: string
    OgImage: string
    Body: string
    Score: number
    IsNew: boolean
    StatusCode: number
    Favicon: string
    ErrorCount: number
    CreatedAt: string
    UpdatedAt: string
    ReferringTweets: ReferringTweets
    HatenaBookmark: HatenaBookmark
}

function truncate(str: string, length: number, omission: string): string {
    str = str ? str : '';
    if (str.length <= length) {
        return str;
    }
    else {
        return str.substring(0, length) + omission;
    }
}

export function getTitle(example: Example, length: number, omission: string): string {
    var title = example.Title ? example.Title : example.Url;
    return truncate(title, length, omission);
}

export function getDomain(example: Example): string {
    var url = example.FinalUrl;
    return url.replace('http://', '').replace('https://', '').split(/[/?#]/)[0];
}

export function getUserName(example: Example): string {
    var domain = getDomain(example);
    var url = example.FinalUrl;
    var paths = url.replace('http://', '').replace('https://', '').split(/[/?#]/);
    if (paths.length === 0) {
        return "";
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
        return "";
    }
}

export function getDescription(example: Example, length: number, omission: string) {
    var title = example.Title ? example.Title : example.Url;
    var body = example.Body ? example.Body : title;
    var desc = example.OgDescription ? example.OgDescription : (example.Description ? example.Description : body);
    return truncate(desc, length, omission);
}

export function filterBookmarksWithComment(example: Example): Bookmark[] {
    return example.HatenaBookmark.bookmarks.filter(function (b: Bookmark) {
        return b.comment !== "";
    });
}

export function getDescriptionForSearchEngine(example: Example) {
    var tweets = example.ReferringTweets.Tweets.map(t => "@" + t.ScreenName + "「" + t.FullText.substr(0, 100) + "...」").slice(0, 3);
    var bookmarks = filterBookmarksWithComment(example).map(b => "id:" + b.user + "「" + b.comment + "」").slice(0, 3);
    return tweets.join("\n") + bookmarks.join("\n");
}

export function getExampleUrl(example: Example) {
    return '/example/' + example.Id;
}