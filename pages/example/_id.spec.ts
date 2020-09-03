import { createLocalVue, shallowMount } from '@vue/test-utils'
import { Auth } from 'aws-amplify';
import { BButton, BCard, BCardBody, BCardFooter, BListGroup, BListGroupItem } from 'bootstrap-vue'
import VueLazyLoad from 'vue-lazyload'
import VueMeta from 'vue-meta'
import ExamplePage from './_id.vue'

const localVue = createLocalVue();
localVue.use(VueLazyLoad);
localVue.component('b-button', BButton)
localVue.component('b-card', BCard)
localVue.component('b-card-body', BCardBody)
localVue.component('b-card-footer', BCardFooter)
localVue.component('b-list-group', BListGroup)
localVue.component('b-list-group-item', BListGroupItem)
// ref: https://stackoverflow.com/questions/59964001/how-to-test-head-in-nuxt-js-using-jest
localVue.use(VueMeta, { keyName: 'head' })

const hatenaBookmark = {
  title: "Embulk & Digdag Online Meetup 2020【イベントレポート】 | trocco(トロッコ)",
  screenshot: "https://b.st-hatena.com/images/v4/public/common/noimage.png?version=8358883de4e4eced6d5c2b49020f40bbbcfdca71",
  entry_url: "https://b.hatena.ne.jp/entry/s/blog.trocco.io/dev-articles/embulk-digdag-meetup-2020",
  bookmarks: [
    {
      "comment": "",
      "user": "machupicchubeta",
      "tags": [],
      "timestamp": "2020/08/27 00:56"
    },
    {
      "tags": [],
      "timestamp": "2020/08/26 19:53",
      "user": "hoppie",
      "comment": ""
    },
    {
      "comment": "“統合”",
      "user": "polamjag",
      "timestamp": "2020/08/26 14:29",
      "tags": []
    },
    {
      "comment": "",
      "user": "yukiyan_w",
      "tags": [
        "digdag",
        "embulk"
      ],
      "timestamp": "2020/08/26 12:15"
    }
  ],
  count: 5,
  url: "https://blog.trocco.io/dev-articles/embulk-digdag-meetup-2020",
  eid: "4690501516113184834"
}

const referringTweets = {
  Tweets: [
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
    },
    {
      ExampleId: 21805423,

      CreatedAt: "2020-08-26T03:13:06Z",
      IdStr: "1111111111",
      FullText: "適当な負例です",
      FavoriteCount: 6,
      RetweetCount: 4,
      Lang: "ja",

      ScreenName: "syou6162",
      Name: "データ分析基盤の総合支援SaaS「trocco(トロッコ)」",
      ProfileImageUrl: "https://pbs.twimg.com/profile_images/1051692644311556096/WmDcCj4Z_normal.jpg",
      Label: 0,
      Score: -2.0
    },
  ], 
  Count: 2
}

const example = {
  Id: 11111,
  Label: 0,
  Url: "https://blog.trocco.io/dev-articles/embulk-digdag-meetup-2020",
  FinalUrl: "https://blog.trocco.io/dev-articles/embulk-digdag-meetup-2020",
  Title: "Embulk & Digdag Online Meetup 2020【イベントレポート】 | trocco(トロッコ)",
  Description: "",
  HatenaBookmark: hatenaBookmark,
  ReferringTweets: referringTweets,
}

const keywords = ["embulk", "digdag"]

describe('ExamplePage', () => {
  jest.spyOn(Auth, 'currentAuthenticatedUser').mockReturnValue(
    Promise.resolve({
      username: "myuser",
      attributes: {
        email: "test@test.com",
        email_verified: true,
        phone_number: "08012345678"
        // ...other attoributes
      }
      // ...other parameters
    })
  );

  const wrapper = shallowMount(ExamplePage, {
    localVue,
  });

  test('メタタグが正しく出ている', async () => {
    await wrapper.setData({
      title: example.Title,
      example: example,
      keywords: keywords,
    })

    expect(wrapper.exists()).toBeTruthy()
    // ref: https://stackoverflow.com/questions/59964001/how-to-test-head-in-nuxt-js-using-jest
    // @ts-ignore
    expect(wrapper.vm.$metaInfo.meta.find((item) => item.name === 'robots').content).toBe('index, follow')
  })
  
  test('負例は検索エンジンから無視してもらう', async () => {
    await wrapper.setData({
      title: example.Title,
      example: {
        Label: -1
      },
      keywords: keywords,
    })

    expect(wrapper.exists()).toBeTruthy()
    // ref: https://stackoverflow.com/questions/59964001/how-to-test-head-in-nuxt-js-using-jest
    // @ts-ignore
    expect(wrapper.vm.$metaInfo.meta.find((item) => item.name === 'robots').content).toBe('noindex, nofollow')
    // @ts-ignore
    expect(wrapper.vm.$metaInfo.meta.find((item) => item.name === 'description').content).toBeDefined()
    // @ts-ignore
    expect(wrapper.vm.$metaInfo.meta.find((item) => item.name === 'og:title').content).toBeDefined()
  })

  test('Twitterの項目が出ている', async () => {
    await wrapper.setData({
      title: example.Title,
      example: example,
      keywords: keywords,
    })

    expect(wrapper.exists()).toBeTruthy()
    expect(wrapper.find('div#referring-tweet').exists()).toBe(true)
    expect(wrapper.element).toMatchSnapshot()
  })
  
  test('はてなブックマークの項目が出ている', async () => {
    await wrapper.setData({
      title: example.Title,
      example: example,
      keywords: keywords,
    })

    expect(wrapper.exists()).toBeTruthy()
    expect(wrapper.find('div#bookmark-comments').exists()).toBe(true)
    expect(wrapper.element).toMatchSnapshot()
  })
})

describe('ExamplePage#Amplify', () => {
  test('amplifyを通っていないときはアノテーションボタンを出さない', async () => {
    jest.spyOn(Auth, 'currentAuthenticatedUser').mockReturnValue(
      Promise.reject(new Error('Not Authenticated User'))
    )
    const wrapper = shallowMount(ExamplePage, {
      localVue,
    });
    await wrapper.setData({
      title: example.Title,
      example: example,
      keywords: keywords,
    })
    expect(wrapper.exists()).toBeTruthy()
    expect(wrapper.find('annotate-buttons-stub').exists()).toBe(false)
    expect(wrapper.find('tweet-annotate-buttons-stub').exists()).toBe(false)
    expect(wrapper.element).toMatchSnapshot()

    // 無理やり書き換えると表示できることを確認
    await wrapper.setData({
      title: example.Title,
      example: example,
      keywords: keywords,
      isAdmin: true
    })
    expect(wrapper.find('annotate-buttons-stub').exists()).toBe(true)
    expect(wrapper.element).toMatchSnapshot()
  })
  test('amplifyを通っているときはアノテーションボタンを出す', async () => {
    jest.spyOn(Auth, 'currentAuthenticatedUser').mockReturnValue(
      Promise.resolve({
        username: "myuser",
        attributes: {
          email: "test@test.com",
          email_verified: true,
          phone_number: "08012345678"
          // ...other attoributes
        }
        // ...other parameters
      })
    );

    const wrapper = shallowMount(ExamplePage, {
      localVue,
    });
    await wrapper.setData({
      title: example.Title,
      example: example,
      keywords: keywords,
    })
    expect(wrapper.exists()).toBeTruthy()
    expect(wrapper.find('annotate-buttons-stub').exists()).toBe(true)
    expect(wrapper.element).toMatchSnapshot()
  })
})
