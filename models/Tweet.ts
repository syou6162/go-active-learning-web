export default interface Tweet {
    ExampleId: number

    CreatedAt: string
    IdStr: string
    FullText: string
    FavoriteCount: number 
    RetweetCount: number 
    Lang: string

    ScreenName: string
    Name: string
    ProfileImageUrl: string
    Label: number
    Score: number
}

export default interface ReferringTweets {
    Tweets: Tweet[]
    Count: number
}