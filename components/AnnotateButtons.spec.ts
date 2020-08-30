import { createLocalVue, mount } from '@vue/test-utils'
import { ButtonPlugin } from 'bootstrap-vue'
import AnnotateButtons from './AnnotateButtons.vue'

const localVue = createLocalVue();
localVue.use(ButtonPlugin)

describe('AnnotateButtons', () => {
  test('ラベルが0のときはunlabeldが有効になっている', () => {
    const wrapper = mount(AnnotateButtons, {
      localVue,
      propsData: {
        example: {
          Id: 11111,
          Label: 0
        }
      },
    })
    expect(wrapper.exists()).toBeTruthy()
    expect(wrapper.findAll('button').length).toBe(3)

    expect(wrapper.findAll('button').at(0).classes()).toContain('btn-outline-primary')
    expect(wrapper.findAll('button').at(1).classes()).toContain('btn-outline-danger')
    expect(wrapper.findAll('button').at(2).classes()).toContain('btn-secondary')

    expect(wrapper.element).toMatchSnapshot()
  })
})
