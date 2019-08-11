import Bookmark from '~/models/Bookmark'

export default interface HatenaBookmark {
    title: string
    bookmarks: Bookmark[]
    screenshot: string
    entry_url: string
    count: number
    url: string
    eid: string
}