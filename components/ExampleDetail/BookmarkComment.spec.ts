import { createLocalVue, shallowMount } from '@vue/test-utils'
import { BListGroupItem } from 'bootstrap-vue'
import VueLazyLoad from 'vue-lazyload'
import BookmarkComment from '~/components/ExampleDetail/BookmarkComment.vue'

const localVue = createLocalVue();
localVue.use(VueLazyLoad);
localVue.component('b-list-group-item', BListGroupItem)

const bookmark = {
  "comment": "ブックマークのテストです",
  "user": "syou6162",
  "timestamp": "2020/08/26 14:29",
  "tags": ["tag1", "tag2"]
}

describe('BookmarkComment', () => {
  test('個別事例ページの中のはてなブックマークコメントが正しく表示できる', () => {
    const wrapper = shallowMount(BookmarkComment, {
      localVue,
      propsData: {
        bookmark: bookmark,
      }
    })

    expect(wrapper.exists()).toBeTruthy()
    expect(wrapper.find('.hatena-bookmark-item').exists()).toBe(true)
    expect(wrapper.find('.hatena-bookmark-user-link').text()).toBe(`id:${bookmark.user}`)
    expect(wrapper.find('.hatena-bookmark-comment').text()).toBe(`${bookmark.comment}`)

    expect(wrapper.element).toMatchSnapshot()
  })
})