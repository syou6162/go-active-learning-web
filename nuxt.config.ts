import NuxtConfiguration from '@nuxt/config'
import pkg from './package.json'

const config: NuxtConfiguration = {
  mode: 'universal',

  /*
  ** Headers of the page
  */
  head: {
    title: pkg.name,
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { property: 'og:type', content: 'website' },
      { name: 'google-site-verification', content: 'H2b0gd2dMVoTm8CCOpRmL9mZYBd9GlkqESRlH2cONhY' }
    ],
    link: [
      { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }
    ]
  },

  /*
  ** Customize the progress-bar color
  */
  loading: { color: '#fff' },

  /*
  ** Global CSS
  */
  css: [
  ],

  /*
  ** Plugins to load before mounting the App
  */
  plugins: [
    { src: "@/plugins/filters.ts", ssr: true },
    { src: "@/plugins/lazyload.ts", ssr: false },
    { src: "@/plugins/amplify.ts", ssr: false }
  ],

  /*
  ** Nuxt.js modules
  */
  modules: [
    // Doc: https://github.com/nuxt-community/axios-module#usage
    '@nuxtjs/axios',
    '@nuxtjs/proxy',
    ['@nuxtjs/google-analytics', {
      id: 'UA-591180-8'
    }],
    ['@nuxtjs/google-adsense', {
      id: 'ca-pub-3611136298489092',
      pageLevelAds: true
    }],
    // Doc: https://bootstrap-vue.js.org/docs/
    'bootstrap-vue/nuxt'
  ],
  /*
  ** Axios module configuration
  */
  axios: {
    // See https://github.com/nuxt-community/axios-module#options
    proxy: true
  },
  proxy: {
    '/api/': 'http://localhost:7778',
  },

  /*
  ** Build configuration
  */
  build: {
    // https://github.com/bootstrap-vue/bootstrap-vue/issues/3397#issuecomment-496835985
    transpile: ['bootstrap-vue'],
    /*
    ** You can extend webpack config here
    */
    extend(config, ctx) {
    }
  }
}

export default config
