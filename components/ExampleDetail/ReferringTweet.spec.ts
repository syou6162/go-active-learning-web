import { createLocalVue, shallowMount } from '@vue/test-utils'
import { BListGroupItem } from 'bootstrap-vue'
import ReferringTweet from '~/components/ExampleDetail/ReferringTweet.vue'

const localVue = createLocalVue();
localVue.component('b-list-group-item', BListGroupItem)

const tweet = {
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

describe('ReferringTweet', () => {
  test('個別事例ページの中のtweetが正しく表示できる', () => {
    const wrapper = shallowMount(ReferringTweet, {
      localVue,
      propsData: {
        tweet: tweet,
        isAdmin: false
      }
    })
    expect(wrapper.exists()).toBeTruthy()
    expect(wrapper.find('.tweet-item').exists()).toBe(true)
    expect(wrapper.find('.tweet-retweet-count').text()).toBe(`${tweet.RetweetCount} RT`)
    expect(wrapper.find('.tweet-favorite-count').text()).toBe(`${tweet.FavoriteCount} Fav`)

    expect(wrapper.element).toMatchSnapshot()
  })

  test('adminの場合はアノテーションボタンが表示される', () => {
    const wrapper = shallowMount(ReferringTweet, {
      localVue,
      propsData: {
        tweet: tweet,
        isAdmin: true
      }
    })
    expect(wrapper.exists()).toBeTruthy()
    expect(wrapper.find('tweet-annotate-buttons-stub').exists()).toBe(true)

    expect(wrapper.element).toMatchSnapshot()
  })
})