import { createLocalVue, mount } from '@vue/test-utils'
import { ButtonPlugin, BCard, BCardBody, BCardFooter } from 'bootstrap-vue'
import AdminAnnotateTweets from './AdminAnnotateTweets.vue'
import VueLazyLoad from 'vue-lazyload'

const localVue = createLocalVue();
localVue.use(VueLazyLoad);
localVue.use(ButtonPlugin)
localVue.component('b-card', BCard)
localVue.component('b-card-body', BCardBody)
localVue.component('b-card-footer', BCardFooter)

const tweets = [
  {
    ExampleId: 21805423,

    CreatedAt: "2020-08-26T03:13:06Z",
    IdStr: "1298458422581448704",
    FullText: "【イベント情報】先月開催された「Embulk &amp; DigDag Online Meetup 2020」のレポート記事を公開しました。\nEmbulkとDigdagをプロダクション環境で利用している事例や、Embulk開発者による発表が盛り沢山です！是非ご一読下さい。\nhttps://t.co/gdtRuoyJtY",
    FavoriteCount: 6,
    RetweetCount: 4,
    Lang: "ja",

    ScreenName: "trocco_jp",
    Name: "データ分析基盤の総合支援SaaS「trocco(トロッコ)」",
    ProfileImageUrl: "https://pbs.twimg.com/profile_images/1051692644311556096/WmDcCj4Z_normal.jpg",
    Label: 1,
    Score: 0.4894515842007917
  }
]

describe('AdminAnnotateTweets', () => {
  test('管理者用のtweetアノテーション画面が表示される', () => {
    const wrapper = mount(AdminAnnotateTweets, {
      localVue,
      propsData: {
        example: {
          Id: 11111,
          Label: 0,
          Url: "https://blog.trocco.io/dev-articles/embulk-digdag-meetup-2020",
          FinalUrl: "https://blog.trocco.io/dev-articles/embulk-digdag-meetup-2020",
          Title: "Embulk & Digdag Online Meetup 2020【イベントレポート】 | trocco(トロッコ)",
          Description: "",
        },
        tweets: tweets,
        isAdmin: true
      },
    })
    expect(wrapper.exists()).toBeTruthy()
    expect(wrapper.element).toMatchSnapshot()
  })
})
