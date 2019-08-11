import HatenaBookmark from '~/models/HatenaBookmark'
import Tweet from '~/models/Tweet'

export default interface Example {
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
    ReferringTweets: Tweet[]
    HatenaBookmark: HatenaBookmark
}