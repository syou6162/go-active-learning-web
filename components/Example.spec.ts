import { createLocalVue, mount } from '@vue/test-utils'
import { BButton, BCard, BCardBody, BCardFooter } from 'bootstrap-vue'
import Example from './Example.vue'

const localVue = createLocalVue();
localVue.component('b-button', BButton)
localVue.component('b-card', BCard)
localVue.component('b-card-body', BCardBody)
localVue.component('b-card-footer', BCardFooter)

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

describe('Example', () => {
  test('ラベルが0のときはunlabeldが有効になっている', () => {
    const wrapper = mount(Example, {
      localVue,
      propsData: {
        example: {
          Id: 11111,
          Label: 0,
          Url: "https://blog.trocco.io/dev-articles/embulk-digdag-meetup-2020",
          FinalUrl: "https://blog.trocco.io/dev-articles/embulk-digdag-meetup-2020",
          Title: "Embulk & Digdag Online Meetup 2020【イベントレポート】 | trocco(トロッコ)",
          Description: "",
          HatenaBookmark: hatenaBookmark,
        },
        tweets: {
          Tweets: [],
          Count: 0
        },
        isAdmin: false
      }
    })
    expect(wrapper.exists()).toBeTruthy()
    expect(wrapper.find('h5').classes()).toContain("card-title")
    expect(wrapper.element).toMatchSnapshot()
  })
})