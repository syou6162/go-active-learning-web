export default interface Tweet {
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