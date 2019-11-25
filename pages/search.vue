<template>
  <div>
    <input
      v-model="query"
      type="text"
      placeholder="Input search query here"
    >
    <div v-if="loading && query">
      Now loading...
    </div>
    <div v-else-if="error">
      Fail to retrieve from API server. Error: {{ error }}
    </div>
    <div v-else-if="query && results.length == 0">
      No search result for '{{ query }}'
    </div>
    <div v-else>
      <b-card-group columns>
        <example 
          v-for="example in results"
          :key="example.Url"
          :example="example"
          :tweets="example.ReferringTweets"
          :is-admin="isAdmin"
        />
      </b-card-group>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator';
import Example from '~/models/Example'
import URLSearchParams from '@ungap/url-search-params'
import Auth from '@aws-amplify/auth';
import { NewExample } from '~/plugins/util';

@Component({
  components: {
    Example: () => import('~/components/Example.vue')
  },
  watchQuery: ['query'],
  async asyncData(context) {
    let params = new URLSearchParams();
    params.append('query', context.route.query.query);
    let data = await context.app.$axios.$post("/api/search", params);
    let isAdmin = false;
    if (process.browser) {
      await Auth.currentAuthenticatedUser()
        .then(user => {
          isAdmin = true;
        })
        .catch(err => console.log(err))
    }
    return {
      query: context.route.query.query,
      results: data.Examples.map(e => NewExample(e)),
      loading: false,
      isAdmin: isAdmin
    }
  },
  head() {
    let query = this.query || '';
    const title = "ML-News - 「" + query + "」に関する検索結果";

    return {
      title: title,
      meta: [
        {
          name: "description",
          content: this.results.map(e => e.Title).join("\n")
        }
      ],
      link: [
        {
          rel: "canonical",
          href: "https://www.machine-learning.news/search?query=" + this.query
        }
      ]
    };
  }
})

export default class SearchPage extends Vue {
  query: string | null = null;
  results: Example[] = [];
  error: string | null = null;
  loading: boolean = true;
  isAdmin: boolean = false;

  @Watch('query')
  onQueryChange(newSearchQuery, oldSearchQuery): void {
    // https://router.vuejs.org/ja/guide/essentials/navigation.html
    this.$router.push({ query: { query: newSearchQuery }});
    this.loading = true;
  }
  mounted() {
    Auth.currentAuthenticatedUser()
      .then(user => {
        this.isAdmin = true;
      })
      .catch(err => console.log(err))
  }
}
</script>
