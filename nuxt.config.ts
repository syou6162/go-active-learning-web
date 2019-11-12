import { Configuration } from '@nuxt/types'
import pkg from './package.json'

const nuxtConfig: Configuration = {
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
    // Doc: https://bootstrap-vue.js.org/docs/
    'bootstrap-vue/nuxt',
    'nuxt-logrocket',
    '@nuxtjs/sentry',
  ],
  // https://github.com/nuxt-community/nuxt-logrocket
  logRocket: {
    logRocketId: '7atsfy/ml-news',
    devModeAllowed: false,
  },
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
  sentry: {
    dsn: process.env.SENTRY_DSN || "",
    publishRelease: process.env.NODE_ENV === "production",
    disabled: process.env.NODE_ENV != "production",
    config: {
      release: process.env.VERSION,
    }
  },

  /*
  ** Build configuration
  */
  build: {
    /*
    ** You can extend webpack config here
    */
    extend(config, ctx) {
      config.devtool = ctx.isClient ? 'eval-source-map' : 'inline-source-map'
    }
  },
  buildModules: ['@nuxt/typescript-build'],
}

export default nuxtConfig
