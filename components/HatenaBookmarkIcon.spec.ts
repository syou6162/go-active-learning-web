import { createLocalVue, mount } from '@vue/test-utils'

import HatenaBookmarkIcon from './HatenaBookmarkIcon.vue'
import VueLazyLoad from 'vue-lazyload'
import { VBPopover } from 'bootstrap-vue'

const localVue = createLocalVue();
localVue.use(VueLazyLoad);
localVue.directive('b-popover', VBPopover)

describe('HatenaBookmarkIcon', () => {
  test('はてブのアイコンが正しく表示される', async () => {
    const wrapper = mount(HatenaBookmarkIcon, {
      localVue,
      propsData: {
        bookmark: {
          user: "syou6162",
          tags: ["tag1", "tag2"]
        }
      },
    })
    expect(wrapper.exists()).toBeTruthy()
    expect(wrapper.props().bookmark.user).toBe("syou6162")
    expect(wrapper.element).toMatchSnapshot()
    expect(wrapper.attributes("alt")).toBe("id:syou6162")
    expect(wrapper.attributes("data-src")).toBeUndefined()

    wrapper.trigger('mouseover')
    await wrapper.vm.$nextTick()
    expect(wrapper.attributes("data-src")).toBe("https://cdn.profile-image.st-hatena.com/users/syou6162/profile.png")
    expect(wrapper.element).toMatchSnapshot()
  })
})