<template>
  <div>
    <input v-model="query" type="text" placeholder="Input search query here"></input>
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
          v-bind:key="example.Url"
          v-bind:example="example"
          v-bind:isAdmin="isAdmin"
          ></example>
      </b-card-group>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import { Auth } from 'aws-amplify';
import Example from '~/components/Example.vue';
import { NewExample } from '~/assets/util';

export default {
  watchQuery: ['query'],
  mounted() {
    Auth.currentAuthenticatedUser()
      .then(user => {
        this.isAdmin = true;
      })
      .catch(err => console.log(err))
  },
  data () {
    return {
      query: this.$route.query.query,
      results: [],
      error: null,
      loading: true,
      isAdmin: false,
    }
  },
  watch: {
    query: function(newSearchQuery, oldSearchQuery) {
      // https://router.vuejs.org/ja/guide/essentials/navigation.html
      this.$router.push({ query: { query: newSearchQuery }})
      this.loading = true;
    },
  },
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
  },
  components: {
    "example": Example
  }
}
</script>
